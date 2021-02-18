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
package process

import (
	proc "github.com/shirou/gopsutil/process"
)

// Process type contains as fields all the information extracted from the kernel.
type Process struct {
	Proc           *proc.Process
	MemoryInfo     *proc.MemoryInfoStat
	PageFault      *proc.PageFaultsStat
	NumCtxSwitches *proc.NumCtxSwitchesStat
	Exe            string
	Name           string
	Status         string
	Children       []*proc.Process
	Gids           []int32
	CPUAffinity    []int32
	CreateTime     int64
	CPUPercent     float64
	Nice           int32
	NumThreads     int32
	MemoryPercent  float32
	IsRunning      bool
	Foreground     bool
	Background     bool
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

func NewProcess(pid int32) (*Process, error) {
	process, err := proc.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	newProcess := &Process{Proc: process}
	return newProcess, nil
}
