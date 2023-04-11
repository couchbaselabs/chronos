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
	"os"
	"sync"
	"time"

	log "github.com/couchbase/clog"
	ui "github.com/gizak/termui/v3"
)

// Colors for different event types
// Percent defined to avoid warning in fmt.Sprintf()
const (
	colorCyan2     ui.Color = 50
	colorSeaGreen1 ui.Color = 84
	colorRed3      ui.Color = 160
	percent        string   = "%"
)

// Assign colors for event types
var eventColors = map[string]ui.Color{
	"Below Threshold": colorCyan2,
	"Above Threshold": colorRed3,
	"Sudden Change":   colorSeaGreen1,
}

// Widget to display a list of alerts
// Each row can be displayed on more than one line
// Allows printing of reports for any alert
type EventDisplay struct {
	*ui.Block

	// List of alerts to be displayed
	Events []*Event

	// Heading for the widget
	header string

	// Cursor position
	SelectedRow int

	// Row currently displayed on the first line
	TopRow int

	// Toggle indicating if widget is currently selected by the user
	selected bool

	// Array to track the number of lines used to display each alert
	rowSize []int

	// Lock for the list of alerts
	EventLock sync.RWMutex
}

// Struct to hold all the information for one alert
type Event struct {

	// Holds data within data padding defined by user
	Data []float64

	// Holds timestamps for each data point
	DataTimes []time.Time

	// Holds timestamps everytime the alert was triggered
	AlertTimes []time.Time

	// Name of the node the alert corresponds to
	Node string

	// Name of the stat the alert corresponds to
	Stat string

	// Type of the alert triggered
	EventType string

	// The threshold for the alert
	Threshold float64

	// The value of data at the time of the latest alert
	ThresholdData float64

	// The amount of change at the time of the latest alert
	ThresholdChange float64

	// The amount of time with which the change is calculated
	ThresholdTime int

	// Description of alert to be displayed
	Description string

	// The time of creation of the alert
	FirstTriggered time.Time

	// The time the alert last triggered
	LastTriggered time.Time

	// The time from which the alert has been collecting data
	DataStart time.Time

	// Toggle to indicate alert has expired
	Stale bool

	// Counter for number of times the alert has triggered
	NumTimes int

	// Toggle to indicate alert data is full
	DataFilled bool

	// Toggle to indicate alert node is no longer in cluster
	Deprecated bool
}

// Initializes a new event display
func NewEventDisplay() *EventDisplay {
	return &EventDisplay{
		Block:       ui.NewBlock(),
		Events:      make([]*Event, 0),
		header:      "Alerts",
		SelectedRow: 0,
		TopRow:      0,
		selected:    false,
		rowSize:     make([]int, 0),
		EventLock:   sync.RWMutex{},
	}
}

// Initializes a new event
func NewEvent(node string, stat string, eventType string,
	thresholdData float64, threshold float64) *Event {
	return &Event{
		Data:           make([]float64, 0),
		DataTimes:      make([]time.Time, 0),
		AlertTimes:     make([]time.Time, 0),
		Node:           node,
		Stat:           stat,
		EventType:      eventType,
		Threshold:      threshold,
		ThresholdData:  thresholdData,
		FirstTriggered: time.Now(),
		LastTriggered:  time.Now(),
		Stale:          false,
		NumTimes:       1,
		Deprecated:     false,
	}
}

