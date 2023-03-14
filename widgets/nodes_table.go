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

// Defining colors for each line dased on its index
var colors = ui.Theme.Plot.Lines

// Widget to display Nodes
// Each row displayed on a separate line
// Allows toggling of any row and a cursor to be present anywhere
type NodesTable struct {
	*ui.Block

	// Heading for the widget
	header string

	// Rows to be displayed
	Rows []string

	// Cursor position
	selectedRow int

	// Row currently displayed on the first line
	topRow int

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
		selectedRow: 0,
		topRow:      0,
		selected:    false,
		Rows:        rows,
		Nodes:       nodes,
	}
}

// Render widget
func (table *NodesTable) Draw(buf *ui.Buffer) {
	table.Block.Draw(buf)

	// Horizontal padding of header from the left edge
	paddingHeader := 10

	// Vertical padding of rows from the header
	paddingRow := 4

	// Display style of the header
	styleHeader := ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)

	// Render header
	buf.SetString(
		table.header, styleHeader,
		image.Pt(table.Inner.Min.X+paddingHeader, table.Inner.Min.Y+1),
	)

	// Loop to render as many rows as possible within the bounds of the widget
	for rowNum := table.topRow; rowNum < table.topRow+table.Inner.Dy()-3 &&
		rowNum < len(table.Rows); rowNum++ {

		row := table.Rows[rowNum]
		y := rowNum + 3 - table.topRow

		// Check if current node is selected
		if table.Nodes[rowNum] {
			// Check if cursor on current node
			if rowNum == table.selectedRow && table.selected {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum], ui.ColorWhite, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						ui.ColorWhite, colors[rowNum], ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			}
		} else {
			// Check if cursor on the current node
			if rowNum == table.selectedRow && table.selected {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum], ui.ColorBlack, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			} else {
				buf.SetString(
					row,
					ui.NewStyle(
						colors[rowNum], ui.ColorClear, ui.ModifierClear,
					),
					image.Pt(
						table.Inner.Min.X+paddingRow, table.Inner.Min.Y+y,
					),
				)
			}
		}
	}
}

// Handler function to scroll up
func (table *NodesTable) ScrollUp() {

	table.selectedRow--

	if table.selectedRow < 0 {
		table.selectedRow = 0
	}

	if table.selectedRow < table.topRow {
		table.topRow = table.selectedRow
	}
}

// Handler function to scroll down
func (table *NodesTable) ScrollDown() {

	table.selectedRow++

	if table.selectedRow > len(table.Rows)-1 {
		table.selectedRow = len(table.Rows) - 1
	}

	if table.selectedRow > table.topRow+table.Inner.Dy()-4 {
		table.topRow = table.selectedRow - (table.Inner.Dy() - 4)
	}
}

// Handler function to indicate if cursor is on widget
func (table *NodesTable) ToggleTableSelect() {
	table.selected = !table.selected
}

// Handler function to initialize toggle info
func (table *NodesTable) SelectStat(stat string,
	nodes map[string]bool, nodeSelectStat string) {

	for i, row := range table.Rows {
		if stat == nodeSelectStat {
			table.Nodes[i] = nodes[row]
		} else {
			table.Nodes[i] = true
		}
	}
}

// Handler function to toggle select for a row
func (table *NodesTable) SelectNode() {
	table.Nodes[table.selectedRow] = !table.Nodes[table.selectedRow]
}

// Handler function to add a new node
func (table *NodesTable) AddNode(node string) {
	table.Nodes = append(table.Nodes, true)
	table.Rows = append(table.Rows, node)
}

// Handler function to remove an existing node
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
