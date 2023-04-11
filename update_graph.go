//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"github.com/couchbaselabs/chronos/widgets"
)

// Handle selection of a stat for a particular graph stat. Copies stat buffers
// for that particular stat and nodes to the widget
func updateGraph(statName string, stats *stats, nodesList []string,
	lineChart *widgets.LineGraph, graphNum int) {

	lineChart.Stat = statName

	stats.bufferLock.RLock()
	for _, nodeData := range lineChart.Nodes {
		nodeData.Line = stats.statBuffers[nodeData.Node][statName]
		nodeData.Active = true
	}
	stats.bufferLock.RUnlock()
}

// Handle changes in nodes selected
func updateUI(stats *stats, lineChart *widgets.LineGraph) {

	stats.bufferLock.RLock()
	if lineChart.Stat != "" {
		for _, nodeData := range lineChart.Nodes {
			if nodeData.Active {
				nodeData.Line = stats.statBuffers[nodeData.Node][lineChart.Stat]
			} else {
				nodeData.Line = make([]float64, 300)
			}
		}
	}
	stats.bufferLock.RUnlock()
}
