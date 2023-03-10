//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package widgets

// Stores stat name and node information for a graph
type GraphNodeSelect struct {
	Stat  string
	Nodes map[string]bool
}

func NewGraphNodeSelect() *GraphNodeSelect {
	return &GraphNodeSelect{
		Stat:  "",
		Nodes: make(map[string]bool),
	}
}

func (graphNodeSelect *GraphNodeSelect) GraphSelectInit(nodesList []string,
	nodesSelected []bool, stat string) {
	if graphNodeSelect.Stat != stat {
		graphNodeSelect.Nodes = make(map[string]bool)
		graphNodeSelect.Stat = stat
	}
	for i, node := range nodesList {
		graphNodeSelect.Nodes[node] = nodesSelected[i]
	}
}
