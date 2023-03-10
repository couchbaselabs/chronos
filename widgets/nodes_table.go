//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package widgets

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

var colors = ui.Theme.Plot.Lines

type NodesTable struct {
	*ui.Block
	Header      string
	Rows        []string
	SelectedRow int
	TopRow      int
	Nodes       []bool
	Selected    bool
}

func NewNodesTable() *NodesTable {
	return &NodesTable{
		Block:       ui.NewBlock(),
		Header:      "List of Nodes",
		SelectedRow: 0,
		TopRow:      0,
		Selected:    false,
	}
}

func (table *NodesTable) Draw(buf *ui.Buffer) {
	table.Block.Draw(buf)

	paddingHeader := 10
	paddingRow := 4
	styleHeader := ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)

	buf.SetString(
		table.Header,
		styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	for rowNum := table.TopRow; rowNum < table.TopRow+table.Inner.Dy()-3 &&
		rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]
		y := rowNum + 3 - table.TopRow

		if table.Nodes[rowNum] {
			if rowNum == table.SelectedRow && table.Selected {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum],
						ui.ColorWhite,
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorWhite,
						colors[rowNum],
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			}
		} else {
			if rowNum == table.SelectedRow && table.Selected {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum],
						ui.ColorBlack,
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum],
						ui.ColorClear,
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			}
		}
	}
}

func (table *NodesTable) ScrollUp() {

	table.SelectedRow--

	if table.SelectedRow < 0 {
		table.SelectedRow = 0
	}

	if table.SelectedRow < table.TopRow {
		table.TopRow = table.SelectedRow
	}
}

func (table *NodesTable) ScrollDown() {

	table.SelectedRow++

	if table.SelectedRow > len(table.Rows)-1 {
		table.SelectedRow = len(table.Rows) - 1
	}

	if table.SelectedRow > table.TopRow+table.Inner.Dy()-4 {
		table.TopRow = table.SelectedRow - (table.Inner.Dy() - 4)
	}
}

func (table *NodesTable) ToggleTableSelect() {
	table.Selected = !table.Selected
}

func (table *NodesTable) SelectStat(stat string, statNodes map[string][]string,
	nodes map[string]bool, nodeSelectStat string) {

	table.Rows = make([]string, 0)
	table.Nodes = make([]bool, 0)

	for _, node := range statNodes[stat] {
		table.Rows = append(table.Rows, node)
		if stat == nodeSelectStat {
			table.Nodes = append(table.Nodes, nodes[node])
		} else {
			table.Nodes = append(table.Nodes, true)
		}
	}
}

func (table *NodesTable) SelectNode() {
	table.Nodes[table.SelectedRow] = !table.Nodes[table.SelectedRow]
}
