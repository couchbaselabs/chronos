//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"fmt"
	"math"
	"time"

	log "github.com/couchbase/clog"
	"github.com/couchbaselabs/chronos/widgets"
)

// Constant used to avoid warning in fmt.Sprintf()
const (
	percent = "%"
)

// Analyse the latest value of a stat for a particular node using
// the alert thresholds
func analyzeStat(stats *stats, node string, stat string,
	eventChannel chan *widgets.Event) {

	stats.bufferLock.RLock()

	length := len(stats.statBuffers[node][stat])
	curVal := stats.statBuffers[node][stat][length-1]
	lastTimeVal := math.NaN()

	// Check if max change can be calculated
	// eg, cannot calculate change with only one valid value
	if *stats.statInfo[stat].MaxChangeTime <= length-1 &&
		!math.IsNaN(*stats.statInfo[stat].MaxChange) {

		lastTimeVal =
			stats.statBuffers[node][stat][length-1-*stats.statInfo[stat].MaxChangeTime]
	}

	stats.bufferLock.RUnlock()

	// Check for minimum threshold
	if curVal < *stats.statInfo[stat].MinVal &&
		!math.IsNaN(*stats.statInfo[stat].MinVal) {

		event := widgets.NewEvent(
			node, stat, "Below Threshold", curVal, *stats.statInfo[stat].MinVal,
		)

		event.Description = makeDescription(event, false)

		triggerEvent(event, eventChannel, stats)

	}

	// Check for maximum threshold
	if curVal > *stats.statInfo[stat].MaxVal &&
		!math.IsNaN(*stats.statInfo[stat].MaxVal) {

		event := widgets.NewEvent(
			node, stat, "Above Threshold", curVal, *stats.statInfo[stat].MaxVal,
		)
		event.Description = makeDescription(event, false)

		triggerEvent(event, eventChannel, stats)

	}

	// Check for maximum change
	if !math.IsNaN(lastTimeVal) {
		if math.Abs(curVal-lastTimeVal)/lastTimeVal >
			*stats.statInfo[stat].MaxChange &&
			!math.IsNaN(*stats.statInfo[stat].MaxChange) {

			event := widgets.NewEvent(
				node, stat, "Sudden Change", curVal,
				*stats.statInfo[stat].MaxChange,
			)
			event.ThresholdChange = math.Abs(curVal-lastTimeVal) / lastTimeVal
			event.ThresholdTime = *stats.statInfo[stat].MaxChangeTime
			event.Description = makeDescription(event, false)

			triggerEvent(event, eventChannel, stats)

		}
	}
}

// Send the event to the event creation handler after basic checks
// Variable to accomodate tests
var triggerEvent = func(event *widgets.Event,
	eventChannel chan *widgets.Event, stats *stats) {

	switch event.EventType {
	case "Sudden Change":
		stats.timeLock.RLock()
		index := len(stats.arrivalTimes[event.Node]) - event.ThresholdTime - 1
		if !stats.arrivalTimes[event.Node][index].IsZero() {
			eventChannel <- event
		}
		stats.timeLock.RUnlock()
	default:
		eventChannel <- event
	}
}

// Handles incoming alerts. Adds to event display if alert is new, updates
// an existing alert if it already exists
func eventCreateHandler(eventChannel chan *widgets.Event,
	eventDisplay *widgets.EventDisplay, stats *stats,
	alerts map[string]*int) {

	// Main loop for event creation handling
	for {
		event := <-eventChannel
		created := false
		eventDisplay.EventLock.Lock()
		for _, prevEvent := range eventDisplay.Events {

			eventTTL := prevEvent.LastTriggered.Add(
				time.Duration(*alerts["ttl"]) * time.Second,
			)

			// Check if alert already exists
			if event.Node == prevEvent.Node &&
				event.Stat == prevEvent.Stat &&
				event.EventType == prevEvent.EventType &&
				event.LastTriggered.Before(eventTTL) {
				updateEvent(prevEvent)
				created = true
				break
			}
		}
		eventDisplay.EventLock.Unlock()

		// Make alert if it doesn't exist
		if !created {

			event = createEvent(event, stats, alerts)

			if event != nil {
				eventDisplay.AddEvent(event)
			}
		}
	}
}

// Update existing alert with the latest alert
func updateEvent(event *widgets.Event) {
	event.AlertTimes = append(event.AlertTimes, time.Now())
	event.LastTriggered = time.Now()
	event.NumTimes++
	event.Description = makeDescription(event, true)
}

