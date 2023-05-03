//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/couchbaselabs/chronos/widgets"
)

func TestAnalyzer(t *testing.T) {

	inputStats := []*stats{
		{
			statBuffers: map[string]map[string][]float64{
				"node1": {
					"stat1": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat1": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node2": {
					"stat2": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat2": {
					MinVal:        math.NaN(),
					MaxVal:        70.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node3": {
					"stat3": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat3": {
					MinVal:        math.NaN(),
					MaxVal:        65.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node4": {
					"stat4": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat4": {
					MinVal:        math.NaN(),
					MaxVal:        55.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node5": {
					"stat5": {
						70.0, 60.0, 50.0, 40.0, 30.0, 20.0, 10.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat5": {
					MinVal:        10.0,
					MaxVal:        math.NaN(),
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node6": {
					"stat6": {
						70.0, 60.0, 50.0, 40.0, 30.0, 20.0, 10.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat6": {
					MinVal:        15.0,
					MaxVal:        math.NaN(),
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node7": {
					"stat7": {
						70.0, 60.0, 50.0, 40.0, 30.0, 20.0, 10.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat7": {
					MinVal:        25.0,
					MaxVal:        math.NaN(),
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node8": {
					"stat8": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat8": {
					MinVal:        70.0,
					MaxVal:        70.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node9": {
					"stat9": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat9": {
					MinVal:        75.0,
					MaxVal:        65.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node10": {
					"stat10": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat10": {
					MinVal:        65.0,
					MaxVal:        75.0,
					MaxChange:     math.NaN(),
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node11": {
					"stat11": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat11": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 0,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node12": {
					"stat12": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat12": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node13": {
					"stat13": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat13": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 2,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node14": {
					"stat14": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat14": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 6,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node15": {
					"stat15": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat15": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 7,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node16": {
					"stat16": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat16": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 8,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node17": {
					"stat17": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 120.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat17": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     1.0,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node18": {
					"stat18": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 120.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat18": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     0.99,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node19": {
					"stat19": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 120.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat19": {
					MinVal:        math.NaN(),
					MaxVal:        math.NaN(),
					MaxChange:     1.01,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node20": {
					"stat20": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat20": {
					MinVal:        75.0,
					MaxVal:        math.NaN(),
					MaxChange:     0.1,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node21": {
					"stat21": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat21": {
					MinVal:        math.NaN(),
					MaxVal:        65.0,
					MaxChange:     0.1,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
		{
			statBuffers: map[string]map[string][]float64{
				"node22": {
					"stat22": {
						10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0,
					},
				},
			},
			statInfo: map[string]*configStatInfo{
				"stat22": {
					MinVal:        75.0,
					MaxVal:        65.0,
					MaxChange:     0.1,
					MaxChangeTime: 1,
				},
			},
			bufferLock: sync.RWMutex{},
		},
	}

	inputNode := []string{
		"node1",
		"node2",
		"node3",
		"node4",
		"node5",
		"node6",
		"node7",
		"node8",
		"node9",
		"node10",
		"node11",
		"node12",
		"node13",
		"node14",
		"node15",
		"node16",
		"node17",
		"node18",
		"node19",
		"node20",
		"node21",
		"node22",
	}

	inputStat := []string{
		"stat1",
		"stat2",
		"stat3",
		"stat4",
		"stat5",
		"stat6",
		"stat7",
		"stat8",
		"stat9",
		"stat10",
		"stat11",
		"stat12",
		"stat13",
		"stat14",
		"stat15",
		"stat16",
		"stat17",
		"stat18",
		"stat19",
		"stat20",
		"stat21",
		"stat22",
	}

	outputEvents := [][]*widgets.Event{
		{},
		{},
		{
			{
				Node:            "node3",
				Stat:            "stat3",
				EventType:       "Above Threshold",
				Threshold:       65.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
		},
		{
			{
				Node:            "node4",
				Stat:            "stat4",
				EventType:       "Above Threshold",
				Threshold:       55.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
		},
		{},
		{
			{
				Node:            "node6",
				Stat:            "stat6",
				EventType:       "Below Threshold",
				Threshold:       15.0,
				ThresholdData:   10.0,
				ThresholdChange: 0.0,
			},
		},
		{
			{
				Node:            "node7",
				Stat:            "stat7",
				EventType:       "Below Threshold",
				Threshold:       25.0,
				ThresholdData:   10.0,
				ThresholdChange: 0.0,
			},
		},
		{},
		{
			{
				Node:            "node9",
				Stat:            "stat9",
				EventType:       "Below Threshold",
				Threshold:       75.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
			{
				Node:            "node9",
				Stat:            "stat9",
				EventType:       "Above Threshold",
				Threshold:       65.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
		},
		{},
		{},
		{
			{
				Node:            "node12",
				Stat:            "stat12",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-60.0) / 60.0,
			},
		},
		{
			{
				Node:            "node13",
				Stat:            "stat13",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-50.0) / 50.0,
			},
		},
		{
			{
				Node:            "node14",
				Stat:            "stat14",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-10.0) / 10.0,
			},
		},
		{},
		{},
		{},
		{
			{
				Node:            "node18",
				Stat:            "stat18",
				EventType:       "Sudden Change",
				Threshold:       0.99,
				ThresholdData:   120.0,
				ThresholdChange: math.Abs(120.0-60.0) / 60.0,
			},
		},
		{},
		{
			{
				Node:            "node20",
				Stat:            "stat20",
				EventType:       "Below Threshold",
				Threshold:       75.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
			{
				Node:            "node20",
				Stat:            "stat20",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-60.0) / 60.0,
			},
		},
		{
			{
				Node:            "node21",
				Stat:            "stat21",
				EventType:       "Above Threshold",
				Threshold:       65.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
			{
				Node:            "node21",
				Stat:            "stat21",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-60.0) / 60.0,
			},
		},
		{
			{
				Node:            "node22",
				Stat:            "stat22",
				EventType:       "Below Threshold",
				Threshold:       75.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
			{
				Node:            "node22",
				Stat:            "stat22",
				EventType:       "Above Threshold",
				Threshold:       65.0,
				ThresholdData:   70.0,
				ThresholdChange: 0.0,
			},
			{
				Node:            "node22",
				Stat:            "stat22",
				EventType:       "Sudden Change",
				Threshold:       0.1,
				ThresholdData:   70.0,
				ThresholdChange: math.Abs(70.0-60.0) / 60.0,
			},
		},
	}

	events := make([]*widgets.Event, 0)

	TriggerEventOri := triggerEvent
	triggerEvent = func(event *widgets.Event, eventChannel chan *widgets.Event, stats *stats) {
		events = append(events, event)
	}

	for i, stats := range inputStats {

		analyzeStat(stats, inputNode[i], inputStat[i], nil)

		correct := true
		if len(events) == len(outputEvents[i]) {

			for j, event := range events {
				if outputEvents[i][j].Node == event.Node &&
					outputEvents[i][j].Stat == event.Stat &&
					outputEvents[i][j].EventType == event.EventType &&
					outputEvents[i][j].Threshold == event.Threshold &&
					outputEvents[i][j].ThresholdData == event.ThresholdData &&
					outputEvents[i][j].ThresholdChange == event.ThresholdChange {
					continue
				} else {
					correct = false
					break
				}
			}
		} else {
			correct = false
		}

		if !correct {
			t.Errorf("Expected %v got %v %d", outputEvents[i], events, i)
		}

		events = make([]*widgets.Event, 0)
	}

	triggerEvent = TriggerEventOri
}

func intPointer(val int) *int {
	return &val
}

func TestTriggerEvent(t *testing.T) {

	eventChannel := make(chan *widgets.Event)
	testCases := []struct {
		event *widgets.Event
		stat  *stats
		count int
	}{
		{
			event: &widgets.Event{
				Node:          "node1",
				EventType:     "Above Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node1": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node2",
				EventType:     "Above Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node2": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node3",
				EventType:     "Above Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node3": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node4",
				EventType:     "Above Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node4": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node5",
				EventType:     "Above Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node5": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node6",
				EventType:     "Above Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node6": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node7",
				EventType:     "Above Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node7": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node8",
				EventType:     "Above Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node8": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node9",
				EventType:     "Above Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node9": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node10",
				EventType:     "Below Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node10": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node11",
				EventType:     "Below Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node11": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node12",
				EventType:     "Below Threshold",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node12": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node13",
				EventType:     "Below Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node13": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node14",
				EventType:     "Below Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node14": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node15",
				EventType:     "Below Threshold",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node15": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node16",
				EventType:     "Below Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node16": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node17",
				EventType:     "Below Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node17": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node18",
				EventType:     "Below Threshold",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node18": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node19",
				EventType:     "Sudden Change",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node19": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node20",
				EventType:     "Sudden Change",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node20": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node21",
				EventType:     "Sudden Change",
				ThresholdTime: 0,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node21": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node22",
				EventType:     "Sudden Change",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node22": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 0,
		},
		{
			event: &widgets.Event{
				Node:          "node23",
				EventType:     "Sudden Change",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node23": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node24",
				EventType:     "Sudden Change",
				ThresholdTime: 1,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node24": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
		{
			event: &widgets.Event{
				Node:          "node25",
				EventType:     "Sudden Change",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node25": {
						time.Time{}, time.Time{}, time.Now(),
					},
				},
			},
			count: 0,
		},
		{
			event: &widgets.Event{
				Node:          "node26",
				EventType:     "Sudden Change",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node26": {
						time.Time{}, time.Now(), time.Now(),
					},
				},
			},
			count: 0,
		},
		{
			event: &widgets.Event{
				Node:          "node27",
				EventType:     "Sudden Change",
				ThresholdTime: 2,
			},
			stat: &stats{
				timeLock: sync.RWMutex{},
				arrivalTimes: map[string][]time.Time{
					"node27": {
						time.Now(), time.Now(), time.Now(),
					},
				},
			},
			count: 1,
		},
	}

	for i, testCase := range testCases {
		count := 0

		go triggerEvent(testCase.event, eventChannel, testCase.stat)

		timeout := time.After(time.Duration(200) * time.Millisecond)

		select {
		case <-eventChannel:
			count++
		case <-timeout:
		}

		if count != testCase.count {
			t.Errorf("Expected %v got %v %d", testCase.count, count, i)
		}
	}
}

func TestEventCreateHandler(t *testing.T) {

	curTime := time.Now()

	testCases := []struct {
		event        *widgets.Event
		prevEvents   []*widgets.Event
		arrivalTimes map[string][]time.Time
		statBuffers  map[string]map[string][]float64
		alerts       map[string]*int
		create       bool
		eventDataLen int
		numTimes     int
	}{
		{
			event: &widgets.Event{
				Node:          "node1",
				Stat:          "stat1",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{
					Node:          "node1",
					Stat:          "stat1",
					EventType:     "Above Threshold",
					LastTriggered: curTime.Add(-time.Second),
					AlertTimes:    make([]time.Time, 0),
					NumTimes:      1,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node1": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node1": {
					"stat1": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			create:       false,
			eventDataLen: 0,
			numTimes:     2,
		},
		{
			event: &widgets.Event{
				Node:          "node2",
				Stat:          "stat2",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{
					Node:          "node2",
					Stat:          "stat2",
					EventType:     "Above Threshold",
					LastTriggered: curTime.Add(-time.Second * time.Duration(2)),
					AlertTimes:    make([]time.Time, 0),
					NumTimes:      1,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node2": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node2": {
					"stat2": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			create:       false,
			eventDataLen: 0,
			numTimes:     2,
		},
		{
			event: &widgets.Event{
				Node:          "node3",
				Stat:          "stat3",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{
					Node:          "node3",
					Stat:          "stat3",
					EventType:     "Above Threshold",
					LastTriggered: curTime.Add(-time.Second * time.Duration(3)),
					AlertTimes:    make([]time.Time, 0),
					NumTimes:      1,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node3": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node3": {
					"stat3": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			create:       false,
			eventDataLen: 0,
			numTimes:     1,
		},
		{
			event: &widgets.Event{
				Node:          "node4",
				Stat:          "stat4",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node4": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node4": {
					"stat4": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(0),
			},
			create:       false,
			eventDataLen: 1,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node5",
				Stat:          "stat5",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node5": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node5": {
					"stat5": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(1),
			},
			create:       false,
			eventDataLen: 2,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node6",
				Stat:          "stat6",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node6": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node6": {
					"stat6": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			create:       false,
			eventDataLen: 3,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node7",
				Stat:          "stat7",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node7": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node7": {
					"stat7": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(3),
			},
			create:       false,
			eventDataLen: 4,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node8",
				Stat:          "stat8",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node8": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node8": {
					"stat8": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(4),
			},
			create:       false,
			eventDataLen: 4,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node9",
				Stat:          "stat9",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node9": {
					time.Time{}, time.Time{}, curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node9": {
					"stat9": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(0),
			},
			create:       false,
			eventDataLen: 1,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node10",
				Stat:          "stat10",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node10": {
					time.Time{}, time.Time{}, curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node10": {
					"stat10": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(1),
			},
			create:       false,
			eventDataLen: 2,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node11",
				Stat:          "stat11",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node11": {
					time.Time{}, time.Time{}, curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node11": {
					"stat11": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			create:       false,
			eventDataLen: 3,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node12",
				Stat:          "stat12",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node12": {
					time.Time{}, time.Time{}, curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node12": {
					"stat12": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(3),
			},
			create:       false,
			eventDataLen: 4,
			numTimes:     0,
		},
		{
			event: &widgets.Event{
				Node:          "node13",
				Stat:          "stat13",
				EventType:     "Above Threshold",
				LastTriggered: curTime,
				AlertTimes:    make([]time.Time, 0),
				DataTimes:     make([]time.Time, 0),
				Data:          make([]float64, 0),
			},
			prevEvents: []*widgets.Event{
				{},
			},
			arrivalTimes: map[string][]time.Time{
				"node13": {
					time.Time{}, time.Time{}, curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node13": {
					"stat13": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(4),
			},
			create:       false,
			eventDataLen: 4,
			numTimes:     0,
		},
	}

	eventChannel := make(chan *widgets.Event)
	eventDisplay := &widgets.EventDisplay{
		EventLock: sync.RWMutex{},
	}
	stats := &stats{
		timeLock:   sync.RWMutex{},
		bufferLock: sync.RWMutex{},
	}
	alerts := make(map[string]*int)

	go eventCreateHandler(eventChannel, eventDisplay, stats, alerts)

	for i, testCase := range testCases {
		eventDisplay.Events = testCase.prevEvents
		stats.arrivalTimes = testCase.arrivalTimes
		stats.statBuffers = testCase.statBuffers
		alerts["ttl"] = testCase.alerts["ttl"]
		alerts["dataPadding"] = testCase.alerts["dataPadding"]

		eventChannel <- testCase.event
		time.Sleep(time.Duration(200) * time.Millisecond)

		if testCase.create && len(eventDisplay.Events[0].Data) != testCase.eventDataLen {
			t.Errorf("Expected %v got %v %d", testCase.eventDataLen, len(eventDisplay.Events[0].Data), i)
		} else if !testCase.create && eventDisplay.Events[0].NumTimes != testCase.numTimes {
			t.Errorf("Expected %v got %v %d", testCase.numTimes, eventDisplay.Events[0].NumTimes, i)
		}
	}
}

func TestEventDataHandler(t *testing.T) {

	curTime := time.Now()

	deletedEvents := make([]*widgets.Event, 0)

	ReportDeleteOri := reportDelete

	reportDelete = func(event *widgets.Event) {
		deletedEvents = append(deletedEvents, event)
	}

	testCases := []struct {
		events       []*widgets.Event
		arrivalTimes map[string][]time.Time
		statBuffers  map[string]map[string][]float64
		alerts       map[string]*int
		delete       bool
		fill         bool
		update       bool
	}{
		{
			events: []*widgets.Event{
				{
					Node:          "node1",
					Stat:          "stat1",
					DataFilled:    true,
					LastTriggered: curTime.Add(-time.Second * time.Duration(2)),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node1": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node1": {
					"stat1": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: false,
			fill:   true,
			update: false,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node2",
					Stat:          "stat2",
					DataFilled:    true,
					LastTriggered: curTime.Add(-time.Second * time.Duration(3)),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node2": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node2": {
					"stat2": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: true,
			fill:   true,
			update: false,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node3",
					Stat:          "stat3",
					DataFilled:    true,
					LastTriggered: curTime.Add(-time.Second * time.Duration(4)),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node3": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node3": {
					"stat3": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: true,
			fill:   true,
			update: false,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node4",
					Stat:          "stat4",
					DataFilled:    false,
					LastTriggered: curTime.Add(-time.Second),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node4": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node4": {
					"stat4": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: false,
			fill:   false,
			update: false,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node5",
					Stat:          "stat5",
					DataFilled:    false,
					LastTriggered: curTime.Add(-time.Second * time.Duration(2)),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node5": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node5": {
					"stat5": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: false,
			fill:   true,
			update: false,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node6",
					Stat:          "stat6",
					DataFilled:    false,
					LastTriggered: curTime.Add(-time.Second),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second),
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node6": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node6": {
					"stat6": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: false,
			fill:   false,
			update: true,
		},
		{
			events: []*widgets.Event{
				{
					Node:          "node7",
					Stat:          "stat7",
					DataFilled:    false,
					LastTriggered: curTime.Add(-time.Second * time.Duration(2)),
					DataTimes: []time.Time{
						curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second),
					},
					Data: []float64{
						0.0, 0.0, 0.0, 0.0,
					},
					Stale: false,
				},
			},
			arrivalTimes: map[string][]time.Time{
				"node7": {
					curTime.Add(-time.Second * time.Duration(3)), curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime,
				},
			},
			statBuffers: map[string]map[string][]float64{
				"node7": {
					"stat7": {
						0.0, 0.0, 0.0, 0.0,
					},
				},
			},
			alerts: map[string]*int{
				"ttl":         intPointer(3),
				"dataPadding": intPointer(2),
			},
			delete: false,
			fill:   true,
			update: true,
		},
	}

	eventDisplay := &widgets.EventDisplay{
		EventLock: sync.RWMutex{},
	}
	stats := &stats{
		timeLock:   sync.RWMutex{},
		bufferLock: sync.RWMutex{},
	}
	alerts := make(map[string]*int)

	for i, testCase := range testCases {

		eventDisplay.Events = testCase.events
		stats.arrivalTimes = testCase.arrivalTimes
		stats.statBuffers = testCase.statBuffers
		alerts["ttl"] = testCase.alerts["ttl"]
		alerts["dataPadding"] = testCase.alerts["dataPadding"]

		updateEventData(eventDisplay, stats, alerts)

		failed := false

		if testCase.delete {
			if len(deletedEvents) != 1 {
				failed = true
				t.Errorf("Expected %v got %v %d", 1, len(deletedEvents), i)
			}
		}

		if testCase.fill && !failed {

			if testCase.delete {
				if !deletedEvents[0].DataFilled {
					failed = true
					t.Errorf("Expected %v got %v %d", true, eventDisplay.Events[0].DataFilled, i)
				}
			} else if !eventDisplay.Events[0].DataFilled {
				failed = true
				t.Errorf("Expected %v got %v %d", true, eventDisplay.Events[0].DataFilled, i)
			}
		}

		if testCase.update && !failed {

			if testCase.delete {
				if deletedEvents[0].DataTimes[len(deletedEvents[0].DataTimes)-1] != testCase.arrivalTimes[testCase.events[0].Node][len(testCase.arrivalTimes[testCase.events[0].Node])-1] {
					t.Errorf("Expected %v got %v %d", testCase.arrivalTimes[testCase.events[0].Node][len(testCase.arrivalTimes[testCase.events[0].Node])-1], eventDisplay.Events[0].DataTimes[len(eventDisplay.Events[0].DataTimes)-1], i)
				}
			} else if eventDisplay.Events[0].DataTimes[len(eventDisplay.Events[0].DataTimes)-1] != testCase.arrivalTimes[testCase.events[0].Node][len(testCase.arrivalTimes[testCase.events[0].Node])-1] {
				t.Errorf("Expected %v got %v %d", testCase.arrivalTimes[testCase.events[0].Node][len(testCase.arrivalTimes[testCase.events[0].Node])-1], eventDisplay.Events[0].DataTimes[len(eventDisplay.Events[0].DataTimes)-1], i)
			}
		}

		deletedEvents = make([]*widgets.Event, 0)
	}

	reportDelete = ReportDeleteOri
}
