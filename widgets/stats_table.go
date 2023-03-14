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

// Widget to display a list of stats
// Each row displayed on one line
// Allows selection of 2 rows with cursor on any row
type StatsTable struct {
	*ui.Block

	// Heading for the widget
	Header string

	// Rows to be displayed
	Rows []string

	// Cursor position
	SelectedRow int

	// Selected row 1
	Stat1 int

	// Selected row 2
	Stat2 int

	// Row currently displayed on the first line
	topRow int

	// Toggle indicating if widget is currently selected by the user
	selected bool
}

// Initializes a new stats table
func NewStatsTable(statsList []string) *StatsTable {

	rows := make([]string, 0)
	rows = append(rows, statsList...)

	return &StatsTable{
		Block:       ui.NewBlock(),
		Header:      "List of Stats",
		SelectedRow: 0,
		topRow:      0,
		Stat1:       -1,
		Stat2:       -1,
		selected:    true,
		Rows:        rows,
	}
}

// Render widget
func (table *StatsTable) Draw(buf *ui.Buffer) {
	table.Block.Draw(buf)

	// Horizontal padding of header from the left edge
	paddingHeader := 10

	// Vertical padding of rows from the header
	paddingRow := 4

	// Display style of the header
	styleHeader := ui.NewStyle(
		ui.Theme.Default.Fg, ui.ColorClear, ui.ModifierBold,
	)

	// Display style of a normal row
	rowStyle := ui.NewStyle(
		ui.Theme.Default.Fg, ui.ColorClear, ui.ModifierClear,
	)

	// Display style of a row with the cursor
	rowStyleSelected := ui.NewStyle(
		ui.ColorBlack, ui.ColorWhite, ui.ModifierClear,
	)

	// Render header
	buf.SetString(
		table.Header, styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	// Loop to render as many rows as possible within the bounds of the widget
	for rowNum := table.topRow; rowNum < table.topRow+table.Inner.Dy()-3 &&
		rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]
		y := rowNum + 3 - table.topRow

		// Check if the current row is selected row 1
		if rowNum == table.Stat1 {
			// Check if cursor on selected row 1
			if rowNum == table.SelectedRow && table.selected {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorCyan, ui.ColorWhite, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.Theme.Default.Fg, ui.ColorCyan, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			}
			// Check if current row is selected row 2
		} else if rowNum == table.Stat2 {
			// Check if cursor on selected row 2
			if rowNum == table.SelectedRow && table.selected {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorMagenta, ui.ColorWhite, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.Theme.Default.Fg, ui.ColorMagenta, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			}
			// Check if cursor on row
		} else if rowNum == table.SelectedRow && table.selected {
			buf.SetString(
				row, rowStyleSelected,
				image.Pt(table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y),
			)
		} else {
			buf.SetString(
				row, rowStyle,
				image.Pt(table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y),
			)
		}
	}
}

// Handler function for scroll uo
func (table *StatsTable) ScrollUp() {

	table.SelectedRow--

	if table.SelectedRow < 0 {
		table.SelectedRow = 0
	}

	if table.SelectedRow < table.topRow {
		table.topRow = table.SelectedRow
	}
}

// Handler function for scroll down
func (table *StatsTable) ScrollDown() {

	table.SelectedRow++

	if table.SelectedRow > len(table.Rows)-1 {
		table.SelectedRow = len(table.Rows) - 1
	}

	if table.SelectedRow > table.topRow+table.Inner.Dy()-4 {
		table.topRow = table.SelectedRow - (table.Inner.Dy() - 4)
	}
}

// Handler function to indicate if cursor is on widget
func (table *StatsTable) ToggleTableSelect() {
	table.selected = !table.selected
}

// Handler function to enable selection of a row
func (table *StatsTable) SelectGraph(graphNum int) {
	if graphNum == 1 && table.Stat2 != table.SelectedRow {
		table.Stat1 = table.SelectedRow
	} else if graphNum == 2 && table.Stat1 != table.SelectedRow {
		table.Stat2 = table.SelectedRow
	}
}