// Create a new alert and append required data and data arrival times
func createEvent(event *widgets.Event, stats *stats,
	alerts map[string]*int) *widgets.Event {

	alertStartTime := event.LastTriggered.Add(
		-time.Duration(*alerts["dataPadding"]) * time.Second,
	)
	event.DataStart = alertStartTime
	startIndex := 0
	event.AlertTimes = append(event.AlertTimes, time.Now())

	stats.timeLock.RLock()

	if _, ok := stats.arrivalTimes[event.Node]; ok {
		for i := len(stats.arrivalTimes[event.Node]) - 1; i >= 0; i-- {
			if stats.arrivalTimes[event.Node][i].IsZero() {
				if i != len(stats.arrivalTimes[event.Node])-1 {
					startIndex = i + 1
				}
				break
			} else if widgets.CompareTimes(
				stats.arrivalTimes[event.Node][i], alertStartTime,
			) {
				startIndex = i
				break
			}
		}

		event.DataTimes = append(
			event.DataTimes, stats.arrivalTimes[event.Node][startIndex:]...,
		)
	} else {
		stats.timeLock.RUnlock()
		return nil
	}

	stats.timeLock.RUnlock()

	stats.bufferLock.RLock()

	if _, ok := stats.statBuffers[event.Node]; ok {
		event.Data = append(
			event.Data, stats.statBuffers[event.Node][event.Stat][startIndex:]...,
		)
	} else {
		stats.bufferLock.RUnlock()
		return nil
	}
	stats.bufferLock.RUnlock()

	return event
}

// Update data routine for all alerts
func eventDataHandler(eventDisplay *widgets.EventDisplay,
	stats *stats, alerts map[string]*int) {

	// Main loop for event data updating
	for range time.Tick(time.Second) {
		updateEventData(eventDisplay, stats, alerts)
	}
}

// One iteration of event data updating
func updateEventData(eventDisplay *widgets.EventDisplay,
	stats *stats, alerts map[string]*int) {

	// Variable to indicate no alerts have expired
	clean := true
	eventDisplay.EventLock.Lock()
	for _, event := range eventDisplay.Events {

		// Check if alert needs to be updated
		// No updating alerts if node no longer in cluster
		// or if alert already at full data capacity
		if !event.DataFilled {

			alertEndTime := event.LastTriggered.Add(
				time.Duration(*alerts["dataPadding"]) * time.Second,
			)

			stats.timeLock.RLock()

			// Check if node still exists
			if _, ok := stats.arrivalTimes[event.Node]; ok {
				i := len(stats.arrivalTimes[event.Node]) - 1

				// Find number of data points to update
				for event.DataTimes[len(event.DataTimes)-1] !=
					stats.arrivalTimes[event.Node][i] {
					i--
				}

				i++
				for i < len(stats.arrivalTimes[event.Node]) {

					stats.bufferLock.RLock()

					// Check if node still exists
					// Update data
					if _, ok := stats.statBuffers[event.Node]; ok {
						event.Data = append(
							event.Data,
							stats.statBuffers[event.Node][event.Stat][i],
						)
						stats.bufferLock.RUnlock()
					} else {
						stats.bufferLock.RUnlock()
						break
					}

					// Update data arrival times
					event.DataTimes = append(
						event.DataTimes, stats.arrivalTimes[event.Node][i],
					)

					// Check for alert fullness
					if widgets.CompareTimes(
						alertEndTime, stats.arrivalTimes[event.Node][i],
					) {
						event.DataFilled = true
						break
					}

					i++
				}
			}

			stats.timeLock.RUnlock()

			if widgets.CompareTimes(
				alertEndTime, event.DataTimes[len(event.DataTimes)-1],
			) {
				event.DataFilled = true
			}
		}

		eventTTL := event.LastTriggered.Add(
			time.Duration(*alerts["ttl"]) * time.Second,
		)

		// Check if alert timed out
		if time.Now().After(eventTTL) {
			event.Stale = true
			clean = false
		}
	}

	eventDisplay.EventLock.Unlock()

	if !clean {

		// Remove stale alerts
		deletedEvents := cleanEvents(eventDisplay)

		// Log removed alerts
		for _, event := range deletedEvents {
			reportDelete(event)
		}
	}
}

