//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/couchbaselabs/chronos/widgets"

	log "github.com/couchbase/clog"
	"github.com/couchbase/gocb/v2"
	ui "github.com/gizak/termui/v3"
)

// Holds all the input from the user given as command line arguments
type config struct {
	username   *string
	password   *string
	ip         *string
	reportPath *string
	stats      map[string]*configStatInfo
	alerts     map[string]*int
}

// Holds all alert related thresholds for a particular stat
type configStatInfo struct {
	MinVal        *float64
	MaxVal        *float64
	MaxChange     *float64
	MaxChangeTime *int
}

// Holds all incoming stat data from the server
type stats struct {

	// The main map holding the stat data (map[node][stat][300][float64])
	statBuffers map[string]map[string][]float64

	// List of stats to monitor
	statsList []string

	// Arrival time of messages from each node (map[node][300][time.Time])
	arrivalTimes map[string][]time.Time

	statInfo map[string]*configStatInfo

	// Lock for statBuffers
	bufferLock sync.RWMutex

	// Lock for arrivalTimes
	timeLock sync.RWMutex
}

// Define and parse flags
func flagsInit() *config {

	config := &config{}

	config.username = flag.String(
		"username", "Administrator", "Provide the username for the cluster",
	)
	config.password = flag.String(
		"password", "123456", "Provide the password for the cluster",
	)
	config.ip = flag.String(
		"connection_string", "127.0.0.1:12000",
		"Provide the ip address for one of the search nodes",
	)
	config.reportPath = flag.String(
		"report", "./", "Provide path to print reports",
	)
	config.stats = make(map[string]*configStatInfo)
	config.alerts = make(map[string]*int)

	statsList := []string{
		"batch_bytes_added",
		"batch_bytes_removed",
		"curr_batches_blocked_by_herder",
		"num_batches_introduced",
		"num_bytes_used_ram",
		"num_gocbcore_dcp_agents",
		"num_gocbcore_stats_agents",
		"pct_cpu_gc",
		"tot_batches_merged",
		"tot_batches_new",
		"tot_bleve_dest_closed",
		"tot_bleve_dest_opened",
		"tot_queryreject_on_memquota",
		"tot_rollback_full",
		"tot_rollback_partial",
		"total_gc",
		"total_queries_rejected_by_herder",
		"utilization:billableUnitsRate",
		"utilization:cpuPercent",
		"utilization:diskBytes",
		"utilization:memoryBytes",
	}

	for _, stat := range statsList {

		configStatInfo := &configStatInfo{}

		configStatInfo.MinVal = flag.Float64(
			stat+"_min_val", math.NaN(),
			"Provide the minimum threshold value for "+stat,
		)
		configStatInfo.MaxVal = flag.Float64(
			stat+"_max_val", math.NaN(),
			"Provide the maximum threshold value for "+stat,
		)
		configStatInfo.MaxChange = flag.Float64(
			stat+"_max_change", math.NaN(),
			"Provide the maximum change permitted for "+stat,
		)
		configStatInfo.MaxChangeTime = flag.Int(
			stat+"_max_change_time", 1,
			"Provide the amount of time for the max change "+stat,
		)

		config.stats[stat] = configStatInfo
	}

	config.alerts["ttl"] = flag.Int(
		"alert_TTL", 120, "Provide number of seconds an alert should live",
	)
	config.alerts["dataPadding"] = flag.Int(
		"alert_data_padding", 20,
		"Provide number of seconds of data before and after an alert",
	)

	flag.Parse()

	// Check to verify alert parameters are within bounds
	checkAlertParams(config.alerts)
	return config
}

