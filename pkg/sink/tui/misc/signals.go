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

var allSignals = [][]string{
	{
		"1",
		"SIGABRT",
	},
	{
		"2",
		"SIGALRM",
	},
	{
		"3",
		"SIGBUS",
	},
	{
		"4",
		"SIGCHLD",
	},
	{
		"5",
		"SIGCLD",
	},
	{
		"6",
		"SIGCONT",
	},
	{
		"7",
		"SIGFPE",
	},
	{
		"8",
		"SIGHUP",
	},
	{
		"9",
		"SIGILL",
	},
	{
		"10",
		"SIGINT",
	},
	{
		"11",
		"SIGIO",
	},
	{
		"12",
		"SIGIOT",
	},
	{
		"13",
		"SIGKILL",
	},
	{
		"14",
		"SIGPIPE",
	},
	{
		"15",
		"SIGPOLL",
	},
	{
		"16",
		"SIGPROF",
	},
	{
		"17",
		"SIGPWR",
	},
	{
		"18",
		"SIGQUIT",
	},
	{
		"19",
		"SIGSEGV",
	},
	{
		"20",
		"SIGSTKFLT",
	},
	{
		"21",
		"SIGSTOP",
	},
	{
		"22",
		"SIGSYS",
	},
	{
		"23",
		"SIGTERM",
	},
	{
		"24",
		"SIGTRAP",
	},
	{
		"25",
		"SIGTSTP",
	},
	{
		"26",
		"SIGTTIN",
	},
	{
		"27",
		"SIGTTOU",
	},
	{
		"28",
		"SIGUNUSED",
	},
	{
		"29",
		"SIGURG",
	},
	{
		"30",
		"SIGUSR1",
	},
	{
		"31",
		"SIGUSR2",
	},
	{
		"32",
		"SIGVTALRM",
	},
	{
		"33",
		"SIGWINCH",
	},
	{
		"34",
		"SIGXCPU",
	},
	{
		"35",
		"SIGXFSZ",
	},
}

var signalMap = map[string]syscall.Signal{
	"SIGABRT":   syscall.SIGABRT,
	"SIGALRM":   syscall.SIGALRM,
	"SIGBUS":    syscall.SIGBUS,
	"SIGCHLD":   syscall.SIGCHLD,
	"SIGCLD":    syscall.SIGCLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPOLL":   syscall.SIGPOLL,
	"SIGPROF":   syscall.SIGPROF,
	"SIGPWR":    syscall.SIGPWR,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTKFLT": syscall.SIGSTKFLT,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGUNUSED": syscall.SIGUNUSED,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ,
}

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
		sigTable.Table.ColWidths = []int{
			4,
			10,
		}
	}
	sigTable.Table.ShowCursor = true
	sigTable.Table.CursorColor = ui.ColorCyan
	sigTable.Table.BorderStyle.Fg = ui.ColorCyan
	return sigTable
}

func (sigTable *SignalTable) SignalFromRow(rowIndex int) syscall.Signal {
	return signalMap[sigTable.Rows[rowIndex][sigNameIdx]]
}

func (sigTable *SignalTable) SelectedSignal() syscall.Signal {
	return signalMap[sigTable.Rows[sigTable.SelectedRow][sigNameIdx]]
}

// Draw puts the required text into the widget
func (sigTable *SignalTable) Draw(buf *ui.Buffer) {
	sigTable.Table.Draw(buf)
}

// ensure interface compliance.
var _ ui.Drawable = (*SignalTable)(nil)