// Deep copies an alert
// Used while generating reports
func CopyEvent(event *Event) *Event {

	eventData := make([]float64, 0)
	eventDataTimes := make([]time.Time, 0)
	eventAlertTimes := make([]time.Time, 0)

	eventData = append(eventData, event.Data...)
	eventDataTimes = append(eventDataTimes, event.DataTimes...)
	eventAlertTimes = append(eventAlertTimes, event.AlertTimes...)

	return &Event{
		Data:            eventData,
		DataTimes:       eventDataTimes,
		AlertTimes:      eventAlertTimes,
		Node:            event.Node,
		Stat:            event.Stat,
		EventType:       event.EventType,
		Threshold:       event.Threshold,
		ThresholdData:   event.ThresholdData,
		ThresholdChange: event.ThresholdChange,
		ThresholdTime:   event.ThresholdTime,
		Description:     event.Description,
		FirstTriggered:  event.FirstTriggered,
		LastTriggered:   event.LastTriggered,
		DataStart:       event.DataStart,
		Stale:           event.Stale,
		NumTimes:        event.NumTimes,
		DataFilled:      event.DataFilled,
		Deprecated:      event.Deprecated,
	}
}

// Handler to generate a report for an event
func MakeReport(event *Event, path string) {

	filePath := fmt.Sprintf(
		path+"Alert Report - %s.txt",
		event.FirstTriggered.Format("2006-01-02 15:04:05.000000"),
	)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("event_display: Failed to create file: %v", err)
		return
	}

	fileInfo := ReportText(event)

	_, err2 := file.WriteString(fileInfo)
	if err2 != nil {
		log.Printf("event_display: Error writing to file: %v", err)
		return
	}
	file.Close()
}

