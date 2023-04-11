//  Copyright 2023-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included
//  in the file licenses/BSL-Couchbase.txt.  As of the Change Date specified
//  in that file, in accordance with the Business Source License, use of this
//  software will be governed by the Apache License, Version 2.0, included in
//  the file licenses/APL2.txt.

package widgets

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

// Widget to display popups
// Works on top of the grid with independent sizing
type PopupManager struct {

	// List of popups to be displayed
	Popups []*popup

	// Width of the terminal window
	Width int

	// Height of the terminal window
	Height int
}

// Initializes a new popup manager
func NewPopupManager() *PopupManager {
	return &PopupManager{
		Popups: make([]*popup, 0),
	}
}

// Handler to add a new 'add node' popup
func (manager *PopupManager) AddNodePopup(node string) {
	popup := NewPopup(
		"Node "+node+" added to cluster", "nodeAdd",
		time.Now().Add(time.Second*time.Duration(5)),
	)
	popup.ProcessText()
	manager.Popups = append(manager.Popups, popup)
}

// Handler to add a new 'remove node' popup
func (manager *PopupManager) RemoveNodePopup(node string) {
	popup := NewPopup(
		"Node "+node+" removed from cluster", "nodeRemove",
		time.Now().Add(time.Second*time.Duration(5)),
	)
	popup.ProcessText()
	manager.Popups = append(manager.Popups, popup)
}

// Handler to add a new popup
func (manager *PopupManager) NewPopup(text string, popupType string, ttl time.Time) {

	// Check if a duplicate already exists
	exist := manager.checkCopy(text, popupType, ttl)

	if !exist {
		popup := NewPopup(text, popupType, ttl)
		popup.ProcessText()
		manager.Popups = append(manager.Popups, popup)
	}
}

// Check for a duplicate. Update ttl if exists
func (manager *PopupManager) checkCopy(text string, popupType string, ttl time.Time) bool {

	exist := false

	for _, popup := range manager.Popups {

		if popup.popupType == popupType && popup.text == text {
			popup.ttl = ttl
			exist = true
			break
		}
	}

	return exist
}

// Handler to render all the popups
func (manager *PopupManager) Render() {

	manager.SetSize(manager.Width, manager.Height)

	for _, popup := range manager.Popups {
		if popup.display {
			ui.Render(popup)
		}
	}
}

// Handler to determine the size of each popup
func (manager *PopupManager) SetSize(windowWidth int, windowHeight int) {

	// Updating terminal dimensions
	manager.Width = windowWidth
	manager.Height = windowHeight

	// External padding for each popup
	paddingExternal := 2

	// Internal padding for each popup
	paddingInternal := 1

	curHeight := paddingExternal
	curWidth := windowWidth - paddingExternal

	for _, popup := range manager.Popups {

		// Get dimensions of popup
		width, height := popup.Dimensions()

		// Calculate if popup fits within terminal bounds
		ok := popup.SetSize(
			curWidth-width-2*paddingInternal, curHeight,
			curWidth, curHeight+height+2*paddingInternal,
			windowWidth, windowHeight,
		)

		// Update width if popup can be rendered
		if ok {
			curWidth = curWidth - width - 2*paddingInternal - paddingExternal
		}
	}
}
