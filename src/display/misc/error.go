/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package misc

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// ErrorString contains the error message
// to display
var ErrorString string

// ErrorBox is a wrapper widget around a List
// meant to display error messages if any
type ErrorBox struct {
	*widgets.List
}

// NewErrorBox is a constructor for the ErrorBox type
func NewErrorBox() *ErrorBox {
	return &ErrorBox{
		List: widgets.NewList(),
	}
}

// Resize resizes the widget based on specified width
// and height
func (errBox *ErrorBox) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range errorKeybindings {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := len(keybindings) + 3
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2
	if x < 0 {
		x = 0
		textWidth = termWidth
	}
	if y < 0 {
		y = 0
		textHeight = termHeight
	}

	errBox.List.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (errBox *ErrorBox) Draw(buf *ui.Buffer) {
	errBox.List.Title = " Error "

	errBox.List.Rows = []string{ErrorString, ""}
	errBox.List.Rows = append(errBox.List.Rows, errorKeybindings...)
	errBox.List.TextStyle = ui.NewStyle(ui.ColorYellow)
	errBox.List.WrapText = false
	errBox.List.Draw(buf)
}

// SetErrorString sets the error string to be displayed
func SetErrorString(errStr string) {
	ErrorString = errStr
}
