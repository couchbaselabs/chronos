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
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
)

// Widget to display a popup
type popup struct {
	*ui.Block

	// Popup text
	text string

	// Popup text after splitting into different rows and aligning
	processedText []string

	// Toggle to indicate if popup should be displayed
	display bool

	// Max time to live for the popup
	ttl time.Time

	// Type of popup
	popupType string
}

// Initializes a new popup
func NewPopup(text string, popupType string, ttl time.Time) *popup {
	return &popup{
		Block:         ui.NewBlock(),
		processedText: make([]string, 0),
		text:          text,
		display:       false,
		ttl:           ttl,
		popupType:     popupType,
	}
}

// Processes popup text into lines within limits and center aligning
func (popup *popup) ProcessText() {

	// Max width of popup (Extended if a word larger than width exists)
	maxWidth := 15

	// Split text into words
	words := strings.Fields(popup.text)
	curWidth := 0

	lines := make([]string, 0)
	line := ""

	// Adjust max width
	for _, word := range words {
		if len(word) > maxWidth {
			maxWidth = len(word)
		}
	}

	// Group words into lines
	for _, word := range words {

		if curWidth == 0 {
			line = word
			curWidth = len(word)
		} else if curWidth+len(word)+1 <= maxWidth {
			line = line + " " + word
			curWidth = curWidth + len(word) + 1
		} else {
			lines = append(lines, line)
			line = word
			curWidth = len(word)
		}
	}

	lines = append(lines, line)

	// Center align by padding with spaces
	for _, line := range lines {
		emptySpace := maxWidth - len(line)
		frontSpace := emptySpace / 2
		backSpace := emptySpace - frontSpace

		popup.processedText = append(
			popup.processedText,
			spaceString(frontSpace)+line+spaceString(backSpace),
		)
	}
}

// Return give number of spaces as string
func spaceString(num int) string {

	spaceString := ""

	for i := 0; i < num; i++ {
		spaceString = spaceString + " "
	}

	return spaceString
}

// Determine the size of processed
func (popup *popup) Dimensions() (int, int) {

	maxWidth := 0

	for _, line := range popup.processedText {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	return maxWidth, len(popup.processedText)
}

// Determine if popup can be rendered with respect to terminal dimensions and
// render coordinates and TTL
func (popup *popup) SetSize(x1 int, y1 int, x2 int, y2 int, width int, height int) bool {

	if time.Now().Before(popup.ttl) {
		if x1 > 0 && x2 < width && y1 > 0 && y2 < height {
			popup.SetRect(x1, y1, x2, y2)
			popup.display = true
			return true
		} else {
			popup.display = false
			return false
		}
	} else {
		popup.display = false
		return false
	}
}

// Render popup
func (popup *popup) Draw(buf *ui.Buffer) {

	// Render border
	popup.drawBorder(buf)

	for y, row := range popup.processedText {

		// Check for a rebalance popup
		if popup.popupType == "rebalance" || popup.popupType == "warning" {
			buf.SetString(
				row,
				ui.NewStyle(ui.ColorRed),
				image.Pt(
					popup.Inner.Min.X, popup.Inner.Min.Y+y,
				),
			)
		} else {
			buf.SetString(
				row,
				ui.NewStyle(ui.ColorBlue),
				image.Pt(
					popup.Inner.Min.X, popup.Inner.Min.Y+y,
				),
			)
		}
	}
}

// Render border of the popup
func (popup *popup) drawBorder(buf *ui.Buffer) {

	var verticalCell ui.Cell
	var horizontalCell ui.Cell

	// Check for rebalance popup
	if popup.popupType == "rebalance" || popup.popupType == "warning" {
		verticalCell = ui.Cell{
			Rune:  ui.VERTICAL_LINE,
			Style: ui.NewStyle(ui.ColorRed),
		}
		horizontalCell = ui.Cell{
			Rune:  ui.HORIZONTAL_LINE,
			Style: ui.NewStyle(ui.ColorRed),
		}
	} else {
		verticalCell = ui.Cell{
			Rune:  ui.VERTICAL_LINE,
			Style: ui.NewStyle(ui.ColorBlue),
		}
		horizontalCell = ui.Cell{
			Rune:  ui.HORIZONTAL_LINE,
			Style: ui.NewStyle(ui.ColorBlue),
		}
	}

	// Render top edge
	buf.Fill(
		horizontalCell,
		image.Rect(popup.Min.X, popup.Min.Y, popup.Max.X, popup.Min.Y+1),
	)

	// Render bottom edge
	buf.Fill(
		horizontalCell,
		image.Rect(popup.Min.X, popup.Max.Y-1, popup.Max.X, popup.Max.Y),
	)

	// Render left edge
	buf.Fill(
		verticalCell,
		image.Rect(popup.Min.X, popup.Min.Y, popup.Min.X+1, popup.Max.Y),
	)

	// Render right edge
	buf.Fill(
		verticalCell,
		image.Rect(popup.Max.X-1, popup.Min.Y, popup.Max.X, popup.Max.Y),
	)

	var style ui.Style

	// Check for rebalance popup
	if popup.popupType == "rebalance" || popup.popupType == "warning" {
		style = ui.NewStyle(ui.ColorRed)
	} else {
		style = ui.NewStyle(ui.ColorBlue)
	}

	// Render top left corner
	buf.SetCell(
		ui.Cell{
			Rune:  ui.TOP_LEFT,
			Style: style,
		},
		popup.Min,
	)

	// Render top right corner
	buf.SetCell(
		ui.Cell{
			Rune:  ui.TOP_RIGHT,
			Style: style,
		},
		image.Pt(popup.Max.X-1, popup.Min.Y),
	)

	// Render bottom left corner
	buf.SetCell(
		ui.Cell{
			Rune:  ui.BOTTOM_LEFT,
			Style: style,
		},
		image.Pt(popup.Min.X, popup.Max.Y-1),
	)

	// Render bottom right corner
	buf.SetCell(
		ui.Cell{
			Rune:  ui.BOTTOM_RIGHT,
			Style: style,
		},
		popup.Max.Sub(image.Pt(1, 1)),
	)
}
