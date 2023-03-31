//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"log"

	"github.com/couchbaselabs/chronos/widgets"

	ui "github.com/gizak/termui/v3"
)

func updateGraph(table *widgets.StatsTable, stats map[string]map[string][]float64,
	statNodes map[string][]string, lineChart *widgets.LineGraph, graphNum int) {

	var graphStat string

	if graphNum == 1 {
		graphStat = table.Rows[table.Stat1]
	} else {
		graphStat = table.Rows[table.Stat2]
	}

	lineChart.Title = graphStat
	lineChart.Data = make([][]float64, len(statNodes[graphStat]))

	for i, node := range statNodes[graphStat] {
		lineChart.Data[i] = stats[node][graphStat]
	}
}

func updateUI(loggerRender *log.Logger, statNodes map[string][]string,
	stats map[string]map[string][]float64, lineChart1 *widgets.LineGraph,
	lineChart2 *widgets.LineGraph, nodeSelect1 *widgets.GraphNodeSelect,
	nodeSelect2 *widgets.GraphNodeSelect) {

	if lineChart1.Title != "" {
		for i, node := range statNodes[lineChart1.Title] {
			if nodeSelect1.Nodes[node] {
				lineChart1.Data[i] = stats[node][lineChart1.Title]
			} else {
				lineChart1.Data[i] = make([]float64, 0)
			}
		}
	}
	if lineChart2.Title != "" {
		for i, node := range statNodes[lineChart2.Title] {
			if nodeSelect2.Nodes[node] {
				lineChart2.Data[i] = stats[node][lineChart2.Title]
			} else {
				lineChart2.Data[i] = make([]float64, 0)
			}
		}
	}
	ui.Render(lineChart1)
	ui.Render(lineChart2)
}
