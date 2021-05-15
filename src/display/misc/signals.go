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
	"syscall"
)

var allSignals = []string{
	"SIGABRT",
	"SIGALRM",
	"SIGBUS",
	"SIGCHLD",
	"SIGCLD",
	"SIGCONT",
	"SIGFPE",
	"SIGHUP",
	"SIGILL",
	"SIGINT",
	"SIGIO",
	"SIGIOT",
	"SIGKILL",
	"SIGPIPE",
	"SIGPOLL",
	"SIGPROF",
	"SIGPWR",
	"SIGQUIT",
	"SIGSEGV",
	"SIGSTKFLT",
	"SIGSTOP",
	"SIGSYS",
	"SIGTERM",
	"SIGTRAP",
	"SIGTSTP",
	"SIGTTIN",
	"SIGTTOU",
	"SIGUNUSED",
	"SIGURG",
	"SIGUSR1",
	"SIGUSR2",
	"SIGVTALRM",
	"SIGWINCH",
	"SIGXCPU",
	"SIGXFSZ",
}

var signalMap = map[string]syscall.Signal {
	"SIGABRT":syscall.SIGABRT,
	"SIGALRM":syscall.SIGALRM,
	"SIGBUS":syscall.SIGBUS,
	"SIGCHLD":syscall.SIGCHLD,
	"SIGCLD":syscall.SIGCLD,
	"SIGCONT":syscall.SIGCONT,
	"SIGFPE":syscall.SIGFPE,
	"SIGHUP":syscall.SIGHUP,
	"SIGILL":syscall.SIGILL,
	"SIGINT":syscall.SIGINT,
	"SIGIO":syscall.SIGIO,
	"SIGIOT":syscall.SIGIOT,
	"SIGKILL":syscall.SIGKILL,
	"SIGPIPE":syscall.SIGPIPE,
	"SIGPOLL":syscall.SIGPOLL,
	"SIGPROF":syscall.SIGPROF,
	"SIGPWR":syscall.SIGPWR,
	"SIGQUIT":syscall.SIGQUIT,
	"SIGSEGV":syscall.SIGSEGV,
	"SIGSTKFLT":syscall.SIGSTKFLT,
	"SIGSTOP":syscall.SIGSTOP,
	"SIGSYS":syscall.SIGSYS,
	"SIGTERM":syscall.SIGTERM,
	"SIGTRAP":syscall.SIGTRAP,
	"SIGTSTP":syscall.SIGTSTP,
	"SIGTTIN":syscall.SIGTTIN,
	"SIGTTOU":syscall.SIGTTOU,
	"SIGUNUSED":syscall.SIGUNUSED,
	"SIGURG":syscall.SIGURG,
	"SIGUSR1":syscall.SIGUSR1,
	"SIGUSR2":syscall.SIGUSR2,
	"SIGVTALRM":syscall.SIGVTALRM,
	"SIGWINCH":syscall.SIGWINCH,
	"SIGXCPU":syscall.SIGXCPU,
	"SIGXFSZ":syscall.SIGXFSZ,
}

// SignalList is a wrapper widget around a List
// meant to display error messages if any
type SignalList struct {
	*widgets.List
}

// NewErrorBox is a constructor for the SignalList type
func NewSignalList() *SignalList {
	return &SignalList{
		List: widgets.NewList(),
	}
}

func (sigList *SignalList) SignalFromRow(rowIndex int) syscall.Signal {
	return signalMap[sigList.Rows[rowIndex]]
}

func (sigList *SignalList) SelectedSignal() syscall.Signal {
	return signalMap[sigList.Rows[sigList.SelectedRow]]
}

// Draw puts the required text into the widget
func (sigList *SignalList) Draw(buf *ui.Buffer) {
	sigList.List.Title = " Select signal "

	sigList.List.Rows = allSignals
	sigList.List.TextStyle = ui.NewStyle(ui.ColorYellow)
	sigList.List.WrapText = false
	sigList.List.Draw(buf)
}
