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
	"io"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Id    int                `json:"id,omitempty"`
	Stats map[string]float64 `json:"stats,omitempty"`
}

func UpdateStats(stats map[string]map[string][]float64, config Config,
	nodeName string, statsList []string, logger *log.Logger,
	errChannel chan string) {

	url := "http://" + nodeName + "/api/statsStream"

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.Username, config.Password)

	if err != nil {
		errChannel <- "update_stats: Cannot connect to server" +
			nodeName + ":" + err.Error()
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		errChannel <- "update_stats: Invalid http response from server" +
			nodeName + ":" + err.Error()
		return
	}

	if resp.StatusCode != http.StatusOK {
		errChannel <- "update_stats: Status code is not OK:" +
			string(resp.StatusCode) + resp.Status
	}

	dec := json.NewDecoder(resp.Body)

	stats[nodeName] = make(map[string][]float64)

	for _, statName := range statsList {
		stats[nodeName][statName] = make([]float64, 110)
	}

	for range time.Tick(time.Second * 1) {
		var m Message
		err := dec.Decode(&m)
		if err != nil {
			if err == io.EOF {
				errChannel <- "update_stats: Server closed connection" +
					nodeName + err.Error()
				return
			}
			errChannel <- "update_stats: Invalid message recieved" +
				nodeName + err.Error()
			return
		}

		for _, statName := range statsList {
			val, ok := m.Stats[statName]
			stats[nodeName][statName] = stats[nodeName][statName][1:]
			if ok {
				stats[nodeName][statName] =
					append(stats[nodeName][statName], val)
			} else {
				stats[nodeName][statName] =
					append(stats[nodeName][statName], 0)
			}
		}
	}
}