// Handler to make the report text as a string
func ReportText(event *Event) string {

	fileInfo := fmt.Sprintf(
		"Node - %s\nStat - %s\n\n", event.Node, event.Stat,
	)

	switch event.EventType {
	case "Sudden Change":
		if event.NumTimes == 1 {
			fileInfo = fileInfo + fmt.Sprintf(
				"Stat changed by more than the threshold limit of %.2f%s at %s."+
					" This change occured over %d second(s).\n\n",
				event.Threshold*100, percent,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
				event.ThresholdTime,
			)
		} else {
			fileInfo = fileInfo + fmt.Sprintf("Stat changed by more than the "+
				"threshold limit of %.2f%s at %s. This change occured over %d "+
				"second(s).\nSimilar changes occured %d times with the last one"+
				" occuring at %s.\n\n",
				event.Threshold*100, percent,
				event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.ThresholdTime, event.NumTimes,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	case "Above Threshold":
		if event.NumTimes == 1 {
			fileInfo = fileInfo + fmt.Sprintf(
				"Stat exceeded threshold limit of %f at %s.\n\n",
				event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		} else {
			fileInfo = fileInfo + fmt.Sprintf(
				"Stat exceeded threshold limit of %f at %s.\nSimilarly, "+
					"the stat exceeded threshold limit %d times with the last"+
					" one occuring at %s.\n\n", event.Threshold,
				event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.NumTimes,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	case "Below Threshold":
		if event.NumTimes == 1 {
			fileInfo = fileInfo + fmt.Sprintf(
				"Stat dropped below threshold limit of %f at %s.\n\n",
				event.Threshold,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		} else {
			fileInfo = fileInfo + fmt.Sprintf(
				"Stat dropped below threshold limit of %f at %s.\n"+
					"Similarly, the stat was below the threshold limit"+
					" %d times with the last one occuring at %s.\n\n",
				event.Threshold,
				event.FirstTriggered.Format("2006-01-02 15:04:05"),
				event.NumTimes,
				event.LastTriggered.Format("2006-01-02 15:04:05"),
			)
		}
	}

	if event.Deprecated {
		fileInfo = fileInfo + "Node corresponding to alert was removed from the cluster.\n\n"
	}

	fileInfo = fileInfo + fmt.Sprintf(
		"Data collected from %s to %s\n\n",
		event.DataTimes[0].Format("2006-01-02 15:04:05"),
		event.DataTimes[len(event.DataTimes)-1].Format("2006-01-02 15:04:05"),
	)

	prevTime := event.DataStart
	var curTime time.Time

	for i, j := 0, 0; i < len(event.DataTimes); i++ {

		curTime = event.DataTimes[i]

		if i == 0 {
			if CompareTimes(prevTime, curTime) {
				fileInfo = fileInfo + fmt.Sprintf(
					"%s - %f", curTime.Format("2006-01-02 15:04:05"),
					event.Data[i],
				)
			} else {
				fileInfo = fileInfo + fmt.Sprintf(
					"No data recieved from server before %s\n",
					curTime.Format("2006-01-02 15:04:05"),
				)
				fileInfo = fileInfo + fmt.Sprintf(
					"%s - %f", curTime.Format("2006-01-02 15:04:05"),
					event.Data[i],
				)
			}
		} else {
			if CompareTimes(prevTime.Add(time.Second), curTime) {
				fileInfo = fileInfo + fmt.Sprintf(
					"%s - %f", curTime.Format("2006-01-02 15:04:05"),
					event.Data[i],
				)
			} else {
				fileInfo = fileInfo + fmt.Sprintf(
					"No data recieved from server between %s and %s\n",
					prevTime.Format("2006-01-02 15:04:05"),
					curTime.Format("2006-01-02 15:04:05"),
				)
				fileInfo = fileInfo + fmt.Sprintf(
					"%s - %f", curTime.Format("2006-01-02 15:04:05"),
					event.Data[i],
				)
			}
		}

		if j < len(event.AlertTimes) {
			if CompareTimes(event.AlertTimes[j], curTime) {
				fileInfo = fileInfo + " ALERT"
				j++
			}
		}

		fileInfo = fileInfo + "\n"
		prevTime = curTime
	}

	return fileInfo
}

// Check if two given times are within 500 milliseconds of each other
func CompareTimes(t1 time.Time, t2 time.Time) bool {

	t1Min := t1.Add(-time.Millisecond * time.Duration(500))
	t1Max := t1.Add(time.Millisecond * time.Duration(500))

	if t2.Before(t1Max) && t2.After(t1Min) {
		return true
	} else {
		return false
	}
}

// Render widget
func (display *EventDisplay) Draw(buf *ui.Buffer) {

	if display.selected {
		display.BorderStyle = ui.NewStyle(ui.ColorWhite)
	} else {
		display.BorderStyle = ui.NewStyle(8)
	}
	display.Block.Draw(buf)

	// Horizontal padding of header from the left edge
	paddingHeader := 10

	// Vertical padding of rows from the header
	paddingRow := 4

	// Display style of the header
	styleHeader := ui.NewStyle(
		display.BorderStyle.Fg, ui.ColorClear, ui.ModifierBold,
	)

	// Render header
	buf.SetString(
		display.header, styleHeader,
		image.Pt(display.Inner.Min.X+paddingHeader, display.Inner.Min.Y+1),
	)

	display.EventLock.RLock()
	display.rowSize = make([]int, len(display.Events))

	// Keep track of when space is not available to render the row
	continueRender := true

	// Loop to render as many rows as possible within the bounds of the widget
	for rowNum, usedSpace := 0, 3; rowNum < len(display.Events); rowNum++ {

		event := display.Events[rowNum]

		var eventCells []ui.Cell

		// Check if current row is selected
		if rowNum == display.SelectedRow && display.selected {
			// Parse row text into cells
			eventCells = ui.ParseStyles(
				event.Description,
				ui.NewStyle(
					ui.ColorBlack, eventColors[event.EventType], ui.ModifierClear,
				),
			)
		} else {
			// Parse row text into cells
			eventCells = ui.ParseStyles(
				event.Description,
				ui.NewStyle(
					eventColors[event.EventType], ui.ColorClear, ui.ModifierClear,
				),
			)
		}

		// Add padding for the rows
		eventCells = ui.WrapCells(
			eventCells, uint(display.Inner.Dx()-2*paddingRow),
		)

		// Split cells into multiple lines
		eventCellRows := ui.SplitCells(eventCells, '\n')

		// Render each cell if all the lines of text fit within the widget
		if len(eventCellRows) < display.Inner.Dy()-usedSpace &&
			rowNum >= display.TopRow && continueRender {
			for i, row := range eventCellRows {
				for _, cx := range ui.BuildCellWithXArray(row) {
					x, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(
						display.Inner.Min.X+paddingRow+x,
						display.Inner.Min.Y+usedSpace+i,
					))
				}
			}

			// Update the size of the text
			display.rowSize[rowNum] = len(eventCellRows)

			// Track the current line number
			usedSpace = usedSpace + display.rowSize[rowNum]
		} else {
			// Update number of rows taken to render even if not rendering
			display.rowSize[rowNum] = len(eventCellRows)

			if rowNum >= display.TopRow {
				// Stop rendering from this row
				continueRender = false
			}
		}
	}
	display.EventLock.RUnlock()
}

// Handler function for scroll up
func (display *EventDisplay) ScrollUp() {

	display.SelectedRow--

	display.CalcPos()
}

// Handler function for scroll down
func (display *EventDisplay) ScrollDown() {

	display.SelectedRow++

	display.CalcPos()
}

// Handler function to ensure cursor is never out of bounds
func (display *EventDisplay) CalcPos() {

	if display.SelectedRow < 0 {
		display.SelectedRow = 0
	}

	if display.SelectedRow < display.TopRow {
		display.TopRow = display.SelectedRow
	}

	display.EventLock.RLock()
	if display.SelectedRow > len(display.Events)-1 {
		display.SelectedRow = len(display.Events) - 1
	}
	display.EventLock.RUnlock()

	if display.SelectedRow >= display.TopRow+display.RowsOnDisplay() {
		space := display.Inner.Dy() - 4

		for i := display.SelectedRow; i >= 0; i-- {
			space = space - display.rowSize[i]
			if space < 0 {
				if i == display.SelectedRow {
					display.TopRow = i
				} else {
					display.TopRow = i + 1
				}
				break
			}
		}
	}
}

// Calculate number of alerts currently on display
func (display *EventDisplay) RowsOnDisplay() int {

	space := display.Inner.Dy() - 4
	rows := 0

	for i := display.TopRow; i < len(display.rowSize); i++ {

		space = space - display.rowSize[i]

		if space > 0 {
			rows = rows + 1
		} else {
			break
		}
	}

	return rows
}

// Handler to add a new event
func (display *EventDisplay) AddEvent(event *Event) {
	display.EventLock.Lock()
	display.Events = append(display.Events, event)
	display.EventLock.Unlock()
}

// Handler to reset cursor
func (display *EventDisplay) ResetSelect() {
	display.SelectedRow = 0
	display.TopRow = 0
}

// Handler to report an event
func (display *EventDisplay) ReportEvent(path string) {

	display.EventLock.RLock()
	var event *Event

	// Check if event exists
	// Can go out of bounds immediately after an event expires
	if display.SelectedRow >= 0 &&
		display.SelectedRow <= len(display.Events)-1 {
		// Copy event to reduce latency of report generation
		event = CopyEvent(display.Events[display.SelectedRow])
	}

	display.EventLock.RUnlock()

	if event != nil {
		// Generate report in a separate routine
		go MakeReport(event, path)
	}
}

// Handler function to indicate if cursor is on widget
func (table *EventDisplay) ToggleTableSelect() {
	table.selected = !table.selected
}

// Handler function for mouse click
func (display *EventDisplay) HandleClick(x int, y int) {
	x = x - display.Min.X
	y = y - display.Min.Y
	if (x > 0 && x <= display.Inner.Dx()) &&
		(y > 0 && y <= display.Inner.Dy()) {
		row := (display.TopRow + y) - 4

		for i, size := range display.rowSize {
			if row-size < 0 {
				display.SelectedRow = i
				break
			} else {
				row = row - size
			}
		}
		display.CalcPos()
	}
}

func (table *EventDisplay) Contains(x int, y int) bool {
	if x > table.Min.X && x <= table.Max.X &&
		y > table.Min.Y && y <= table.Max.Y {
		return true
	} else {
		return false
	}
}
