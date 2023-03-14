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
func updateGraph(table *widgets.StatsTable, stats *stats, nodesList []string,
	lineChart *widgets.LineGraph, graphNum int) {

	var graphStat string

	if graphNum == 1 {
		graphStat = table.Rows[table.Stat1]
	} else {
		graphStat = table.Rows[table.Stat2]
	}

	lineChart.Title = graphStat
	lineChart.Data = make([][]float64, len(nodesList))

	stats.bufferLock.RLock()
	for i, node := range nodesList {
		lineChart.Data[i] = stats.statBuffers[node][graphStat]
	}
	stats.bufferLock.RUnlock()
}

// Handle changes in the list of nodes in a cluster
func updateGraphNodes(nodesList []string, stats *stats,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph,
	nodeSelect1 *widgets.GraphNodeSelect, nodeSelect2 *widgets.GraphNodeSelect) {

	stats.bufferLock.RLock()
	if lineChart1.Title != "" {
		lineChart1.Data = make([][]float64, len(nodesList))
		for i, node := range nodesList {
			if nodeSelect1.Nodes[node] {
				lineChart1.Data[i] = stats.statBuffers[node][lineChart1.Title]
			} else {
				lineChart1.Data[i] = make([]float64, 0)
			}
		}
	}
	if lineChart2.Title != "" {
		for i, node := range nodesList {
			lineChart2.Data = make([][]float64, len(nodesList))
			if nodeSelect2.Nodes[node] {
				lineChart2.Data[i] = stats.statBuffers[node][lineChart2.Title]
			} else {
				lineChart2.Data[i] = make([]float64, 0)
			}
		}
	}
	stats.bufferLock.RUnlock()
}

// Update graph if node has been toggled in nodesTable
func updateUI(nodesList []string,
	stats *stats, lineChart1 *widgets.LineGraph,
	lineChart2 *widgets.LineGraph, nodeSelect1 *widgets.GraphNodeSelect,
	nodeSelect2 *widgets.GraphNodeSelect) {

	stats.bufferLock.RLock()
	if lineChart1.Title != "" {
		for i, node := range nodesList {
			if nodeSelect1.Nodes[node] {
				lineChart1.Data[i] = stats.statBuffers[node][lineChart1.Title]
			} else {
				lineChart1.Data[i] = make([]float64, 0)
			}
		}
	}
	if lineChart2.Title != "" {
		for i, node := range nodesList {
			if nodeSelect2.Nodes[node] {
				lineChart2.Data[i] = stats.statBuffers[node][lineChart2.Title]
			} else {
				lineChart2.Data[i] = make([]float64, 0)
			}
		}
	}
	stats.bufferLock.RUnlock()
}
