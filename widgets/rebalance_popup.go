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
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
)

// Widget to display a rebalance popup
type RebalancePopup struct {
	*ui.Block

	// Popup text
	text string

	// Toggle to indicate if popup should be displayed
	Display bool

	// Max time to live for the popup
	ttl time.Time

	// Lock for the popup
	DisplayLock sync.RWMutex
}

// Create a new rebalance popup
func NewRebalancePopup() *RebalancePopup {
	return &RebalancePopup{
		Block:   ui.NewBlock(),
		text:    "Cluster under rebalance",
		Display: false,
	}
}

// Handler function to resize popup with respect to terminal dimensions
func (popup *RebalancePopup) Resize(width int, height int) {
	popup.Block.SetRect(
		int(0.9*float64(width)), 0, width, int(0.1*float64(height)),
	)
}

// Handler to update the rebalance popup
func (popup *RebalancePopup) SetRebalance() {
	popup.DisplayLock.Lock()
	popup.Display = true
	popup.ttl = time.Now().Add(time.Duration(5) * time.Second)
	popup.DisplayLock.Unlock()
}

// Render the widget
func (popup *RebalancePopup) Draw(buf *ui.Buffer) {

	popup.DisplayLock.RLock()
	// If popup needs to be displayed
	if popup.Display {
		// If popup is active
		if time.Now().Before(popup.ttl) {
			popup.Block.Draw(buf)

			// Parse text into cells
			cells := ui.ParseStyles(
				popup.text,
				ui.NewStyle(
					ui.ColorWhite, ui.ColorClear, ui.ModifierClear,
				),
			)
			// Pad cells within the widget
			cells = ui.WrapCells(cells, uint(popup.Inner.Dx()))

			// Split cells into rows
			rows := ui.SplitCells(cells, '\n')

			// Render row by row
			for y, row := range rows {
				if y+popup.Inner.Min.Y >= popup.Inner.Max.Y {
					break
				}
				row = ui.TrimCells(row, popup.Inner.Dx())
				for _, cx := range ui.BuildCellWithXArray(row) {
					x, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(x, y).Add(popup.Inner.Min))
				}
			}
		} else {
			popup.Display = false
		}
	}
	popup.DisplayLock.RUnlock()
}
