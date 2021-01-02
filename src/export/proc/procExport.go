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

package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	procInfo "github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
)

type contextSwitches struct {
	Voluntary   int64 `json:"voluntary"`
	InVoluntary int64 `json:"involuntary"`
}

type pageFaults struct {
	Minor uint64 `json:"minor"`
	Major uint64 `json:"major"`
}

type memPidStats struct {
	RSS   uint64 `json:"RSS"`
	Data  uint64 `json:"Data"`
	Stack uint64 `json:"Stack"`
	Swap  uint64 `json:"Swap"`
}

type ProcDetails struct {
	Name              string `json:"Name"`
	Command           string `json:"Command"`
	Status            string `json:"Status"`
	Background        string `json:"Background"`
	Running           string `json:"Running"`
	CreationTime      string `json:"CreationTime"`
	NiceValue         string `json:"NiceValue"`
	ThreadCount       string `json:"ThreadCount"`
	ChildProcessCount string `json:"ChildProcessCount"`
}

type PidStats struct {
	Cpu         float64         `json:"cpu"`
	Mem         float64         `json:"mem"`
	PidDetails  ProcDetails     `json:"pid"`
	CtxSwitches contextSwitches `json:"ctxSwitches"`
	PageFaults  pageFaults      `json:"pageFaults"`
	MemStats    memPidStats     `json:"memStats"`
	ChildProcs  map[int]string  `json:"childProcs"`
}

var statusMap map[string]string = map[string]string{
	"R": "Running",
	"S": "Sleep",
	"Z": "Zombie",
	"T": "Stop",
	"I": "Idle",
	"W": "Wait",
	"L": "Lock",
}

func getPidDataJSON(pid int32) (PidStats, error) {
	proc, err := procInfo.NewProcess(pid)
	if err != nil {
		return PidStats{}, err
	}
	proc.UpdateProcInfo()

	pidDetails := ProcDetails{
		Name:              proc.Name,
		Command:           proc.Exe,
		Status:            statusMap[proc.Status],
		Background:        strconv.FormatBool(proc.Background),
		Running:           strconv.FormatBool(proc.IsRunning),
		CreationTime:      strconv.FormatInt(proc.CreateTime, 10),
		NiceValue:         strconv.Itoa(int(proc.Nice)),
		ThreadCount:       strconv.Itoa(int(proc.NumThreads)),
		ChildProcessCount: strconv.Itoa(len(proc.Children)),
	}

	ctxSwitches := contextSwitches{
		Voluntary:   proc.NumCtxSwitches.Voluntary,
		InVoluntary: proc.NumCtxSwitches.Involuntary,
	}

	pgFaults := pageFaults{
		Minor: proc.PageFault.MinorFaults,
		Major: proc.PageFault.MajorFaults,
	}

	memStats := memPidStats{
		RSS:   proc.MemoryInfo.RSS,
		Data:  proc.MemoryInfo.Data,
		Stack: proc.MemoryInfo.Stack,
		Swap:  proc.MemoryInfo.Swap,
	}

	pidData := PidStats{
		Cpu:         utils.RoundFloat(proc.CPUPercent, "NONE", 2),
		Mem:         utils.RoundFloat(float64(proc.MemoryPercent), "NONE", 2),
		PidDetails:  pidDetails,
		CtxSwitches: ctxSwitches,
		PageFaults:  pgFaults,
		MemStats:    memStats,
	}
	return pidData, nil
}

func ExportPidJSON(pid int32, filename string, iter uint32, refreshRate uint64) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("Previous profile with name %s exists. Overwrite? (Y/N) ", filename)
		var choice string
		fmt.Scanf("%s", &choice)

		choice = strings.ToLower(choice)
		if choice != "y" {
			return nil
		} else {
			os.Remove(filename)
		}
	}

	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	encoder := json.NewEncoder(logFile)

	for i := uint32(0); i < iter; i++ {
		data, err := getPidDataJSON(pid)
		if err != nil {
			fmt.Println("Error in iteration", i, "Error:", err)
		} else {
			err = encoder.Encode(data)
			if err != nil {
				fmt.Println("Error in iteration", i, "Error:", err)
			}
		}

		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}

	return nil
}
