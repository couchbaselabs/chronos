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

	// Data slices holding the values to be displayed
	Data [][]float64

	// Maximum value in the entire data slice
	// Used to vertically size the graph
	maxVal float64

	// Defines colors for each line
	lineColors []ui.Color

	// Defines colors for the axes
	axesColor ui.Color

	// Horizontal scaling of the graph
	horizontalScale int
}

// Values of different paddings
const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

// Initializes a new line graph
func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block:           ui.NewBlock(),
		lineColors:      ui.Theme.Plot.Lines,
		axesColor:       ui.Theme.Plot.Axes,
		horizontalScale: 1,
		Data:            [][]float64{},
	}
}

// Function to render the lines
func (graph *LineGraph) renderBraille(buf *ui.Buffer, drawArea image.Rectangle,
	maxVal float64) {

	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	for i, line := range graph.Data {
		// Check to prevent rendering of empty lines
		if len(line) > 1 {
			previousHeight :=
				int((line[1] / maxVal) * float64(drawArea.Dy()-1))
			for j, val := range line[1:] {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				canvas.SetLine(
					image.Pt(
						(drawArea.Min.X+(j*graph.horizontalScale))*2,
						(drawArea.Max.Y-previousHeight-1)*4,
					),
					image.Pt(
						(drawArea.Min.X+((j+1)*graph.horizontalScale))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					ui.SelectColor(graph.lineColors, i),
				)
				previousHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

// Render axes of the graph
func (graph *LineGraph) plotAxes(buf *ui.Buffer, maxVal float64) {

	// Render origin
	buf.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, ui.NewStyle(ui.ColorWhite)),
		image.Pt(
			graph.Inner.Min.X+yAxisLabelsWidth,
			graph.Inner.Max.Y-xAxisLabelsHeight-1,
		),
	)
	// Render x axis
	for i := yAxisLabelsWidth + 1; i < graph.Inner.Dx(); i++ {
		buf.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(
				i+graph.Inner.Min.X, graph.Inner.Max.Y-xAxisLabelsHeight-1,
			),
		)
	}
	// Render y axis
	for i := 0; i < graph.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, i+graph.Inner.Min.Y),
		)
	}

	// Render 0
	buf.SetString(
		"0",
		ui.NewStyle(ui.ColorWhite),
		image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, graph.Inner.Max.Y-1),
	)
	// Render other x axis labels
	for x := graph.Inner.Min.X + yAxisLabelsWidth +
		(xAxisLabelsGap)*graph.horizontalScale + 1; x < graph.Inner.Max.X-1; {

		label := fmt.Sprintf(
			"%d",
			(x-(graph.Inner.Min.X+yAxisLabelsWidth)-1)/
				(graph.horizontalScale)+1,
		)
		buf.SetString(
			label, ui.NewStyle(ui.ColorWhite), image.Pt(x, graph.Inner.Max.Y-1),
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
			ui.NewStyle(ui.ColorWhite),
			image.Pt(
				graph.Inner.Min.X, graph.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2,
			),
		)
	}
}

// Render widget
func (graph *LineGraph) Draw(buf *ui.Buffer) {
	graph.Block.Draw(buf)

	// Identify max value within the data that can be displayed
	maxVal := graph.maxVal
	maxData, _ := ui.GetMaxFloat64From2dSlice(graph.Data)

	// Set new max value if out of bounds
	if maxVal > 2*maxData || maxVal < maxData {
		maxVal = maxData * 1.25
		graph.maxVal = maxVal
	}

	// Render axes
	graph.plotAxes(buf, maxVal)

	// Identify leftover space
	drawArea := image.Rect(
		graph.Inner.Min.X+yAxisLabelsWidth+1, graph.Inner.Min.Y,
		graph.Inner.Max.X, graph.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	// Render lines
	graph.renderBraille(buf, drawArea, maxVal)

}
