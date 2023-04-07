//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package widgets

import (
	"testing"
	"time"
)

func TestReportText(t *testing.T) {

	curTime, _ := time.Parse("2006-01-02 15:04:05", "2001-01-01 01:01:30")

	testCases := []struct {
		event      *Event
		reportText string
	}{
		{
			event: &Event{
				Node:          "node1",
				Stat:          "stat1",
				EventType:     "Sudden Change",
				NumTimes:      1,
				Threshold:     0.2,
				LastTriggered: curTime,
				ThresholdTime: 1,
				DataTimes: []time.Time{
					curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime, curTime.Add(time.Second), curTime.Add(time.Second * time.Duration(2)),
				},
				DataStart: curTime.Add(-time.Second * time.Duration(2)),
				AlertTimes: []time.Time{
					curTime, curTime.Add(time.Second),
				},
				Data: []float64{
					1, 2, 3, 4, 5,
				},
			},
			reportText: "Node - node1\n" +
				"Stat - stat1\n\n" +
				"Stat changed by more than the threshold limit of 20.00% at 2001-01-01 01:01:30. This change occured over 1 second(s).\n\n" +
				"Data collected from 2001-01-01 01:01:28 to 2001-01-01 01:01:32\n\n" +
				"2001-01-01 01:01:28 - 1.000000\n" +
				"2001-01-01 01:01:29 - 2.000000\n" +
				"2001-01-01 01:01:30 - 3.000000 ALERT\n" +
				"2001-01-01 01:01:31 - 4.000000 ALERT\n" +
				"2001-01-01 01:01:32 - 5.000000\n",
		},
		{
			event: &Event{
				Node:           "node2",
				Stat:           "stat2",
				EventType:      "Sudden Change",
				NumTimes:       3,
				Threshold:      0.2,
				FirstTriggered: curTime.Add(-time.Second),
				LastTriggered:  curTime.Add(time.Second),
				ThresholdTime:  1,
				DataTimes: []time.Time{
					curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime, curTime.Add(time.Second), curTime.Add(time.Second * time.Duration(2)),
				},
				DataStart: curTime.Add(-time.Second * time.Duration(3)),
				AlertTimes: []time.Time{
					curTime.Add(-time.Second), curTime, curTime.Add(time.Second),
				},
				Data: []float64{
					1, 2, 3, 4, 5,
				},
			},
			reportText: "Node - node2\n" +
				"Stat - stat2\n\n" +
				"Stat changed by more than the threshold limit of 20.00% at 2001-01-01 01:01:29. This change occured over 1 second(s).\n" +
				"Similar changes occured 3 times with the last one occuring at 2001-01-01 01:01:31.\n\n" +
				"Data collected from 2001-01-01 01:01:28 to 2001-01-01 01:01:32\n\n" +
				"No data recieved from server before 2001-01-01 01:01:28\n" +
				"2001-01-01 01:01:28 - 1.000000\n" +
				"2001-01-01 01:01:29 - 2.000000 ALERT\n" +
				"2001-01-01 01:01:30 - 3.000000 ALERT\n" +
				"2001-01-01 01:01:31 - 4.000000 ALERT\n" +
				"2001-01-01 01:01:32 - 5.000000\n",
		},
		{
			event: &Event{
				Node:          "node3",
				Stat:          "stat3",
				EventType:     "Above Threshold",
				NumTimes:      1,
				Threshold:     10,
				LastTriggered: curTime,
				DataTimes: []time.Time{
					curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime, curTime.Add(time.Second * time.Duration(2)),
				},
				DataStart: curTime.Add(-time.Second * time.Duration(2)),
				AlertTimes: []time.Time{
					curTime, curTime.Add(time.Second * time.Duration(2)),
				},
				Data: []float64{
					1, 2, 30, 50,
				},
			},
			reportText: "Node - node3\n" +
				"Stat - stat3\n\n" +
				"Stat exceeded threshold limit of 10.000000 at 2001-01-01 01:01:30.\n\n" +
				"Data collected from 2001-01-01 01:01:28 to 2001-01-01 01:01:32\n\n" +
				"2001-01-01 01:01:28 - 1.000000\n" +
				"2001-01-01 01:01:29 - 2.000000\n" +
				"2001-01-01 01:01:30 - 30.000000 ALERT\n" +
				"No data recieved from server between 2001-01-01 01:01:30 and 2001-01-01 01:01:32\n" +
				"2001-01-01 01:01:32 - 50.000000 ALERT\n",
		},
		{
			event: &Event{
				Node:           "node4",
				Stat:           "stat4",
				EventType:      "Above Threshold",
				NumTimes:       2,
				Threshold:      10,
				FirstTriggered: curTime.Add(-time.Second),
				LastTriggered:  curTime.Add(time.Second * time.Duration(2)),
				DataTimes: []time.Time{
					curTime.Add(-time.Second * time.Duration(2)), curTime.Add(-time.Second), curTime, curTime.Add(time.Second * time.Duration(2)),
				},
				DataStart: curTime.Add(-time.Second * time.Duration(3)),
				AlertTimes: []time.Time{
					curTime.Add(-time.Second), curTime.Add(time.Second * time.Duration(2)),
				},
				Data: []float64{
					1, 20, 3, 50,
				},
			},
			reportText: "Node - node4\n" +
				"Stat - stat4\n\n" +
				"Stat exceeded threshold limit of 10.000000 at 2001-01-01 01:01:29.\n" +
				"Similarly, the stat exceeded threshold limit 2 times with the last one occuring at 2001-01-01 01:01:32.\n\n" +
				"Data collected from 2001-01-01 01:01:28 to 2001-01-01 01:01:32\n\n" +
				"No data recieved from server before 2001-01-01 01:01:28\n" +
				"2001-01-01 01:01:28 - 1.000000\n" +
				"2001-01-01 01:01:29 - 20.000000 ALERT\n" +
				"2001-01-01 01:01:30 - 3.000000\n" +
				"No data recieved from server between 2001-01-01 01:01:30 and 2001-01-01 01:01:32\n" +
				"2001-01-01 01:01:32 - 50.000000 ALERT\n",
		},
		{
			event: &Event{
				Node:          "node5",
				Stat:          "stat5",
				EventType:     "Below Threshold",
				NumTimes:      1,
				Threshold:     10,
				LastTriggered: curTime,
				DataTimes: []time.Time{
					curTime,
				},
				DataStart: curTime.Add(-time.Second * time.Duration(2)),
				AlertTimes: []time.Time{
					curTime,
				},
				Data: []float64{
					1,
				},
			},
			reportText: "Node - node5\n" +
				"Stat - stat5\n\n" +
				"Stat dropped below threshold limit of 10.000000 at 2001-01-01 01:01:30.\n\n" +
				"Data collected from 2001-01-01 01:01:30 to 2001-01-01 01:01:30\n\n" +
				"No data recieved from server before 2001-01-01 01:01:30\n" +
				"2001-01-01 01:01:30 - 1.000000 ALERT\n",
		},
		{
			event: &Event{
				Node:           "node6",
				Stat:           "stat6",
				EventType:      "Below Threshold",
				NumTimes:       2,
				Threshold:      10,
				FirstTriggered: curTime.Add(-time.Second),
				LastTriggered:  curTime.Add(time.Second),
				DataTimes: []time.Time{
					curTime.Add(-time.Second), curTime.Add(time.Second),
				},
				DataStart: curTime.Add(-time.Second * time.Duration(3)),
				AlertTimes: []time.Time{
					curTime.Add(-time.Second), curTime.Add(time.Second),
				},
				Data: []float64{
					1, 3,
				},
			},
			reportText: "Node - node6\n" +
				"Stat - stat6\n\n" +
				"Stat dropped below threshold limit of 10.000000 at 2001-01-01 01:01:29.\n" +
				"Similarly, the stat was below the threshold limit 2 times with the last one occuring at 2001-01-01 01:01:31.\n\n" +
				"Data collected from 2001-01-01 01:01:29 to 2001-01-01 01:01:31\n\n" +
				"No data recieved from server before 2001-01-01 01:01:29\n" +
				"2001-01-01 01:01:29 - 1.000000 ALERT\n" +
				"No data recieved from server between 2001-01-01 01:01:29 and 2001-01-01 01:01:31\n" +
				"2001-01-01 01:01:31 - 3.000000 ALERT\n",
		},
	}

	for i, testCase := range testCases {

		testString := ReportText(testCase.event)

		if testString != testCase.reportText {
			t.Errorf("%d Expected \n%v\n got \n%v\n", i, testCase.reportText, testString)
		}
	}
}
