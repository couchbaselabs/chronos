//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/couchbaselabs/chronos/widgets"
)

// Structure to hold all the necessary parameters to monitor the cluster for
// nodes being added or removed and create polls for them appropriately
type manager struct {
	stats *stats

	// Holds the kill switch channels to all the polling routines
	nodes map[string]chan bool

	// To track of nodes currently being polled
	polledNodes map[string]bool

	errChannel   chan *errorMsg
	eventChannel chan *widgets.Event
	popupChannel chan string

	// Used to signal the addition or removal of any node or stat to the main routine
	updateChannel chan updateMessage

	username string
	password string
	cluster  *gocb.Cluster
}

// Struct to indicate the addition or removal of any node or stat
type updateMessage struct {
	add  bool
	node string
	stat string
}

// Create and initialize a new manager
func newManager(nodesList []string, stats *stats,
	username string, password string, cluster *gocb.Cluster) *manager {

	manager := &manager{
		stats:         stats,
		nodes:         make(map[string]chan bool),
		polledNodes:   make(map[string]bool),
		errChannel:    make(chan *errorMsg),
		eventChannel:  make(chan *widgets.Event),
		popupChannel:  make(chan string),
		updateChannel: make(chan updateMessage),
		username:      username,
		password:      password,
		cluster:       cluster,
	}

	for _, node := range nodesList {
		manager.nodes[node] = make(chan bool)
	}

	return manager
}

// Manages all the polling routines, starting new ones and closing old ones
// in response to changes in the search nodes of the cluster
func monitorCluster(manager *manager) {

	// Frequency of checking the server for cluster changes
	monitorTicker := time.NewTicker(time.Second).C

	// Start the polling routines for the first time
	startPolls(manager)

	// Main loop of the routine
	for {

		<-monitorTicker

		// Set all existing nodes to false
		for node := range manager.polledNodes {
			manager.polledNodes[node] = false
		}

		// Get new set of nodes from the cluster
		nodes, err := nodesListInit(manager.cluster)
		if err != nil {
			manager.errChannel <- newErrorMsg(
				err,
				"manager: unable to get node configuration from the server: "+
					err.Error(), true,
			)
			return
		}
		if len(nodes) == 0 {
			manager.errChannel <- newErrorMsg(
				err,
				"manager: unable to detect any active search nodes in the cluster",
				true,
			)
			return
		}

		// Set nodes in both lists to true
		// Create new polls for nodes not in the current list
		for _, node := range nodes {
			if _, ok := manager.polledNodes[node]; ok {
				manager.polledNodes[node] = true
			} else {
				createPoll(manager, node)
				manager.updateChannel <- updateMessage{
					add:  true,
					node: node,
					stat: "",
				}
			}
		}

		// Check and delete polls for nodes no longer in current list
		for node, status := range manager.polledNodes {
			if !status {
				deletePoll(manager, node)
				manager.updateChannel <- updateMessage{
					add:  false,
					node: node,
					stat: "",
				}
			}
		}
	}
}

// Start polls for all the nodes
func startPolls(manager *manager) {

	for node, killSwitch := range manager.nodes {
		updateStatsParams := newUpdateStatsParams(
			manager.username, manager.password, manager.stats, node,
			manager.errChannel, manager.eventChannel,
			manager.popupChannel, manager.updateChannel, killSwitch,
		)

		go updateStatsExponentialBackoff(updateStatsParams)

		manager.polledNodes[node] = true
	}
}

// Start poll for a new node while creating slices to hold its information
func createPoll(manager *manager, node string) {

	manager.polledNodes[node] = true
	manager.nodes[node] = make(chan bool)

	statsList := getStatsList(manager.stats)
	manager.stats.bufferLock.Lock()
	manager.stats.statBuffers[node] = make(map[string][]float64)
	for _, stat := range statsList {
		manager.stats.statBuffers[node][stat] = make([]float64, 300)
	}
	manager.stats.bufferLock.Unlock()

	manager.stats.timeLock.Lock()
	manager.stats.arrivalTimes[node] = make([]time.Time, 300)
	manager.stats.timeLock.Unlock()

	updateStatsParams := newUpdateStatsParams(
		manager.username, manager.password, manager.stats, node,
		manager.errChannel, manager.eventChannel,
		manager.popupChannel, manager.updateChannel, manager.nodes[node],
	)

	go updateStatsExponentialBackoff(updateStatsParams)
}

// Remove an exisitng poll and delete all the data related to that node
func deletePoll(manager *manager, node string) {

	manager.nodes[node] <- true
	delete(manager.nodes, node)
	delete(manager.polledNodes, node)

	manager.stats.bufferLock.Lock()

	for stat := range manager.stats.statBuffers[node] {
		delete(manager.stats.statBuffers[node], stat)
	}
	delete(manager.stats.statBuffers, node)

	manager.stats.bufferLock.Unlock()

	manager.stats.timeLock.Lock()
	delete(manager.stats.arrivalTimes, node)
	manager.stats.timeLock.Unlock()
}

// Copies the stats list to reduce amount of time each routine holds the lock
func getStatsList(stats *stats) []string {

	stats.statsListLock.RLock()
	statsList := make([]string, len(stats.statsList))
	copy(statsList, stats.statsList)
	stats.statsListLock.RUnlock()

	return statsList
}
