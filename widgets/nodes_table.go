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

// Widget to display nodes
// Each row displayed on a separate line
// Allows toggling of any row and a cursor to be present anywhere
type NodesTable struct {
	*ui.Block

	// Heading for the widget
	header string

	// Rows to be displayed
	Rows []string

	// Array to track the number of lines used to display each row
	RowSize []int

	// Cursor position
	SelectedRow int

	// Row currently displayed on the first line
	TopRow int

	// Toggle information for each node
	Nodes []bool

	// Toggle indicating if widget is currently selected by the user
	selected bool
}

// Initializes a new Nodes table
func NewNodesTable(NodesList []string) *NodesTable {

	rows := make([]string, 0)
	rows = append(rows, NodesList...)
	nodes := make([]bool, 0)

	for i := 0; i < len(NodesList); i++ {
		nodes = append(nodes, true)
	}

	return &NodesTable{
		Block:       ui.NewBlock(),
		header:      "List of Nodes",
		SelectedRow: 0,
		TopRow:      0,
		selected:    false,
		Rows:        rows,
		RowSize:     make([]int, 0),
		Nodes:       nodes,
	}
}

// Render widget
func (table *NodesTable) Draw(buf *ui.Buffer) {

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
	styleHeader := ui.NewStyle(table.BorderStyle.Fg, ui.ColorClear, ui.ModifierBold)

	// Render header
	buf.SetString(
		table.header, styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	table.RowSize = make([]int, len(table.Rows))

	// Keep track of when space is not available to render the row
	continueRender := true

	// Loop to render as many rows as possible within the bounds of the widget
	for rowNum, usedSpace := 0, 3; rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]

		var rowCells []ui.Cell

		// Check if current node is selected
		if table.Nodes[rowNum] {
			// Check if cursor on current node
			if rowNum == table.SelectedRow && table.selected {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						getColor(row), ui.ColorWhite, ui.ModifierClear,
					),
				)
			} else {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						getTextColor(getColor(row)), getColor(row), ui.ModifierClear,
					),
				)
			}
		} else {
			// Check if cursor on the current node
			if rowNum == table.SelectedRow && table.selected {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						ui.ColorBlack, ui.ColorWhite, ui.ModifierClear,
					),
				)
			} else {
				rowCells = ui.ParseStyles(
					row,
					ui.NewStyle(
						table.BorderStyle.Fg, ui.ColorClear, ui.ModifierClear,
					),
				)
			}
		}

		// Add padding for the rows
		rowCells = WrapCells(
			rowCells, uint(table.Inner.Dx()-2*paddingRow),
		)

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
func (table *NodesTable) ScrollUp() {

	table.SelectedRow--

	table.CalcPos()
}

// Handler function for scroll down
func (table *NodesTable) ScrollDown() {

	table.SelectedRow++

	table.CalcPos()
}

// Handler function for mouse click
func (table *NodesTable) HandleClick(x int, y int) {
	x = x - table.Min.X
	y = y - table.Min.Y
	if (x > 0 && x <= table.Inner.Dx()) && (y > 0 && y <= table.Inner.Dy()) {
		table.SelectedRow = (table.TopRow + y) - 4
		table.CalcPos()
	}
}

// Handler function to ensure cursor is never out of bounds
func (table *NodesTable) CalcPos() {

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
func (table *NodesTable) RowsOnDisplay() int {

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

// Handler function to indicate if cursor is on widget
func (table *NodesTable) ToggleTableSelect() {
	table.selected = !table.selected
}

// Handler function to initialize toggle info
func (table *NodesTable) SelectStat(stat string,
	nodes []*NodeData, graphStat string) {

	if stat == graphStat {
		for i, row := range table.Rows {
			for _, node := range nodes {
				if node.Node == row {
					table.Nodes[i] = node.Active
					break
				}
			}
		}
	} else {
		for i := range table.Rows {
			table.Nodes[i] = true
		}
	}
}

// Handler function to toggle select for a row
func (table *NodesTable) SelectNode() {
	table.Nodes[table.SelectedRow] = !table.Nodes[table.SelectedRow]
}

// Handler function to add a new row
func (table *NodesTable) AddNode(node string) {
	table.Nodes = append(table.Nodes, true)
	table.Rows = append(table.Rows, node)
}

// Handler function to remove a row
func (table *NodesTable) RemoveNode(node string) {

	tempRows := make([]string, 0)
	tempNodes := make([]bool, 0)

	for i, nodeName := range table.Rows {
		if nodeName != node {
			tempRows = append(tempRows, nodeName)
			tempNodes = append(tempNodes, table.Nodes[i])
		}
	}

	table.Nodes = tempNodes
	table.Rows = tempRows
}

// Handler function to determins if pixel is within the widget
func (table *NodesTable) Contains(x int, y int) bool {
	if x > table.Min.X && x <= table.Max.X &&
		y > table.Min.Y && y <= table.Max.Y {
		return true
	} else {
		return false
	}
}

// Function to wrap long words without line breaks into multiple lines
func WrapCells(cells []ui.Cell, width uint) []ui.Cell {
	str := ui.CellsToString(cells)
	newStr := ""
	wrappedCells := make([]ui.Cell, 0)

	for {
		if len(str) > int(width) {
			newStr = newStr + str[:width] + "\n"
			str = str[width:]
		} else {
			newStr = newStr + str
			break
		}
	}

	i := 0
	for _, _rune := range newStr {
		if _rune == '\n' {
			wrappedCells = append(
				wrappedCells,
				ui.Cell{
					Rune:  _rune,
					Style: ui.StyleClear,
				})
		} else {
			wrappedCells = append(
				wrappedCells,
				ui.Cell{
					Rune:  _rune,
					Style: cells[i].Style,
				})
			i++
		}
	}
	return wrappedCells
}
