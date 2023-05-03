//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	flag "github.com/couchbaselabs/chronos/cflag"

	"github.com/couchbaselabs/chronos/widgets"

	log "github.com/couchbase/clog"
	"github.com/couchbase/gocb/v2"
	ui "github.com/gizak/termui/v3"
)

var logNum = 1
var maxLogSize = 5000000
var logFile *os.File

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
	MinVal        float64
	MaxVal        float64
	MaxChange     float64
	MaxChangeTime int
}

// Holds all incoming stat data from the server
type stats struct {

	// The main map holding the stat data (map[node][stat][300][float64])
	statBuffers map[string]map[string][]float64

	// List of stats to monitor
	statsList []string

	// Arrival time of messages from each node (map[node][300][time.Time])
	arrivalTimes map[string][]time.Time

	// A copy of stat alert information
	statInfo map[string]*configStatInfo

	// Lock for statInfo
	statInfoLock sync.RWMutex

	// Lock for statBuffers
	bufferLock sync.RWMutex

	// Lock for arrivalTimes
	timeLock sync.RWMutex

	// Lock for the list of stats
	statsListLock sync.RWMutex

	// Flag for first updation
	updated bool
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
		"connection_string",
		"couchbase://127.0.0.1:12000",
		"Provide the ip address for one of the search nodes",
	)
	config.reportPath = flag.String(
		"report", "./", "Provide path to print reports",
	)
	config.stats = make(map[string]*configStatInfo)
	config.alerts = make(map[string]*int)

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

	for {
		file, err := os.OpenFile(
			fmt.Sprintf("chronos%d.log", logNum),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)

		if err != nil {
			return err
		}

		fileStats, err := file.Stat()

		if err != nil {
			return err
		}

		size := fileStats.Size()

		if size > int64(maxLogSize) {
			logNum++
		} else {
			logFile = file
			break
		}
	}

	log.SetOutput(logFile)
	log.SetLoggerCallback(loggerFunc)

	return nil
}

// Setting log format
func loggerFunc(level, format string, args ...interface{}) string {

	ts := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	prefix := ts + " [" + level + "] "

	if file, err := os.Stat(fmt.Sprintf("chronos%d.log", logNum)); err == nil && level != "FATA" {
		size := file.Size()

		if size > int64(maxLogSize) {
			newLogFile, err := os.OpenFile(
				fmt.Sprintf("chronos%d.log", logNum+1),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0644,
			)

			if err != nil {
				log.Fatalf("log: Max log limit reached. Unable to create new log file %v", err)
			}

			logFile.Close()
			logNum++
			logFile = newLogFile
			log.SetOutput(logFile)
		}
	} else {
		fmt.Println(err)
	}
	if format != "" {
		return prefix + fmt.Sprintf(format, args...)
	}
	return prefix + fmt.Sprint(args...)
}

// Starting a connection to the server
func clusterInit(connectionString string, username string,
	password string) (*gocb.Cluster, error) {

	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
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

// Adding additional threshold values derived from the server
func addThresholds(node string, username string, password string, stats *stats) {

	url := node + "/api/manager"

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)

	if err != nil {
		log.Warnf("init: /api/manager request creation failed %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Warnf("init: /api/manager request failed %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnf("init: /api/manager response status not ok %v", err)
		return
	}

	dec := json.NewDecoder(resp.Body)

	var respMsg map[string]interface{}

	err = dec.Decode(&respMsg)

	if err != nil {
		log.Warnf("init: /api/manager response parsing failed %v", err)
		return
	}

	var num_bytes_used_ram_max_val float64
	if val, exists := respMsg["mgr"]; exists {
		if mgr, ok := val.(map[string]interface{}); ok {
			if val, exists := mgr["options"]; exists {
				if options, ok := val.(map[string]interface{}); ok {
					if val, exists := options["ftsMemoryQuota"]; exists {
						if val, ok := val.(string); ok {
							if thresholdVal, err :=
								strconv.ParseFloat(val, 64); err == nil {
								num_bytes_used_ram_max_val =
									thresholdVal
							}
						}
					}
				}
			}
		}
	}

	// Only use the threshold value if the user did not give one
	if num_bytes_used_ram_max_val == 0 {
		log.Warnf("init: getting max threshold from couchbase server for " +
			"num_bytes_used_ram failed")
	} else {
		stats.statInfoLock.Lock()
		if math.IsNaN(stats.statInfo["num_bytes_used_ram"].MaxVal) {
			stats.statInfo["num_bytes_used_ram"].MaxVal = num_bytes_used_ram_max_val
		}
		stats.statInfoLock.Unlock()
	}
}

// Initialize the stats struct with empty slices
func statsInit(config *config, nodesList []string) *stats {

	statBuffers := make(map[string]map[string][]float64)
	arrivalTimes := make(map[string][]time.Time)

	for _, node := range nodesList {
		statBuffers[node] = make(map[string][]float64)
		arrivalTimes[node] = make([]time.Time, 300)
	}

	return &stats{
		statBuffers:   statBuffers,
		statsList:     make([]string, 0),
		statInfo:      config.stats,
		arrivalTimes:  arrivalTimes,
		statInfoLock:  sync.RWMutex{},
		bufferLock:    sync.RWMutex{},
		timeLock:      sync.RWMutex{},
		statsListLock: sync.RWMutex{},
		updated:       false,
	}
}

// Setting up the widgets in a grid with relative ratios and positions
func gridInit(nodesTable *widgets.NodesTable, statsTable *widgets.StatsTable,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph,
	eventDisplay *widgets.EventDisplay,
	popupManager *widgets.PopupManager) *ui.Grid {

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
	popupManager.SetSize(termWidth, termHeight)
	return grid
}
