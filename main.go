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
	lineChart1   *widgets.LineGraph
	lineChart2   *widgets.LineGraph
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
	lineChart1 = widgets.NewLineGraph(nodesList, 1)
	lineChart2 = widgets.NewLineGraph(nodesList, 2)
	eventDisplay = widgets.NewEventDisplay()
	popupManager := widgets.NewPopupManager()

	// Starting the routine to check and accept incoming events
	go eventCreateHandler(
		manager.eventChannel, eventDisplay, stats, config.alerts,
	)

	// Starting the routine to track and update event data
	go eventDataHandler(eventDisplay, stats, config.alerts)

	// Grid is used to hold all the UI widgets together and maintain relative sizes
	grid := gridInit(
		nodesTable, statsTable, lineChart1, lineChart2,
		eventDisplay, popupManager,
	)

	// Used to determine the table currently in use
	tableSelect := leftTable

	// Used to determine the graph corresponding to the nodesTable
	var graphNum int

	refreshUI(
		statsTable, nodesTable, lineChart1, lineChart2,
		eventDisplay, popupManager, grid,
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
				ui.Render(table)
				popupManager.Render()

			// Scroll Down
			case "<Down>", "j", "J":
				table := getSelectedTable(tableSelect)
				table.ScrollDown()
				ui.Render(table)
				popupManager.Render()

			// Move to the table to the right
			case "<Right>", "l", "L":
				if tableSelect != rightTable {
					table := getSelectedTable(tableSelect)
					table.ToggleTableSelect()
					ui.Render(table)

					tableSelect += 1
					table = getSelectedTable(tableSelect)
					table.ToggleTableSelect()
					ui.Render(table)
					popupManager.Render()
				}
			// Move to the table to the left
			case "<Left>", "h", "H":
				if tableSelect != leftTable {
					table := getSelectedTable(tableSelect)
					table.ToggleTableSelect()
					ui.Render(table)

					tableSelect -= 1
					table = getSelectedTable(tableSelect)
					table.ToggleTableSelect()
					ui.Render(table)

					popupManager.Render()
				}
			// Select stat for the left line chart
			case "a", "A":
				if tableSelect == leftTable &&
					statsTable.SelectedRow != statsTable.Stat2 {

					graphNum = leftGraph
					lineChart1.SelectGraph()
					lineChart2.UnSelectGraph()

					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat1],
						lineChart1.Nodes, lineChart1.Stat,
					)

					if lineChart1.Stat != statsTable.Rows[statsTable.Stat1] {
						updateGraph(
							statsTable.Rows[statsTable.Stat1], stats, nodesList,
							lineChart1, graphNum,
						)
					}
					ui.Render(lineChart1)
					ui.Render(lineChart2)
					ui.Render(statsTable)
					ui.Render(nodesTable)
					popupManager.Render()
				}

			// Select stat for the right line chart
			case "d", "D":
				if tableSelect == leftTable &&
					statsTable.SelectedRow != statsTable.Stat1 {

					graphNum = rightGraph
					lineChart2.SelectGraph()
					lineChart1.UnSelectGraph()

					statsTable.SelectGraph(graphNum)
					nodesTable.SelectStat(
						statsTable.Rows[statsTable.Stat2],
						lineChart2.Nodes, lineChart2.Stat,
					)

					if lineChart2.Stat != statsTable.Rows[statsTable.Stat2] {
						updateGraph(
							statsTable.Rows[statsTable.Stat2], stats, nodesList,
							lineChart2, graphNum,
						)
					}
					ui.Render(lineChart1)
					ui.Render(lineChart2)
					ui.Render(statsTable)
					ui.Render(nodesTable)
					popupManager.Render()
				}
			// Interact with a table
			case "<Enter>":
				switch tableSelect {
				case middleTable:
					nodesTable.SelectNode()
					lineChart := getSelectedGraph(graphNum)
					lineChart.SelectNode(nodesTable.Rows[nodesTable.SelectedRow])
					updateUI(stats, lineChart)
					ui.Render(nodesTable)
					ui.Render(lineChart)
					popupManager.Render()
				case rightTable:
					eventDisplay.ReportEvent(*config.reportPath)
				}
			// Toggle legend for the selected graph
			case "p", "P":
				lineChart := getSelectedGraph(graphNum)
				lineChart.ToggleLegend()
				ui.Render(lineChart)
			// Scroll down on the hovered table
			case "<MouseWheelDown>":

				payload := e.Payload.(ui.Mouse)
				x := payload.X
				y := payload.Y

				var onTable bool

				prevTableSelect := tableSelect

				if tableSelect, onTable = hoveredTable(x, y, tableSelect); onTable {

					if prevTableSelect == tableSelect {
						table := getSelectedTable(tableSelect)
						table.ScrollDown()
						ui.Render(table)
					} else {
						table := getSelectedTable(prevTableSelect)
						table.ToggleTableSelect()
						ui.Render(table)

						table = getSelectedTable(tableSelect)
						table.ToggleTableSelect()
						table.ScrollDown()
						ui.Render(table)
					}

					popupManager.Render()
				}
			// Scroll up on the hovered table
			case "<MouseWheelUp>":

				payload := e.Payload.(ui.Mouse)
				x := payload.X
				y := payload.Y

				var onTable bool

				prevTableSelect := tableSelect

				if tableSelect, onTable = hoveredTable(x, y, tableSelect); onTable {

					if prevTableSelect == tableSelect {
						table := getSelectedTable(tableSelect)
						table.ScrollUp()
						ui.Render(table)
					} else {
						table := getSelectedTable(prevTableSelect)
						table.ToggleTableSelect()
						ui.Render(table)

						table = getSelectedTable(tableSelect)
						table.ToggleTableSelect()
						table.ScrollUp()
						ui.Render(table)
					}

					popupManager.Render()
				}
			// Select the clicked line in any table
			case "<MouseLeft>":

				payload := e.Payload.(ui.Mouse)
				x := payload.X
				y := payload.Y

				var onTable bool

				prevTableSelect := tableSelect

				if tableSelect, onTable = hoveredTable(x, y, tableSelect); onTable {

					if prevTableSelect == tableSelect {
						table := getSelectedTable(tableSelect)
						table.HandleClick(x, y)
						ui.Render(table)
					} else {
						table := getSelectedTable(prevTableSelect)
						table.ToggleTableSelect()
						ui.Render(table)

						table = getSelectedTable(tableSelect)
						table.ToggleTableSelect()
						table.HandleClick(x, y)
						ui.Render(table)
					}

					popupManager.Render()
				}
			// Event called when terminal window is resized
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				popupManager.SetSize(payload.Width, payload.Height)
				refreshUI(
					statsTable, nodesTable, lineChart1, lineChart2,
					eventDisplay, popupManager, grid,
				)
			case "t", "T": // To simulate a rebalance
				manager.popupChannel <- "rebalance"
			}

		// Update line charts and re-render UI
		case <-updateTicker:
			updateUI(stats, lineChart1)
			updateUI(stats, lineChart2)
			refreshUI(
				statsTable, nodesTable, lineChart1, lineChart2,
				eventDisplay, popupManager, grid,
			)

		// Handle errors from the update_stats routine
		case errorMsg := <-manager.errChannel:
			if errorMsg.terminate {
				log.Fatalf(errorMsg.description)
				return
			} else {
				log.Warnf(errorMsg.description)
			}
		// Handle warnings from the update_stats routine
		case msg := <-manager.popupChannel:

			if msg == "rebalance" {
				popupManager.NewPopup(
					"Cluster undergoing rebalance", "rebalance",
					time.Now().Add(time.Millisecond*time.Duration(1500)),
				)
				log.Warnf("Cluster undergoing rebalance")
			} else {
				popupManager.NewPopup(
					"Slow response from "+msg, "warning",
					time.Now().Add(time.Millisecond*time.Duration(3000)),
				)
				log.Warnf("Slow response from " + msg)
			}

			refreshUI(
				statsTable, nodesTable, lineChart1, lineChart2,
				eventDisplay, popupManager, grid,
			)
		// Handle node changes from the manager routine
		case info := <-manager.updateChannel:
			if info.add {
				nodesTable.AddNode(info.node)
				lineChart1.AddNode(info.node)
				lineChart2.AddNode(info.node)
				popupManager.AddNodePopup(info.node)
			} else {
				nodesTable.RemoveNode(info.node)
				lineChart1.RemoveNode(info.node)
				lineChart2.RemoveNode(info.node)
				popupManager.RemoveNodePopup(info.node)
			}

			nodesList = nodesTable.Rows
			updateUI(stats, lineChart1)
			updateUI(stats, lineChart2)
		}
	}
}