// Validating given alert time to live (TTL) and data padding
// and setting default or max values if out of bounds
func checkAlertParams(alerts map[string]*int) {
	defaultTTL := 120
	defaultDataPadding := 20
	maxTTL := 600
	maxDataPadding := 60

	if val, ok := alerts["ttl"]; !ok {
		alerts["ttl"] = &defaultTTL
	} else if *val <= 0 {
		alerts["ttl"] = &defaultTTL
	} else if *val > maxTTL {
		alerts["ttl"] = &maxTTL
	}

	if val, ok := alerts["dataPadding"]; !ok {
		alerts["dataPadding"] = &defaultDataPadding
	} else if *val <= 0 {
		alerts["dataPadding"] = &defaultDataPadding
	} else if *val > maxDataPadding {
		alerts["dataPadding"] = &maxDataPadding
	}
}

// Initializing the logger
func logsInit() error {

	logFile, err := os.OpenFile(
		"chronos.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	log.SetOutput(logFile)
	log.SetLoggerCallback(loggerFunc)

	return nil
}

// Setting log format
func loggerFunc(level, format string, args ...interface{}) string {

	ts := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	prefix := ts + " [" + level + "] "
	if format != "" {
		return prefix + fmt.Sprintf(format, args...)
	}
	return prefix + fmt.Sprint(args...)
}

// Starting a connection to the server
func clusterInit(connectionString string, username string,
	password string) (*gocb.Cluster, error) {

	cluster, err := gocb.Connect(
		"couchbase://"+connectionString, gocb.ClusterOptions{
			Authenticator: gocb.PasswordAuthenticator{
				Username: username,
				Password: password,
			},
		})

	if err != nil {
		return cluster, err
	}

	return cluster, nil
}

// Getting a list of existing search nodes from the cluster
func nodesListInit(cluster *gocb.Cluster) ([]string, error) {

	pings, err := cluster.Ping(&gocb.PingOptions{
		ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeSearch},
	})
	if err != nil {
		return nil, err
	}

	nodesList := make([]string, 0)

	for service, pingReports := range pings.Services {
		if service == gocb.ServiceTypeSearch {
			for _, pingReport := range pingReports {
				if pingReport.State == gocb.PingStateOk {
					nodesList = append(nodesList, pingReport.Remote)
				}
			}
		}
	}

	return nodesList, nil
}

// Getting the list of stats from the config
func statsListInit(stats map[string]*configStatInfo) []string {

	statsList := make([]string, 0)

	for stat := range stats {
		statsList = append(statsList, stat)
	}

	return statsList
}

// Initialize the stats struct with empty slices
func statsInit(config *config, nodesList []string, statsList []string) *stats {

	statBuffers := make(map[string]map[string][]float64)
	arrivalTimes := make(map[string][]time.Time)

	for _, node := range nodesList {

		statBuffers[node] = make(map[string][]float64)
		arrivalTimes[node] = make([]time.Time, 110)

		for _, stat := range statsList {
			statBuffers[node][stat] = make([]float64, 110)
		}
	}

	return &stats{
		statBuffers:  statBuffers,
		statsList:    statsList,
		statInfo:     config.stats,
		arrivalTimes: arrivalTimes,
		bufferLock:   sync.RWMutex{},
		timeLock:     sync.RWMutex{},
	}
}

// Setting up the widgets in a grid with relative ratios and positions
func gridInit(nodesTable *widgets.NodesTable, statsTable *widgets.StatsTable,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph,
	eventDisplay *widgets.EventDisplay,
	rebalancePopup *widgets.RebalancePopup) *ui.Grid {

	lineChart1.Title = ""
	lineChart2.Title = ""
	lineChart1.Data = make([][]float64, 0)
	lineChart2.Data = make([][]float64, 0)

	grid := ui.NewGrid()

	grid.Set(
		ui.NewRow(3.0/5,
			ui.NewCol(1.0/2, lineChart1),
			ui.NewCol(1.0/2, lineChart2),
		),
		ui.NewRow(2.0/5,
			ui.NewCol(1.0/4, statsTable),
			ui.NewCol(1.0/4, nodesTable),
			ui.NewCol(1.0/2, eventDisplay),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	rebalancePopup.Resize(termWidth, termHeight)
	return grid
}
