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
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/couchbaselabs/chronos/widgets"

	ui "github.com/gizak/termui/v3"
)

type Config struct {
	Username string              `json:"username,omitempty"`
	Password string              `json:"password,omitempty"`
	Nodes    map[string][]string `json:"nodes,omitempty"`
}

type StatInfo struct {
	MinVal    float64
	MaxVal    float64
	MaxChange float64
}

func flagsInit() string {

	configPath := flag.String(
		"config",
		"config.json",
		"Provide the relative path to the config file",
	)
	flag.Parse()

	return *configPath
}

func logsInit() (*log.Logger, *log.Logger, error) {

	logRender, err := os.OpenFile(
		"Render.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, nil, err
	}

	loggerRender := log.New(logRender, "Render Log", log.LstdFlags)

	logUpdate, err := os.OpenFile(
		"Update.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, nil, err
	}

	loggerUpdate := log.New(logUpdate, "Update Log", log.LstdFlags)

	return loggerRender, loggerUpdate, nil
}

func configInit(configPath string, logger *log.Logger) (Config, error) {

	configRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	err = json.Unmarshal([]byte(configRaw), &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func uiInit() error {

	return ui.Init()

}

func statNodeInit(config Config) map[string][]string {

	statNodes := make(map[string][]string)

	for nodeName, statList := range config.Nodes {
		for _, stat := range statList {

			if _, ok := statNodes[stat]; !ok {
				statNodes[stat] = make([]string, 0)
			}

			statNodes[stat] = append(statNodes[stat], nodeName)
		}
	}

	return statNodes
}

func gridInit(nodesTable *widgets.NodesTable, statsTable *widgets.StatsTable,
	lineChart1 *widgets.LineGraph, lineChart2 *widgets.LineGraph) *ui.Grid {

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
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	return grid
}
