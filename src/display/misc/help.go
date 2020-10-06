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
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var KEYBINDS string

type HelpMenu struct {
	*widgets.List
}

func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		List: widgets.NewList(),
	}
}

func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 53
	for _, line := range strings.Split(KEYBINDS, "\n") {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := strings.Count(KEYBINDS, "\n") + 1
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	help.Block.SetRect(x, y, textWidth+x, textHeight+y)
}

func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.Block.Draw(buf)

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, rune := range line {
			buf.SetCell(
				ui.NewCell(rune, ui.Theme.Default),
				image.Pt(help.Inner.Min.X+x, help.Inner.Min.Y+y-1),
			)
		}
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
