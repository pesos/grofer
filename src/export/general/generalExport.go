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
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	cpuInfo "github.com/pesos/grofer/src/general"
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

// NewOverallStats returns a pointer to
// an empty OverallStats struct
func NewOverallStats() *OverallStats {
	return &OverallStats{}
}

func roundFloat(num float64, base string, precision int) float64 {
	x := num
	div := math.Pow10(precision)
	switch base {
	case "K":
		x /= 1024
	case "M":
		x /= (1024 * 1024)
	case "G":
		x /= (1024 * 1024 * 1024)
	}
	return math.Round(x*div) / div
}

func roundUint(num uint64, base string, precision int) float64 {
	x := float64(num)
	return roundFloat(x, base, precision)
}

func (data *OverallStats) updateData() error {
	startUpdateTime := uint64(time.Now().Unix())

	cpuRates, err := cpu.Percent(time.Second, true)
	if err == nil {
		for i, rate := range cpuRates {
			cpuRates[i] = roundFloat(rate, "NONE", 2)
		}
		data.CpuStats = cpuRates
	} else {
		return err
	}

	// Memory values in giga
	memory, err := mem.VirtualMemory()
	memRates := memStats{
		roundUint(memory.Total, "G", 2),
		roundUint(memory.Available, "G", 2),
		roundUint(memory.Used, "G", 2),
		roundUint(memory.Free, "G", 2),
	}

	if err == nil {
		data.MemStats = memRates
	} else {
		return err
	}

	// Disk values in giga
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
				total := roundUint(usageVals.Total, "G", 2)
				used := roundUint(usageVals.Used, "G", 2)
				usedPercent := roundFloat(usageVals.UsedPercent, "NONE", 2)
				free := roundUint(usageVals.Free, "G", 2)
				fs := usageVals.Fstype

				temp := diskStats{path, total, used, usedPercent, free, fs}
				tempParts = append(tempParts, temp)
			}
		}
		data.DiskStats = tempParts
	} else {
		return err
	}

	// Net values in kilo
	netData, err := net.IOCounters(false)
	if err == nil {
		IO := make(map[string]netStats)
		for _, IOStat := range netData {
			nic := netStats{
				roundUint(IOStat.BytesSent, "K", 2),
				roundUint(IOStat.BytesRecv, "K", 2),
			}
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

// ExportJSON exports data to a JSON file for a specified number of iterations
// and a specified refreshed rate.
func ExportJSON(filename string, iter uint32, refreshRate uint64) error {
	fmt.Println("Here", filename)
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("Previous profile with name 'grofer_profile' exists. Overwrite? (Y/N) ")
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
	stats := NewOverallStats()
	var i uint32

	for i = 0; i < iter; i++ {
		err := stats.updateData()
		if err != nil {
			return err
		}

		err = encoder.Encode(&stats)
		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}

	return nil
}
