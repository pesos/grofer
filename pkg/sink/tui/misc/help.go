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

// HelpMenu is a wrapper widget around a List meant
// to display the help menu for a command. HelpMenu
// implements the ui.Drawable interface.
type HelpMenu struct {
	*widgets.List
	keybindings []string
}

// NewHelpMenu is a constructor for the HelpMenu type.
func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		List:        widgets.NewList(),
		keybindings: getDefaultHelpKeybinding(),
	}
}

// ForCommand sets the keybindings to be displayed as part of the help
// for a specific command and returns the modified HelpMenu.
func (help *HelpMenu) ForCommand(command HelpKeybindingType) *HelpMenu {
	help.keybindings = getHelpKeybindingsForCommand(command)
	return help
}

// Resize resizes the widget based on specified width
// and height
func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range help.keybindings {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := len(help.keybindings) + 3
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

	help.List.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget.
func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.List.Title = " Keybindings "

	help.List.Rows = help.keybindings
	help.List.TextStyle = ui.NewStyle(ui.ColorYellow)
	help.List.WrapText = false
	help.List.Draw(buf)
}

// ensure interface compliance.
var _ ui.Drawable = (*HelpMenu)(nil)
