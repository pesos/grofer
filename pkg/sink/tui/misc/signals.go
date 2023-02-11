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
	"syscall"

	ui "github.com/gizak/termui/v3"
	viz "github.com/pesos/grofer/pkg/utils/visualization"
)

const sigNameIdx int = 1

// SignalTable is a wrapper widget around a Table
// meant to display error messages if any
type SignalTable struct {
	*viz.Table
}

// NewSignalTable is a constructor for the SignalTable type
func NewSignalTable() *SignalTable {
	sigTable := &SignalTable{
		Table: viz.NewTable(),
	}
	sigTable.Table.Title = " Select Signal "
	sigTable.Table.Header = []string{"ID", "Signal"}
	sigTable.Table.Rows = allSignals
	sigTable.Table.ColWidths = []int{4, 10}
	sigTable.Table.ColResizer = func() {
		x := sigTable.Table.Inner.Dx()
		sigTable.Table.ColWidths = []int{
			3 * x / 10,
			7 * x / 10,
		}
	}
	sigTable.Table.ShowCursor = true
	sigTable.Table.CursorColor = ui.ColorCyan
	sigTable.Table.BorderStyle.Fg = ui.ColorCyan
	return sigTable
}

// SignalFromRow returns the symbol at a given row index
func (sigTable *SignalTable) SignalFromRow(rowIndex int) syscall.Signal {
	return signalMap[sigTable.Rows[rowIndex][sigNameIdx]]
}

// SelectedSignal returns the signal at the currently selected row index
func (sigTable *SignalTable) SelectedSignal() syscall.Signal {
	return signalMap[sigTable.Rows[sigTable.SelectedRow][sigNameIdx]]
}

// Draw puts the required text into the widget
func (sigTable *SignalTable) Draw(buf *ui.Buffer) {
	sigTable.Table.Draw(buf)
}

// ensure interface compliance.
var _ ui.Drawable = (*SignalTable)(nil)
