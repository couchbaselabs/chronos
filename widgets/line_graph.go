//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package widgets

import (
	"fmt"
	"image"
	"strconv"

	ui "github.com/gizak/termui/v3"
)

// Widget to display a set of lines as a graph
// Each line displayed in a different color
type LineGraph struct {
	*ui.Block

	// Maximum value in the entire data slice
	// Used to vertically size the graph
	maxVal float64

	// Horizontal scaling of the graph
	horizontalScale int

	// Stat name corresponding to the graph
	Stat string

	// Stat color corresponding to the graph
	statColor ui.Color

	// Nodes list for the graph
	Nodes []*NodeData

	// Toggle to display legend
	legend bool

	// Toggle indicating if widget is currently selected by the user
	Selected bool
}

// Stores the data for a single node
type NodeData struct {
	Node   string
	Line   []float64
	color  ui.Color
	Active bool
}

// Values of different paddings
const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

var (
	// List of distinctly different colors for the lines
	// Each number represents an Xterm color
	colors = []ui.Color{
		46, 201, 33, 208, 108, 125, 64, 57, 51, 226, 213, 24, 48, 196, 98, 216,
		118, 117, 204, 30, 242, 88, 142, 198, 157, 38, 130, 128, 53, 76, 146,
		85, 185, 170, 237, 27, 173, 54, 19, 67, 178, 133, 94, 197, 42, 69, 228,
		99, 83, 21,
	}

	// List of text colors for the corresponding color highlights
	// 15 is White
	// 0 is Black
	textColors = []ui.Color{
		0, 15, 15, 15, 0, 15, 15, 15, 0, 0, 0, 15, 0, 15, 15, 0, 0, 0, 15, 15, 15, 15, 15, 15,
		0, 15, 15, 15, 15, 15, 0, 0, 0, 15, 15, 15, 0, 15, 15, 15, 0, 15, 15, 15, 15, 15, 0, 15,
		0, 15,
	}

	// Struct to hold the node to color mappings
	nodeColors = make(map[string]struct {
		color     ui.Color
		textColor ui.Color
	})
)

// Initializes a new line graph
func NewLineGraph(nodesList []string, graphNum int) *LineGraph {

	var color ui.Color

	if graphNum == 1 {
		color = stat1Color
	} else {
		color = stat2Color
	}

	lineGraph := &LineGraph{
		Block:           ui.NewBlock(),
		Stat:            "",
		statColor:       color,
		Nodes:           make([]*NodeData, 0),
		horizontalScale: 1,
		legend:          false,
		Selected:        false,
	}

	for _, node := range nodesList {

		color := getColor(node)

		nodeData := &NodeData{
			Node:   node,
			Line:   make([]float64, 0),
			color:  color,
			Active: true,
		}

		lineGraph.Nodes = append(lineGraph.Nodes, nodeData)
	}

	return lineGraph
}