// Separate all alerts that needs to be deleted from the event display
func cleanEvents(eventDisplay *widgets.EventDisplay) []*widgets.Event {

	eventDisplay.EventLock.Lock()
	cleanEvents := make([]*widgets.Event, 0)
	deletedEvents := make([]*widgets.Event, 0)

	for _, event := range eventDisplay.Events {
		if event.Stale {
			deletedEvents = append(deletedEvents, event)
		} else {
			cleanEvents = append(cleanEvents, event)
		}
	}

	eventDisplay.Events = cleanEvents

	eventDisplay.ResetSelect()
	eventDisplay.EventLock.Unlock()

	return deletedEvents
}

// Makes the description to be displayed in the UI for an event
func makeDescription(event *widgets.Event, update bool) string {

	switch update {
	case true:
		switch event.EventType {
		case "Below Threshold", "Above Threshold":
			return fmt.Sprintf(
				"%s:- %s:%s Event - %s, Value - %f, Threshold Value - %f"+
					", Last occured at %s, Occured %d times",
				event.FirstTriggered.Format("2006-01-02 15:04:05"), event.Node,
				event.Stat, event.EventType, event.ThresholdData, event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"), event.NumTimes,
			)
		case "Sudden Change":
			return fmt.Sprintf(
				"%s:- %s:%s Event - %s, Value - %f, Threshold Value - %f"+
					", Last occured at %s with a change of %f over %d seconds,"+
					" Occured %d times",
				event.FirstTriggered.Format("2006-01-02 15:04:05"), event.Node,
				event.Stat, event.EventType, event.ThresholdData, event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
				event.ThresholdChange, event.ThresholdTime, event.NumTimes,
			)
		}
	case false:
		switch event.EventType {
		case "Below Threshold", "Above Threshold":
			return fmt.Sprintf(
				"%s:- %s:%s Event - %s, Value - %f, Threshold Value - %f",
				event.FirstTriggered.Format("2006-01-02 15:04:05"), event.Node,
				event.Stat, event.EventType, event.ThresholdData, event.Threshold,
			)
		case "Sudden Change":
			return fmt.Sprintf(
				"%s:- %s:%s Event - %s, Value - %f, Threshold Value - %f"+
					", Occured over %d seconds with a change of %f",
				event.FirstTriggered.Format("2006-01-02 15:04:05"), event.Node,
				event.Stat, event.EventType, event.ThresholdData, event.Threshold,
				event.ThresholdTime, event.ThresholdChange,
			)
		}
	}
	return ""
}

// Makes the description for an event to be used while logging
var reportDelete = func(event *widgets.Event) {

	switch event.EventType {
	case "Sudden Change":
		if event.NumTimes == 1 {
			log.Printf(
				"Node - %s, Stat - %s. Stat changed by more "+
					"than the threshold limit of %f%s at %s. This"+
					"change occured over %d second(s)",
				event.Node, event.Stat, event.Threshold*100,
				percent, event.LastTriggered.Format("2006-01-02 15:04:05"),
				event.ThresholdTime,
			)
		} else {
			log.Printf(
				"Node - %s, Stat - %s. Stat changed by more than"+
					" the threshold limit of %f%s at %s. This change"+
					" occured over %d second(s). Similar changes occured"+
					" %d times with the last one occuring at %s",
				event.Node, event.Stat, event.Threshold*100,
				percent, event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.ThresholdTime, event.NumTimes,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	case "Below Threshold":
		if event.NumTimes == 1 {
			log.Printf("Node - %s, Stat - %s. Stat exceeded "+
				"threshold limit of %f at %s",
				event.Node, event.Stat, event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		} else {
			log.Printf("Node - %s, Stat - %s. Stat exceeded "+
				"threshold limit of %f at %s. Similarly, the stat exceeded"+
				" threshold limit %d times with the last one occuring at %s",
				event.Node, event.Stat, event.Threshold,
				event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.NumTimes, event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	case "Above Threshold":
		if event.NumTimes == 1 {
			log.Printf("Node - %s, Stat - %s. Stat dropped below threshold"+
				" limit of %f at %s",
				event.Node, event.Stat, event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		} else {
			log.Printf("Node - %s, Stat - %s. Stat dropped below threshold"+
				" limit of %f at %s. Similarly, the stat was below the threshold"+
				" limit %d times with the last one occuring at %s",
				event.Node, event.Stat, event.Threshold,
				event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.NumTimes, event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	}
}
