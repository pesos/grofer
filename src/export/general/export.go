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

package general

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	cpuInfo "github.com/pesos/grofer/src/general"
	procInfo "github.com/pesos/grofer/src/process"
	"github.com/pesos/grofer/src/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type diskStats struct {
	Path     string  `json:"path"`
	Total    float64 `json:"total"`
	Used     float64 `json:"used"`
	UsedPerc float64 `json:"usedPerc"`
	Free     float64 `json:"free"`
	Fs       string  `json:"fs"`
}

type netStats struct {
	Sent float64 `json:"sent"`
	Recv float64 `json:"recv"`
}

type memStats struct {
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Used      float64 `json:"used"`
	Free      float64 `json:"free"`
}

// OverallStats describes the structure of each exported json object.
type OverallStats struct {
	Epoch     uint64              `json:"epoch"`
	CpuStats  []float64           `json:"cpu"`
	MemStats  memStats            `json:"mem"`
	DiskStats []diskStats         `json:"disk"`
	NetStats  map[string]netStats `json:"net"`
	CpuLoad   cpuInfo.CPULoad     `json:"cpuLoad"`
}

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
	Mem         float32         `json:"mem"`
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

// NewOverallStats returns a pointer to
// an empty OverallStats struct
func NewOverallStats() *OverallStats {
	return &OverallStats{}
}

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
}

func (data *OverallStats) updateData() error {
	startUpdateTime := uint64(time.Now().Unix())

	cpuRates, err := cpu.Percent(time.Second, true)
	if err == nil {
		data.CpuStats = cpuRates
	} else {
		return err
	}

	memory, err := mem.VirtualMemory()
	memRates := memStats{roundOff(memory.Total),
		roundOff(memory.Available),
		roundOff(memory.Used),
		roundOff(memory.Free),
	}

	if err == nil {
		data.MemStats = memRates
	} else {
		return err
	}

	partitions, err := disk.Partitions(false)
	if err == nil {
		var tempParts []diskStats
		for _, value := range partitions {
			usageVals, _ := disk.Usage(value.Mountpoint)

			if strings.HasPrefix(value.Device, "/dev/loop") {
				continue
			} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
				continue
			} else {
				path := usageVals.Path
				total := float64(usageVals.Total) / (1024 * 1024 * 1024)
				used := float64(usageVals.Used) / (1024 * 1024 * 1024)
				usedPercent := usageVals.UsedPercent
				free := float64(usageVals.Free) / (1024 * 1024 * 1024)
				fs := usageVals.Fstype
				temp := diskStats{path, total, used, usedPercent, free, fs}
				tempParts = append(tempParts, temp)
			}
		}
		data.DiskStats = tempParts
	} else {
		return err
	}

	netData, err := net.IOCounters(false)
	if err == nil {
		IO := make(map[string]netStats)
		for _, IOStat := range netData {
			nic := netStats{float64(IOStat.BytesSent) / (1024), float64(IOStat.BytesRecv) / (1024)}
			IO[IOStat.Name] = nic
		}
		data.NetStats = IO
	} else {
		return err
	}

	cpuLoad := cpuInfo.NewCPULoad()
	err = cpuLoad.UpdateCPULoad()
	if err != nil {
		return err
	}
	data.CpuLoad = *cpuLoad

	endUpdateTime := uint64(time.Now().Unix())
	avg := uint64((startUpdateTime + endUpdateTime) / 2)
	data.Epoch = avg
	return nil
}

func getJSONData(iter uint32, refreshRate uint64) ([]OverallStats, error) {
	var data []OverallStats
	var i uint32
	stats := NewOverallStats()
	for i = 0; i < iter; i++ {
		err := stats.updateData()
		if err != nil {
			return data, err
		}
		data = append(data, *stats)
		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}
	return data, nil
}

// ExportJSON exports data to a JSON file for a specified number of iterations
// and a specified refreshed rate.
func ExportJSON(fileName string, iter uint32, refreshRate uint64) error {
	toWrite, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer toWrite.Close()
	encoder := json.NewEncoder(toWrite)
	data, err := getJSONData(iter, refreshRate)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
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
		CreationTime:      utils.GetDateFromUnix(proc.CreateTime),
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
		Cpu:         proc.CPUPercent,
		Mem:         proc.MemoryPercent,
		PidDetails:  pidDetails,
		CtxSwitches: ctxSwitches,
		PageFaults:  pgFaults,
		MemStats:    memStats,
	}
	return pidData, nil
}

func ExportPidJSON(pid int32, filename string, iter uint32, refreshRate uint64) error {
	toWrite, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer toWrite.Close()
	encoder := json.NewEncoder(toWrite)
	var jsonData []PidStats
	for i := uint32(0); i < iter; i++ {
		data, err := getPidDataJSON(pid)
		if err != nil {
			break
		}
		jsonData = append(jsonData, data)
		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}

	err := encoder.Encode(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func ExportCSV(fileName string, iter uint32, refreshRate uint64) error {
	// TODO
	return nil
}
