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

type LineGraph struct {
	*ui.Block

	Data            [][]float64
	MaxVal          float64
	LineColors      []ui.Color
	AxesColor       ui.Color
	HorizontalScale int
	DotMarkerRune   rune
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block:           ui.NewBlock(),
		LineColors:      ui.Theme.Plot.Lines,
		AxesColor:       ui.Theme.Plot.Axes,
		HorizontalScale: 1,
		DotMarkerRune:   ui.DOT,
		Data:            [][]float64{},
	}
}

func (graph *LineGraph) renderBraille(buf *ui.Buffer, drawArea image.Rectangle,
	maxVal float64) {

	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	for i, line := range graph.Data {
		if len(line) > 1 {
			previousHeight :=
				int((line[1] / maxVal) * float64(drawArea.Dy()-1))
			for j, val := range line[1:] {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				canvas.SetLine(
					image.Pt(
						(drawArea.Min.X+(j*graph.HorizontalScale))*2,
						(drawArea.Max.Y-previousHeight-1)*4,
					),
					image.Pt(
						(drawArea.Min.X+((j+1)*graph.HorizontalScale))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					ui.SelectColor(graph.LineColors, i),
				)
				previousHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

func (graph *LineGraph) plotAxes(buf *ui.Buffer, maxVal float64) {
	// draw origin cell
	buf.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, ui.NewStyle(ui.ColorWhite)),
		image.Pt(
			graph.Inner.Min.X+yAxisLabelsWidth,
			graph.Inner.Max.Y-xAxisLabelsHeight-1,
		),
	)
	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < graph.Inner.Dx(); i++ {
		buf.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(
				i+graph.Inner.Min.X,
				graph.Inner.Max.Y-xAxisLabelsHeight-1,
			),
		)
	}
	// draw y axis line
	for i := 0; i < graph.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, i+graph.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		"0",
		ui.NewStyle(ui.ColorWhite),
		image.Pt(graph.Inner.Min.X+yAxisLabelsWidth, graph.Inner.Max.Y-1),
	)
	// draw rest
	for x := graph.Inner.Min.X + yAxisLabelsWidth +
		(xAxisLabelsGap)*graph.HorizontalScale + 1; x < graph.Inner.Max.X-1; {

		label := fmt.Sprintf(
			"%d",
			(x-(graph.Inner.Min.X+yAxisLabelsWidth)-1)/
				(graph.HorizontalScale)+1,
		)
		buf.SetString(
			label,
			ui.NewStyle(ui.ColorWhite),
			image.Pt(x, graph.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * graph.HorizontalScale
	}
	// draw y axis labels
	verticalScale := maxVal / float64(graph.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < graph.Inner.Dy()-1; i++ {
		buf.SetString(
			strconv.FormatFloat(
				float64(i)*verticalScale*(yAxisLabelsGap+1),
				'E',
				2,
				64,
			),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(
				graph.Inner.Min.X,
				graph.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2,
			),
		)
	}
}

func (graph *LineGraph) Draw(buf *ui.Buffer) {
	graph.Block.Draw(buf)

	maxVal := graph.MaxVal
	maxData, _ := ui.GetMaxFloat64From2dSlice(graph.Data)

	if maxVal > 2*maxData || maxVal < maxData {
		maxVal = maxData * 1.25
		graph.MaxVal = maxVal
	}

	graph.plotAxes(buf, maxVal)

	drawArea := image.Rect(
		graph.Inner.Min.X+yAxisLabelsWidth+1, graph.Inner.Min.Y,
		graph.Inner.Max.X, graph.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	graph.renderBraille(buf, drawArea, maxVal)

}
