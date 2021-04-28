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

package help

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var KEYBINDS []string

type HelpMenu struct {
	*widgets.List
}

func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		List: widgets.NewList(),
	}
}

func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range KEYBINDS {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := len(KEYBINDS) + 3
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

func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.List.Title = " Keybindings "

	help.List.Rows = KEYBINDS
	help.List.TextStyle = ui.NewStyle(ui.ColorYellow)
	help.List.WrapText = false
	help.List.Draw(buf)
}

func SelectHelpMenu(page string) {
	switch page {
	case "proc":
		KEYBINDS = PROC_KEYBINDS
	case "proc_pid":
		KEYBINDS = PROC_PID_KEYBINDS
	case "main":
		KEYBINDS = MAIN_KEYBINDS
	case "cont":
		KEYBINDS = CONT_KEYBINDS
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
