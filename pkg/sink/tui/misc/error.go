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
	vz "github.com/pesos/grofer/pkg/utils/visualization"
)

// ErrorBox is a wrapper widget around a List
// meant to display error messages if any. It
// implements the ui.Drawable interface.
type ErrorBox struct {
	*vz.Table
	errorMessage string
	err          error
	keybindings  [][]string
}

// NewErrorBox is a constructor for the ErrorBox type.
func NewErrorBox() *ErrorBox {
	return &ErrorBox{
		Table:       vz.NewTable(),
		keybindings: getErrorKeybindings(),
	}
}

// Resize resizes the widget based on specified width
// and height.
func (errBox *ErrorBox) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range errBox.Rows {
		if textWidth < len(line[0]) {
			textWidth = len(line[0]) + 2
		}
	}
	textHeight := len(errBox.keybindings) + 5
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

	errBox.Table.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget.
func (errBox *ErrorBox) Draw(buf *ui.Buffer) {
	errBox.Table.Title = " Error "
	errBox.Table.Header = []string{errBox.errorMessage}
	errBox.Table.Rows = [][]string{{errBox.err.Error()}}
	errBox.Table.Rows = append(errBox.Table.Rows, errBox.keybindings...)
	errBox.Table.BorderStyle.Fg = ui.ColorCyan
	errBox.Table.BorderStyle.Bg = ui.ColorClear
	errBox.Table.ColResizer = func() {
		x := errBox.Table.Inner.Dx()
		errBox.Table.ColWidths = []int{x}
	}
	errBox.Table.Draw(buf)
}

// SetErrorString sets the error string to be displayed.
func (errBox *ErrorBox) SetErrorString(errStr string, err error) {
	errBox.errorMessage = errStr
	errBox.err = err
}

// ensure interface compliance.
var _ ui.Drawable = (*ErrorBox)(nil)