// Function to render the lines
func (graph *LineGraph) renderBraille(buf *ui.Buffer, drawArea image.Rectangle,
	maxVal float64) {

	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	for _, nodeData := range graph.Nodes {

		line := nodeData.Line

		// Check to prevent rendering of empty or deselected lines
		if len(line) != 0 && nodeData.Active {

			prevHeight := int((line[len(line)-1] / maxVal) *
				float64(drawArea.Dy()-1))

			for j := 1; j < len(line)-1; j++ {

				val := line[len(line)-1-j]
				height := int((val / maxVal) * float64(drawArea.Dy()-1))

				if (drawArea.Max.X - ((j - 1) * graph.horizontalScale)) >=
					drawArea.Min.X {

					canvas.SetLine(
						image.Pt((drawArea.Max.X-(j*graph.horizontalScale))*2,
							(drawArea.Max.Y-height-1)*4),
						image.Pt(
							(drawArea.Max.X-((j-1)*graph.horizontalScale))*2,
							(drawArea.Max.Y-prevHeight-1)*4,
						),
						nodeData.color,
					)
				}
				prevHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

// Render axes of the graph
func (graph *LineGraph) plotAxes(buf *ui.Buffer, maxVal float64) {

	var axesStyle ui.Style

	if graph.Selected {
		axesStyle = ui.NewStyle(ui.ColorWhite)
	} else {
		axesStyle = ui.NewStyle(8)
	}

	// Render origin
	buf.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, axesStyle),
		image.Pt(
			graph.Inner.Min.X+yAxisLabelsWidth,
			graph.Inner.Max.Y-xAxisLabelsHeight-1,
		),
	)

	// Render x axis
	for i := yAxisLabelsWidth + 1; i < graph.Inner.Dx(); i++ {
		buf.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, axesStyle),
			image.Pt(
				i+graph.Inner.Min.X, graph.Inner.Max.Y-xAxisLabelsHeight-1,
			),
		)
	}

	// Render y axis
	for i := 0; i < graph.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, axesStyle),
			image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, i+graph.Inner.Min.Y),
		)
	}

	// Render 0
	buf.SetString(
		"0",
		axesStyle,
		image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, graph.Inner.Max.Y-1),
	)

	// Render other x axis label
	for x := graph.Inner.Min.X + yAxisLabelsWidth +
		(xAxisLabelsGap)*graph.horizontalScale + 1; x < graph.Inner.Max.X-1; {

		label := fmt.Sprintf(
			"%d", (x-(graph.Inner.Min.X+yAxisLabelsWidth)-1)/
				(graph.horizontalScale)+1,
		)
		buf.SetString(
			label, axesStyle, image.Pt(x, graph.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * graph.horizontalScale
	}

	// Render y axis labels
	verticalScale := maxVal / float64(graph.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < graph.Inner.Dy()-1; i++ {
		buf.SetString(
			strconv.FormatFloat(
				float64(i)*verticalScale*(yAxisLabelsGap+1), 'E', 2, 64,
			),
			axesStyle,
			image.Pt(
				graph.Inner.Min.X, graph.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2,
			),
		)
	}
}

// Render widget
func (graph *LineGraph) Draw(buf *ui.Buffer) {

	if graph.Selected {
		graph.BorderStyle = ui.NewStyle(ui.ColorWhite)
	} else {
		graph.BorderStyle = ui.NewStyle(8)
	}

	graph.Block.Draw(buf)

	// Identify length of the display
	dispLength := graph.DispLength()

	// Identify max value within the data that can be displayed
	maxData := graph.MaxData(dispLength)

	// Set new max value if out of bounds
	if graph.maxVal > 2*maxData || graph.maxVal < maxData {
		graph.maxVal = maxData * 1.25
	}

	// Render axes
	graph.plotAxes(buf, graph.maxVal)

	// Identify leftover space
	drawArea := image.Rect(
		graph.Inner.Min.X+yAxisLabelsWidth+1, graph.Inner.Min.Y,
		graph.Inner.Max.X, graph.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	// Render lines
	graph.renderBraille(buf, drawArea, graph.maxVal)

	// Render legend if any stat is selected
	if graph.Stat != "" {
		graph.renderStat(buf)
		if graph.legend {
			graph.renderLegend(buf)
		}
	}
}

// Render the stat corresponding to the graph
func (graph *LineGraph) renderStat(buf *ui.Buffer) {
	buf.SetString(
		graph.Stat,
		ui.NewStyle(
			graph.statColor, ui.ColorClear, ui.ModifierClear,
		),
		image.Pt(
			graph.Inner.Min.X+11, graph.Inner.Min.Y+1,
		),
	)
}

// Render the node names in its line color
func (graph *LineGraph) renderLegend(buf *ui.Buffer) {

	for i, nodeData := range graph.Nodes {
		buf.SetString(
			nodeData.Node,
			ui.NewStyle(
				nodeData.color, ui.ColorClear, ui.ModifierClear,
			),
			image.Pt(
				graph.Inner.Min.X+13, graph.Inner.Min.Y+3+i,
			),
		)
	}
}

// Calculate length of data that can be rendered
func (graph *LineGraph) DispLength() int {

	return int((graph.Inner.Max.X-(graph.Inner.Min.X+yAxisLabelsWidth+1))/
		graph.horizontalScale) + 1

}

// Calculate the max value of the data that can be rendered
func (graph *LineGraph) MaxData(dispLength int) float64 {

	var maxData float64

	for _, nodeData := range graph.Nodes {
		line := nodeData.Line
		if len(line) != 0 {
			if len(line)-dispLength-1 >= 0 {
				for _, val := range line[len(line)-dispLength-1:] {
					if val > maxData {
						maxData = val
					}
				}
			} else {
				for _, val := range line {
					if val > maxData {
						maxData = val
					}
				}
			}
		}
	}

	return maxData
}

// Handler to toggle display of the line corresponding to a node
func (graph *LineGraph) SelectNode(node string) {

	for _, nodeData := range graph.Nodes {
		if nodeData.Node == node {
			nodeData.Active = !nodeData.Active
		}
	}
}

// Handler to add a new node
func (graph *LineGraph) AddNode(node string) {

	newNodeData := &NodeData{
		Node:   node,
		Line:   make([]float64, 0),
		color:  getColor(node),
		Active: true,
	}

	graph.Nodes = append(graph.Nodes, newNodeData)
}

// Handler to remove an existing node
func (graph *LineGraph) RemoveNode(node string) {

	delete(nodeColors, node)

	for i, nodeData := range graph.Nodes {
		if nodeData.Node == node {
			graph.Nodes = append(graph.Nodes[:i], graph.Nodes[i+1:]...)
		}
	}
}

// Handler to toggle the display of legend
func (graph *LineGraph) ToggleLegend() {
	graph.legend = !graph.legend
}

// Function to return the corresponding color of a node
// Assigns a new color if called for the first time for a node
func getColor(node string) ui.Color {

	color, assigned := nodeColors[node]

	if !assigned {
		for i, defaultColor := range colors {
			used := false
			for _, nodeColor := range nodeColors {
				if nodeColor.color == defaultColor {
					used = true
					break
				}
			}
			if !used {
				color.color = defaultColor
				color.textColor = textColors[i]
				nodeColors[node] = color
				break
			}
		}
	}

	return color.color
}

// Function to return the color of the text for a particular highlight color
func getTextColor(color ui.Color) ui.Color {

	for i, highlightColor := range colors {
		if highlightColor == color {
			return textColors[i]
		}
	}
	return 0
}

// Handler function to indicate if graph is connected to nodes table
func (graph *LineGraph) SelectGraph() {
	graph.Selected = true
}

// Handler function to indicate if cursor is not connected to nodes table
func (graph *LineGraph) UnSelectGraph() {
	graph.Selected = false
}
