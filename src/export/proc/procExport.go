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
	ChildProcs  map[int]string  `json:"childProcs"`
	PidDetails  ProcDetails     `json:"pid"`
	MemStats    memPidStats     `json:"memStats"`
	CtxSwitches contextSwitches `json:"ctxSwitches"`
	PageFaults  pageFaults      `json:"pageFaults"`
	Mem         float64         `json:"mem"`
	Cpu         float64         `json:"cpu"`
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

// getPidDataJSON returns a PidStats structure populated with information about the process specified by pid
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

// ExportJSON exports data particular to a given process (given by pid) to a
// JSON file for a specified number of iterations and a specified
// refreshed rate.
func ExportPidJSON(pid int32, filename string, iter uint32, refreshRate uint64) error {
	// Verify if previous profile exists and whether or not to overwrite
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

	// Open file pointer to file to be written
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Encoder to encode JSON data into file
	encoder := json.NewEncoder(logFile)

	// Encode JSON object by object into file
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
