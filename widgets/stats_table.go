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

type StatsTable struct {
	*ui.Block
	Header      string
	Rows        []string
	SelectedRow int
	Stat1       int
	Stat2       int
	TopRow      int
	Selected    bool
}

func NewStatsTable() *StatsTable {
	return &StatsTable{
		Block:       ui.NewBlock(),
		Header:      "List of Stats",
		SelectedRow: 0,
		TopRow:      0,
		Stat1:       -1,
		Stat2:       -1,
		Selected:    true,
	}
}

func (table *StatsTable) Draw(buf *ui.Buffer) {
	table.Block.Draw(buf)

	paddingHeader := 10
	paddingRow := 4
	styleHeader := ui.NewStyle(
		ui.Theme.Default.Fg,
		ui.ColorClear,
		ui.ModifierBold,
	)
	rowStyle := ui.NewStyle(
		ui.Theme.Default.Fg,
		ui.ColorClear,
		ui.ModifierClear,
	)
	rowStyleSelected := ui.NewStyle(
		ui.ColorBlack,
		ui.ColorWhite,
		ui.ModifierClear,
	)

	buf.SetString(
		table.Header,
		styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	for rowNum := table.TopRow; rowNum < table.TopRow+table.Inner.Dy()-3 &&
		rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]
		y := rowNum + 3 - table.TopRow

		if rowNum == table.Stat1 {
			if rowNum == table.SelectedRow && table.Selected {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorCyan,
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
						ui.Theme.Default.Fg,
						ui.ColorCyan,
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			}
		} else if rowNum == table.Stat2 {
			if rowNum == table.SelectedRow && table.Selected {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorMagenta,
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
						ui.Theme.Default.Fg,
						ui.ColorMagenta,
						ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow,
						table.Inner.Min.Y+y,
					),
				)
			}
		} else if rowNum == table.SelectedRow && table.Selected {
			buf.SetString(
				row,
				rowStyleSelected,
				image.Pt(table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y),
			)
		} else {
			buf.SetString(
				row,
				rowStyle,
				image.Pt(table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y),
			)
		}
	}
}

func (table *StatsTable) ScrollUp() {

	table.SelectedRow--

	if table.SelectedRow < 0 {
		table.SelectedRow = 0
	}

	if table.SelectedRow < table.TopRow {
		table.TopRow = table.SelectedRow
	}
}

func (table *StatsTable) ScrollDown() {

	table.SelectedRow++

	if table.SelectedRow > len(table.Rows)-1 {
		table.SelectedRow = len(table.Rows) - 1
	}

	if table.SelectedRow > table.TopRow+table.Inner.Dy()-4 {
		table.TopRow = table.SelectedRow - (table.Inner.Dy() - 4)
	}
}

func (table *StatsTable) ToggleTableSelect() {
	table.Selected = !table.Selected
}

func (table *StatsTable) SelectGraph(graphNum int) {
	if graphNum == 1 && table.Stat2 != table.SelectedRow {
		table.Stat1 = table.SelectedRow
	} else if graphNum == 2 && table.Stat1 != table.SelectedRow {
		table.Stat2 = table.SelectedRow
	}
}
