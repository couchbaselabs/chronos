//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	log "github.com/couchbase/clog"
	"github.com/couchbaselabs/chronos/widgets"
)

// Structure of chunk sent by the server
type message struct {
	Stats               map[string]float64 `json:"stats,omitempty"`
	RebalanceInProgress bool               `json:"rebalance,omitempty"`
}

// Parameters used by each polling routine
type updateStatsParams struct {
	username     string
	password     string
	stats        *stats
	nodeName     string
	errChannel   chan *errorMsg
	eventChannel chan *widgets.Event
	popupChannel chan string
	killSwitch   chan bool
	timeDiff     float64
}

// Errors encountered sent to the main routine
type errorMsg struct {
	name        error
	description string

	// Flag to signal termination of the program
	terminate bool
}

// Function to consolidate all the input parameters into a struct
func newUpdateStatsParams(username string, password string,
	stats *stats, node string, errChannel chan *errorMsg,
	eventChannel chan *widgets.Event,
	popupChannel chan string, killSwitch chan bool) *updateStatsParams {

	return &updateStatsParams{
		username:     username,
		password:     password,
		stats:        stats,
		nodeName:     node,
		errChannel:   errChannel,
		eventChannel: eventChannel,
		popupChannel: popupChannel,
		killSwitch:   killSwitch,
		timeDiff:     0,
	}
}

// Initializes errors into a struct to be sent through channels
func newErrorMsg(name error, description string, terminate bool) *errorMsg {

	return &errorMsg{
		name:        name,
		description: description,
		terminate:   terminate,
	}
}

// Polling stats from a node every second and updating to stats struct
func updateStats(params *updateStatsParams) int {

	log.Printf("Started %v", params.nodeName)
	url := params.nodeName + "/api/statsStream"

	// Making the http request
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(params.username, params.password)

	if err != nil {
		params.errChannel <- newErrorMsg(
			err, "update_stats: Cannot connect to server"+
				params.nodeName+":"+err.Error(), false,
		)
		return 0
	}

	// Sending the http request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		params.errChannel <- newErrorMsg(
			err, "update_stats: Invalid http response from server"+
				params.nodeName+":"+err.Error(), false,
		)
		return 0
	}

	// If status code is not OK
	if resp.StatusCode != http.StatusOK {
		params.errChannel <- newErrorMsg(
			err, "update_stats: Status code is not OK:"+
				fmt.Sprintf("%d", resp.StatusCode)+resp.Status, false,
		)
		return 0
	}

	dec := json.NewDecoder(resp.Body)

	// Setting the timer for checking if server sent a chunk through the response
	// Should be equal to the time defined in the server
	updateTicker := time.NewTicker(time.Second).C

	// Main polling loop
	for {
		select {
		case <-updateTicker:

			// Decode chunk sent by the server
			var m message
			err := dec.Decode(&m)
			if err != nil {
				if err == io.EOF {
					params.errChannel <- newErrorMsg(
						err, "update_stats: Server closed connection"+
							params.nodeName+err.Error(), false,
					)
					return 0
				}
				params.errChannel <- newErrorMsg(
					err, "update_stats: Invalid message recieved"+
						params.nodeName+err.Error(), false,
				)
				return 0
			}

			// Send a message to the main routine if node is under rebalance
			if m.RebalanceInProgress {
				params.popupChannel <- "rebalance"
			}

			// Note time before updating for accurate calculations across commands
			curTime := time.Now()

			// Update arrival time of the chunk
			params.stats.timeLock.Lock()

			sec := 1

			// If this is not the first update for the node
			if !params.stats.arrivalTimes[params.nodeName][len(params.stats.arrivalTimes[params.nodeName])-1].IsZero() {

				// Time passed from the last time stats were updated
				diffTime := curTime.Sub(
					params.stats.arrivalTimes[params.nodeName][len(params.stats.arrivalTimes[params.nodeName])-1],
				)
				diffSec := diffTime.Seconds()

				// Number of seconds to update
				sec = int(math.Round(diffSec + params.timeDiff))

				// Excess time ignored in rounding
				params.timeDiff = diffSec + params.timeDiff - float64(sec)
			}

			// If response is delayed
			if sec > 1 {
				params.popupChannel <- params.nodeName
			}

			// Update unknown times
			for i := 0; i < sec-1; i++ {
				params.stats.arrivalTimes[params.nodeName] =
					params.stats.arrivalTimes[params.nodeName][1:]

				params.stats.arrivalTimes[params.nodeName] =
					append(
						params.stats.arrivalTimes[params.nodeName],
						time.Time{},
					)
			}

			// Update current time
			params.stats.arrivalTimes[params.nodeName] =
				params.stats.arrivalTimes[params.nodeName][1:]

			params.stats.arrivalTimes[params.nodeName] =
				append(params.stats.arrivalTimes[params.nodeName], curTime)

			params.stats.timeLock.Unlock()

			params.stats.bufferLock.Lock()

			// Update unknown stats
			for i := 0; i < sec-1; i++ {
				for _, stat := range params.stats.statsList {
					params.stats.statBuffers[params.nodeName][stat] =
						params.stats.statBuffers[params.nodeName][stat][1:]

					params.stats.statBuffers[params.nodeName][stat] =
						append(
							params.stats.statBuffers[params.nodeName][stat],
							params.stats.statBuffers[params.nodeName][stat][len(params.stats.statBuffers[params.nodeName][stat])-1],
						)
				}
			}
			params.stats.bufferLock.Unlock()

			// Update current stat slices
			for _, stat := range params.stats.statsList {
				val, ok := m.Stats[stat]
				params.stats.bufferLock.Lock()

				// Remove first element
				params.stats.statBuffers[params.nodeName][stat] =
					params.stats.statBuffers[params.nodeName][stat][1:]

				// If chunk has the stat, update the buffer with the value
				// else update with 0 as default (should never occur)
				if ok {
					params.stats.statBuffers[params.nodeName][stat] =
						append(
							params.stats.statBuffers[params.nodeName][stat], val,
						)
				} else {
					params.stats.statBuffers[params.nodeName][stat] =
						append(
							params.stats.statBuffers[params.nodeName][stat], 0,
						)
				}
				params.stats.bufferLock.Unlock()

				// Run analysis on the newest data
				analyzeStat(
					params.stats, params.nodeName, stat, params.eventChannel,
				)
			}
		// Kill the routine if node is no longer part of the cluster
		case <-params.killSwitch:
			params.errChannel <- newErrorMsg(
				nil, "update_stats: Node not part of cluster anymore: "+
					params.nodeName, false,
			)
			// Indicates to the exponential loop to exit out of it
			return -1
		}
	}
}

// Exponential backoff loop for connection to the node
func updateStatsExponentialBackoff(params *updateStatsParams) {

	startSleepMS := 500
	backoffFactor := 1.5
	maxSleepMS := 5000
	nextSleepMS := startSleepMS

	for {
		select {
		// Exit out of backoff loop if node no longer in the cluster
		case <-params.killSwitch:
			params.errChannel <- newErrorMsg(
				nil, "update_stats: Node not part of cluster anymore: "+
					params.nodeName, false)
			return
		// Retry with exponential backoff
		default:
			val := updateStats(params)

			time.Sleep(time.Duration(nextSleepMS) * time.Millisecond)
			nextSleepMS = int(float64(nextSleepMS) * backoffFactor)

			if nextSleepMS > maxSleepMS {
				nextSleepMS = maxSleepMS
			}

			if val == -1 {
				return
			}
		}
	}
}
