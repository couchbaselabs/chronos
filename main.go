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
	"os"
	"time"

	log "github.com/couchbase/clog"
	"github.com/couchbaselabs/chronos/widgets"
	ui "github.com/gizak/termui/v3"
)

var (
	nodesTable   *widgets.NodesTable
	statsTable   *widgets.StatsTable
	eventDisplay *widgets.EventDisplay
	nodeSelect1  *widgets.GraphNodeSelect
	nodeSelect2  *widgets.GraphNodeSelect
)

// This variable is used to track which table is currently selected
const (
	leftTable   = 1
	middleTable = 2
	rightTable  = 3
)

// This variable is used to track which graph's nodes are currently displayed
const (
	leftGraph  = 1
	rightGraph = 2
)

func main() {

	// Parse all the flags into a config struct
	config := flagsInit()

	// Initialize the loggers
	err := logsInit()
	if err != nil {
		fmt.Println("main: Unable to initialize loggers:", err)
		return
	}

	// Connect to cluster using the go sdk connector
	cluster, err := clusterInit(
		*config.ip, *config.username, *config.password,
	)
	if err != nil {
		log.Fatalf("main: unable to connect to cluster: %v", err)
	}

	// Get a list of search node hostnames
	nodesList, err := nodesListInit(cluster)
	if err != nil {
		log.Fatalf(
			"main: unable to get node configuration from the server: %v", err,
		)
		return
	}

	// If the cluster has no search nodes
	if len(nodesList) == 0 {
		log.Fatalf(
			"main: unable to detect any active search nodes in the cluster",
		)
	}

	// Obtain a list of stats from the config
	statsList := statsListInit(config.stats)

	// Initialize the stats struct with default values
	stats := statsInit(config, nodesList, statsList)

	// Create a manager instance with all the
	// parameters necessary to spawn new polls
	manager := newManager(
		nodesList, statsList, stats, *config.username,
		*config.password, cluster,
	)

	// Start the manager routine
	// This starts all the polls and enters an infinite
	// loop to check the list of search nodes for any
	// changes using the go sdk
	go monitorCluster(manager)

	// Initialize termui
	err2 := ui.Init()
	defer ui.Close()
	if err2 != nil {
		log.Fatalf("main: Unable to initialize ui: %v", err2)
		return
	}

	// Refresh rate ticker for the UI
	updateTicker := time.NewTicker(time.Second).C

	// Channel for UI events
	uiEvents := ui.PollEvents()

	// Initializing all widgets
	statsTable = widgets.NewStatsTable(statsList)
	nodesTable = widgets.NewNodesTable(nodesList)
	lineChart1 := widgets.NewLineGraph()
	lineChart2 := widgets.NewLineGraph()
	nodeSelect1 = widgets.NewGraphNodeSelect()
	nodeSelect2 = widgets.NewGraphNodeSelect()
	eventDisplay = widgets.NewEventDisplay()
	rebalancePopup := widgets.NewRebalancePopup()

	// Starting the routine to check and accept incoming events
	go eventCreateHandler(
		manager.eventChannel, eventDisplay, stats, config.alerts,
	)

	// Starting the routine to track and update event data
	go eventDataHandler(eventDisplay, stats, config.alerts)

	// Initializing the grid
	grid := gridInit(
		nodesTable, statsTable, lineChart1, lineChart2,
		eventDisplay, rebalancePopup,
	)

	tableSelect := leftTable

	var graphNum int

	refreshUI(
		statsTable, nodesTable, lineChart1, lineChart2,
		eventDisplay, rebalancePopup, grid,
	)

	defer func() {
		if r := recover(); r != nil {
			// Add log line here
			fmt.Printf("main panicked!, r:%v\n", r)
			os.Exit(2)
		}
	}()

	// Main event loop
	for {
		select {

		// UI events
		case e := <-uiEvents:
			switch e.ID {

			// Exit out of the program
			case "q", "Q", "<C-c>":
				return

			// Scroll up
			case "<Up>", "k", "K":
				table := getSelectedTable(tableSelect)
				table.ScrollUp()
				refreshUI(
					statsTable, nodesTable, lineChart1, lineChart2,
					eventDisplay, rebalancePopup, grid,
				)

			// Scroll Down
			case "<Down>", "j", "J":
				table := getSelectedTable(tableSelect)
				table.ScrollDown()
				refreshUI(
					statsTable, nodesTable, lineChart1, lineChart2,
					eventDisplay, rebalancePopup, grid,
				)

			// Move to the table to the right
			case "<Right>", "l", "L":
				if tableSelect != rightTable {
					table := getSelectedTable(tableSelect)
					table.ToggleTableSelect()

					tableSelect += 1
					table = getSelectedTable(tableSelect)
					table.ToggleTableSelect()

					refreshUI(
						statsTable, nodesTable, lineChart1, lineChart2,
						eventDisplay, rebalancePopup, grid,
					)
				}

			// Move to the table to the left
			case "<Left>", "h", "H":
				if tableSelect != leftTable {
					table := getSelectedTable(tableSelect)
					table.ToggleTableSelect()

					tableSelect -= 1
					table = getSelectedTable(tableSelect)
					table.ToggleTableSelect()

					refreshUI(
						statsTable, nodesTable, lineChart1, lineChart2,
						eventDisplay, rebalancePopup, grid,
					)
				}

			// Select stat for the left line chart
			case "a", "A":
				if tableSelect == leftTable &&
					statsTable.SelectedRow != statsTable.Stat2 {

					graphNum = leftGraph
					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat1],
						nodeSelect1.Nodes, nodeSelect1.Stat,
					)
					nodeSelect1.GraphSelectInit(
						nodesTable.Rows, nodesTable.Nodes,
						statsTable.Rows[statsTable.Stat1],
					)
					updateGraph(
						statsTable, stats, nodesList,
						lineChart1, graphNum,
					)
					refreshUI(
						statsTable, nodesTable, lineChart1, lineChart2,
						eventDisplay, rebalancePopup, grid,
					)
				}

			// Select stat for the right line chart
			case "d", "D":
				if tableSelect == leftTable &&
					statsTable.SelectedRow != statsTable.Stat1 {

					graphNum = rightGraph
					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat2],
						nodeSelect2.Nodes, nodeSelect2.Stat,
					)
					nodeSelect2.GraphSelectInit(
						nodesTable.Rows, nodesTable.Nodes,
						statsTable.Rows[statsTable.Stat2],
					)
					updateGraph(
						statsTable, stats, nodesList,
						lineChart2, graphNum,
					)
					refreshUI(
						statsTable, nodesTable, lineChart1, lineChart2,
						eventDisplay, rebalancePopup, grid,
					)
				}

			// Toggle selection of a node
			case "s", "S":
				if tableSelect == middleTable {
					nodesTable.SelectNode()
					nodeSelect, statNum := getSelectedGraphInfo(graphNum)
					nodeSelect.GraphSelectInit(nodesTable.Rows, nodesTable.Nodes, statsTable.Rows[statNum])
					updateUI(
						nodesList, stats, lineChart1, lineChart2,
						nodeSelect1, nodeSelect2,
					)
					refreshUI(
						statsTable, nodesTable, lineChart1, lineChart2,
						eventDisplay, rebalancePopup, grid,
					)
				}

			// Print a report for an alert
			case "p", "P":
				if tableSelect == middleTable {
					eventDisplay.ReportEvent(*config.reportPath)
				}

			// Triggered on resize of the terminal window
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				rebalancePopup.Resize(payload.Width, payload.Height)
				refreshUI(
					statsTable, nodesTable, lineChart1, lineChart2,
					eventDisplay, rebalancePopup, grid,
				)
			}

		// Update line charts and re-render UI
		case <-updateTicker:
			updateUI(
				nodesList, stats, lineChart1, lineChart2,
				nodeSelect1, nodeSelect2,
			)
			refreshUI(
				statsTable, nodesTable, lineChart1, lineChart2,
				eventDisplay, rebalancePopup, grid,
			)

		// Recieve errors from the other routines
		case errorMsg := <-manager.errChannel:
			if errorMsg.terminate {
				log.Fatalf(errorMsg.description)
				return
			} else {
				log.Warnf(errorMsg.description)
			}

		// Recieve rebalance notifications from the polling routines
		case <-manager.rebalanceChannel:
			log.Warnf("Cluster undergoing rebalance")
			rebalancePopup.SetRebalance()
			refreshUI(
				statsTable, nodesTable, lineChart1, lineChart2,
				eventDisplay, rebalancePopup, grid,
			)

		// Recieve node addition or removal notification from manager
		case info := <-manager.updateChannel:
			if info.add {
				nodesTable.AddNode(info.node)
				nodeSelect1.AddNode(info.node)
				nodeSelect2.AddNode(info.node)
			} else {
				nodesTable.RemoveNode(info.node)
				nodeSelect1.RemoveNode(info.node)
				nodeSelect2.RemoveNode(info.node)
			}

			nodesList = nodesTable.Rows
			updateGraphNodes(
				nodesList, stats, lineChart1, lineChart2,
				nodeSelect1, nodeSelect2,
			)
		}
	}
}

// Re-render UI without updating line charts
func refreshUI(statsTable *widgets.StatsTable, nodesTable *widgets.NodesTable,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph,
	eventDisplay *widgets.EventDisplay, rebalancePopup *widgets.RebalancePopup,
	grid *ui.Grid) {

	ui.Render(grid)
	ui.Render(statsTable)
	ui.Render(nodesTable)
	ui.Render(lineChart1)
	ui.Render(lineChart2)
	ui.Render(eventDisplay)

	rebalancePopup.DisplayLock.RLock()
	if rebalancePopup.Display {
		ui.Render(rebalancePopup)
	}
	rebalancePopup.DisplayLock.RUnlock()
}

func getSelectedTable(tableSelect int) widgets.Table {
	if tableSelect == leftTable {
		return statsTable
	} else if tableSelect == middleTable {
		return nodesTable
	} else {
		return eventDisplay
	}
}

func getSelectedGraphInfo(graphNum int) (*widgets.GraphNodeSelect, int) {
	if graphNum == leftGraph {
		return nodeSelect1, statsTable.Stat1
	} else {
		return nodeSelect2, statsTable.Stat2
	}
}
