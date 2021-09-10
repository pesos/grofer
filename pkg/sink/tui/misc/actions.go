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
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

var allActions = [][]string{
	{
		"PAUSE",
	},
	{
		"UNPAUSE",
	},
	{
		"RESTART",
	},
	{
		"STOP",
	},
	{
		"KILL",
	},
	{
		"REMOVE",
	},
}

const actionNameIdx = 0

// ActionTable is a wrapper widget around a Table
// meant to display error messages if any
type ActionTable struct {
	*viz.Table
}

// NewActionTable is a constructor for the ActionTable type
func NewActionTable() *ActionTable {
	actionTable := &ActionTable{
		Table: viz.NewTable(),
	}
	actionTable.Table.Title = " Select Action "
	actionTable.Table.Header = []string{"Action"}
	actionTable.Table.Rows = allActions
	actionTable.Table.ColResizer = func() {
		x := actionTable.Table.Inner.Dx()
		actionTable.Table.ColWidths = []int{x}
	}
	actionTable.Table.ShowCursor = true
	actionTable.Table.CursorColor = ui.ColorCyan
	actionTable.Table.BorderStyle.Fg = ui.ColorCyan
	return actionTable
}

// SelectedAction returns an action as string from the selected row of the action table
func (actionTable *ActionTable) SelectedAction() string {
	return actionTable.Rows[actionTable.SelectedRow][actionNameIdx]
}

// Draw puts the required text into the widget
func (actionTable *ActionTable) Draw(buf *ui.Buffer) {
	actionTable.Table.Draw(buf)
}

// ensure interface compliance.
var _ ui.Drawable = (*ActionTable)(nil)
