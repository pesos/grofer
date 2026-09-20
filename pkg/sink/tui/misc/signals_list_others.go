//go:build android || aix || dragonfly || freebsd || (js && wasm) || linux || nacl || netbsd || plan9 || windows

package misc

import "syscall"

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
