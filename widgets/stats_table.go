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
	"sort"

	ui "github.com/gizak/termui/v3"
)

// Colors corresponding to selected stats
const (
	stat1Color ui.Color = 46  // Green1
	stat2Color ui.Color = 165 // Magenta2
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

	// Array to track the number of lines used to display each row
	RowSize []int

	// Cursor position
	SelectedRow int

	// Selected row 1
	Stat1 int

	// Selected row 2
	Stat2 int

	// Row currently displayed on the first line
	TopRow int

	// Toggle indicating if widget is currently selected by the user
	selected bool
}

// Initializes a new stats table
func NewStatsTable() *StatsTable {

	return &StatsTable{
		Block:       ui.NewBlock(),
		Header:      "List of Stats",
		SelectedRow: 0,
		TopRow:      0,
		Stat1:       -1,
		Stat2:       -1,
		selected:    true,
		Rows:        make([]string, 0),
		RowSize:     make([]int, 0),
	}
}

// Render widget
func (table *StatsTable) Draw(buf *ui.Buffer) {

	if table.selected {
		table.BorderStyle = ui.NewStyle(ui.ColorWhite)
	} else {
		table.BorderStyle = ui.NewStyle(8)
	}
	table.Block.Draw(buf)

	// Horizontal padding of header from the left edge
	paddingHeader := 10

	// Vertical padding of rows from the header
	paddingRow := 4

	// Display style of the header
	styleHeader := ui.NewStyle(
		table.BorderStyle.Fg, ui.ColorClear, ui.ModifierBold,
	)

	// Display style of a normal row
	rowStyle := ui.NewStyle(
		table.BorderStyle.Fg, ui.ColorClear, ui.ModifierClear,
	)

	// Display style of a row with the cursor
	rowStyleSelected := ui.NewStyle(
		ui.ColorBlack, table.BorderStyle.Fg, ui.ModifierClear,
	)

	// Render header
	buf.SetString(
		table.Header, styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	table.RowSize = make([]int, len(table.Rows))

	// Keep track of when space is not available to render the row
	continueRender := true

	// Loop to render as many rows as possible within the bounds of the widget
	for rowNum, usedSpace := 0, 3; rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]

		var rowCells []ui.Cell

		// Check if the row is stat 1
		if rowNum == table.Stat1 {
			// Check if cursor is on row
			if rowNum == table.SelectedRow && table.selected {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						stat1Color, ui.ColorWhite, ui.ModifierClear,
					),
				)
			} else {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						ui.ColorBlack, stat1Color, ui.ModifierClear,
					),
				)
			}
			// Check if the row is stat 2
		} else if rowNum == table.Stat2 {
			// Check if cursor is on row
			if rowNum == table.SelectedRow && table.selected {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						stat2Color, ui.ColorWhite, ui.ModifierClear,
					),
				)
			} else {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						ui.ColorBlack, stat2Color, ui.ModifierClear,
					),
				)
			}
			// Check if cursor is on row
		} else if rowNum == table.SelectedRow && table.selected {
			rowCells = ui.ParseStyles(row, rowStyleSelected)
		} else {
			rowCells = ui.ParseStyles(row, rowStyle)
		}

		// Add padding for the rows
		rowCells = WrapCells(rowCells, uint(table.Inner.Dx()-2*paddingRow))

		// Split cells into multiple lines
		rowCellRows := ui.SplitCells(rowCells, '\n')

		// Render each cell if all the lines of text fit within the widget
		if len(rowCellRows) < table.Inner.Dy()-usedSpace &&
			rowNum >= table.TopRow && continueRender {
			for i, row := range rowCellRows {
				for _, cx := range ui.BuildCellWithXArray(row) {
					x, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(
						table.Inner.Min.X+paddingRow+x,
						table.Inner.Min.Y+usedSpace+i,
					))
				}
			}

			// Update the size of the text
			table.RowSize[rowNum] = len(rowCellRows)

			// Track the current line number
			usedSpace = usedSpace + table.RowSize[rowNum]
		} else {
			// Update number of rows taken to render even if not rendering
			table.RowSize[rowNum] = len(rowCellRows)

			if rowNum >= table.TopRow {
				// Stop rendering from this row
				continueRender = false
			}
		}
	}
}

// Handler function for scroll up
func (table *StatsTable) ScrollUp() {

	table.SelectedRow--

	table.CalcPos()
}

// Handler function for scroll down
func (table *StatsTable) ScrollDown() {

	table.SelectedRow++

	table.CalcPos()
}

// Handler function to indicate if cursor is on widget
func (table *StatsTable) ToggleTableSelect() {
	table.selected = !table.selected
}

// Handler function to enable selection of a row
func (table *StatsTable) SelectGraph(graphNum int) {

	if table.Stat2 != table.SelectedRow && table.Stat1 != table.SelectedRow {
		if graphNum == 1 {
			table.Stat1 = table.SelectedRow
		} else {
			table.Stat2 = table.SelectedRow
		}
	}
}

// Handler function for mouse click
func (table *StatsTable) HandleClick(x int, y int) {
	x = x - table.Min.X
	y = y - table.Min.Y
	if (x > 0 && x <= table.Inner.Dx()) && (y > 0 && y <= table.Inner.Dy()) {
		table.SelectedRow = (table.TopRow + y) - 4
		table.CalcPos()
	}
}

// Handler function to ensure cursor is never out of bounds
func (table *StatsTable) CalcPos() {

	if table.SelectedRow < 0 {
		table.SelectedRow = 0
	}

	if table.SelectedRow < table.TopRow {
		table.TopRow = table.SelectedRow
	}

	if table.SelectedRow > len(table.Rows)-1 {
		table.SelectedRow = len(table.Rows) - 1
	}
	if table.SelectedRow > table.TopRow+(table.Inner.Dy()-4) {
		table.TopRow = table.SelectedRow - (table.Inner.Dy() - 4)
	}

	if table.SelectedRow >= table.TopRow+table.RowsOnDisplay() {
		space := table.Inner.Dy() - 4

		for i := table.SelectedRow; i >= 0; i-- {
			space = space - table.RowSize[i]
			if space < 0 {
				if i == table.SelectedRow {
					table.TopRow = i
				} else {
					table.TopRow = i + 1
				}
				break
			}
		}
	}
}

// Calculate number of alerts currently on display
func (table *StatsTable) RowsOnDisplay() int {

	space := table.Inner.Dy() - 4
	rows := 0

	for i := table.TopRow; i < len(table.RowSize); i++ {

		space = space - table.RowSize[i]

		if space > 0 {
			rows = rows + 1
		} else {
			break
		}
	}

	return rows
}

// Handler function to determins if pixel is within the widget
func (table *StatsTable) Contains(x int, y int) bool {
	if x > table.Min.X && x <= table.Max.X &&
		y > table.Min.Y && y <= table.Max.Y {
		return true
	} else {
		return false
	}
}

// Handler function to add a new stat in a sorted order
func (table *StatsTable) AddStat(stat string) {

	index := sort.SearchStrings(table.Rows, stat)
	table.Rows = append(table.Rows, "")
	copy(table.Rows[index+1:], table.Rows[index:])
	table.Rows[index] = stat

}

// Handler function to remove a stat
func (table *StatsTable) RemoveStat(stat string) {

	for i, row := range table.Rows {
		if row == stat {
			table.Rows = append(table.Rows[:i], table.Rows[i+1:]...)
			return
		}
	}
}