// Re-render UI without updating line charts
func refreshUI(statsTable *widgets.StatsTable, nodesTable *widgets.NodesTable,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph,
	eventDisplay *widgets.EventDisplay, popupManager *widgets.PopupManager,
	grid *ui.Grid) {

	ui.Render(grid)
	ui.Render(statsTable)
	ui.Render(nodesTable)
	ui.Render(lineChart1)
	ui.Render(lineChart2)
	ui.Render(eventDisplay)

	popupManager.Render()
}

// Returns the table currently selected
func getSelectedTable(tableSelect int) widgets.Table {
	if tableSelect == leftTable {
		return statsTable
	} else if tableSelect == middleTable {
		return nodesTable
	} else {
		return eventDisplay
	}
}

// Returns the table the mouse is currently on
func hoveredTable(x int, y int, tableSelect int) (int, bool) {

	if statsTable.Contains(x, y) {
		return leftTable, true
	} else if nodesTable.Contains(x, y) {
		return middleTable, true
	} else if eventDisplay.Contains(x, y) {
		return rightTable, true
	} else {
		return tableSelect, false
	}
}

// Returns the graph the nodes table is currently linked to
func getSelectedGraph(graphNum int) *widgets.LineGraph {
	if graphNum == leftGraph {
		return lineChart1
	} else {
		return lineChart2
	}
}
