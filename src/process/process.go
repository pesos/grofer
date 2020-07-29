package process

import (
	proc "github.com/shirou/gopsutil/process"
)

// Process type contains as fields all the information extracted from the kernel.
type Process struct {
	Proc           *proc.Process
	Background     bool
	Foreground     bool
	IsRunning      bool
	CPUPercent     float64
	Children       []*proc.Process
	CreateTime     int64
	Gids           []int32
	MemoryInfo     *proc.MemoryInfoStat
	MemoryPercent  float32
	Name           string
	Nice           int32
	NumCtxSwitches *proc.NumCtxSwitchesStat
	NumThreads     int32
	PageFault      *proc.PageFaultsStat
	Status         string
	Exe            string
	CPUAffinity    []int32
}

// InitAllProcs initialises the set of currently running processes in the system.
func InitAllProcs() (map[int32]*Process, error) {
	var processes map[int32]*Process = make(map[int32]*Process)
	pids, err := proc.Processes()

	if err != nil {
		return processes, err
	}

	for _, proc := range pids {
		tempProc := &Process{Proc: proc}
		processes[proc.Pid] = tempProc
	}
	return processes, nil
}
