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
	"time"

	"github.com/chronos/widgets"

	ui "github.com/gizak/termui/v3"
)

func main() {

	configPath := flagsInit()

	loggerRender, loggerUpdate, err := logsInit()
	if err != nil {
		fmt.Println("main: Unable to initialize loggers:", err)
		return
	}

	stats := make(map[string]map[string][]float64)

	config, err := configInit(configPath, loggerRender)
	if err != nil {
		loggerRender.Fatalf(
			"main: Unable to initialize configuration file:", err,
		)
		return
	}

	statNodes := statNodeInit(config)

	errChannel := make(chan string)
	for nodeName, statsList := range config.Nodes {
		go UpdateStats(
			stats, config, nodeName, statsList, loggerUpdate, errChannel,
		)
	}

	err2 := uiInit()
	defer ui.Close()
	if err2 != nil {
		loggerRender.Fatalf(
			"main: Unable to initialize ui:", err2,
		)
	}

	ticker := time.NewTicker(time.Second).C
	uiEvents := ui.PollEvents()
	statsTable := widgets.NewStatsTable()
	nodesTable := widgets.NewNodesTable()
	lineChart1 := widgets.NewLineGraph()
	lineChart2 := widgets.NewLineGraph()
	nodeSelect1 := widgets.NewGraphNodeSelect()
	nodeSelect2 := widgets.NewGraphNodeSelect()

	for _, statsList := range config.Nodes {
		for _, stat := range statsList {
			contains := false

			for _, statInTable := range statsTable.Rows {
				if stat == statInTable {
					contains = true
					break
				}
			}

			if !contains {
				statsTable.Rows = append(statsTable.Rows, stat)
			}
		}
	}

	grid := gridInit(nodesTable, statsTable, lineChart1, lineChart2)

	tableSelect := 2
	var graphNum int

	ui.Render(grid)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "Q", "<C-c>":
				ui.Close()
				return
			case "<Up>":
				if tableSelect == 1 {
					nodesTable.ScrollUp()
					ui.Render(grid)
				} else {
					statsTable.ScrollUp()
					ui.Render(grid)
				}
			case "<Down>":
				if tableSelect == 1 {
					nodesTable.ScrollDown()
					ui.Render(grid)
				} else {
					statsTable.ScrollDown()
					ui.Render(grid)
				}
			case "<Right>":
				if tableSelect == 2 {
					tableSelect = 1
					nodesTable.ToggleTableSelect()
					statsTable.ToggleTableSelect()
				}
				ui.Render(grid)
			case "<Left>":
				if tableSelect == 1 {
					tableSelect = 2
					nodesTable.ToggleTableSelect()
					statsTable.ToggleTableSelect()
				}
				ui.Render(grid)
			case "a", "A":
				if tableSelect == 2 &&
					statsTable.SelectedRow != statsTable.Stat2 {

					graphNum = 1
					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat1],
						statNodes,
						nodeSelect1.Nodes,
						nodeSelect1.Stat,
					)
					nodeSelect1.GraphSelectInit(
						nodesTable.Rows,
						nodesTable.Nodes,
						statsTable.Rows[statsTable.Stat1],
					)
					updateGraph(
						statsTable,
						stats,
						statNodes,
						lineChart1,
						graphNum,
					)
					ui.Render(grid)
				}
			case "d", "D":
				if tableSelect == 2 &&
					statsTable.SelectedRow != statsTable.Stat1 {

					graphNum = 2
					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat2],
						statNodes,
						nodeSelect2.Nodes,
						nodeSelect2.Stat,
					)
					nodeSelect2.GraphSelectInit(
						nodesTable.Rows,
						nodesTable.Nodes,
						statsTable.Rows[statsTable.Stat2],
					)
					updateGraph(
						statsTable,
						stats,
						statNodes,
						lineChart2,
						graphNum,
					)
					ui.Render(grid)
				}
			case "s", "S":
				if tableSelect == 1 {
					nodesTable.SelectNode()
					if graphNum == 1 {
						nodeSelect1.GraphSelectInit(
							nodesTable.Rows,
							nodesTable.Nodes,
							statsTable.Rows[statsTable.Stat1],
						)
					} else {
						nodeSelect2.GraphSelectInit(
							nodesTable.Rows,
							nodesTable.Nodes,
							statsTable.Rows[statsTable.Stat2],
						)
					}
					ui.Render(grid)
					updateUI(
						loggerRender,
						statNodes,
						stats,
						lineChart1,
						lineChart2,
						nodeSelect1,
						nodeSelect2,
					)
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Render(grid)
				ui.Render(statsTable)
				ui.Render(nodesTable)
			}
		case <-ticker:
			updateUI(
				loggerRender,
				statNodes,
				stats,
				lineChart1,
				lineChart2,
				nodeSelect1,
				nodeSelect2,
			)
		case errorString := <-errChannel:
			ui.Close()
			loggerUpdate.Fatalf(errorString)
		}
	}
}
