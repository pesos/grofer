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
	"strings"
	"time"

	cpuInfo "github.com/pesos/grofer/pkg/metrics/general"
	"github.com/pesos/grofer/pkg/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type diskStats struct {
	Path     string  `json:"path"`
	Fs       string  `json:"fs"`
	Total    float64 `json:"total"`
	Used     float64 `json:"used"`
	UsedPerc float64 `json:"usedPerc"`
	Free     float64 `json:"free"`
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
	NetStats  map[string]netStats `json:"net"`
	CPUStats  []float64           `json:"cpu"`
	DiskStats []diskStats         `json:"disk"`
	CPULoad   cpuInfo.CPULoad     `json:"cpuLoad"`
	MemStats  memStats            `json:"mem"`
	Epoch     uint64              `json:"epoch"`
}

// NewOverallStats returns a pointer to an empty OverallStats struct
func NewOverallStats() *OverallStats {
	return &OverallStats{}
}

// updateData updates values of a received OverallStats struct, returns error on failure of updates
func (data *OverallStats) updateData() error {
	startUpdateTime := uint64(time.Now().Unix())

	cpuRates, err := cpu.Percent(time.Second, true)
	if err == nil {
		for i, rate := range cpuRates {
			cpuRates[i] = utils.RoundFloat(rate, "NONE", 2)
		}
		data.CPUStats = cpuRates
	} else {
		return err
	}

	// Memory values in giga
	memory, err := mem.VirtualMemory()
	memRates := memStats{
		utils.RoundUint(memory.Total, "G", 2),
		utils.RoundUint(memory.Available, "G", 2),
		utils.RoundUint(memory.Used, "G", 2),
		utils.RoundUint(memory.Free, "G", 2),
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
				total := utils.RoundUint(usageVals.Total, "G", 2)
				used := utils.RoundUint(usageVals.Used, "G", 2)
				usedPercent := utils.RoundFloat(usageVals.UsedPercent, "NONE", 2)
				free := utils.RoundUint(usageVals.Free, "G", 2)
				fs := usageVals.Fstype
				temp := diskStats{path, fs, total, used, usedPercent, free}
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
				utils.RoundUint(IOStat.BytesSent, "K", 2),
				utils.RoundUint(IOStat.BytesRecv, "K", 2),
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
	data.CPULoad = *cpuLoad

	endUpdateTime := uint64(time.Now().Unix())
	avg := uint64((startUpdateTime + endUpdateTime) / 2)
	data.Epoch = avg

	return nil
}

// ToJSON exports data to a JSON file for a specified number of iterations
// and a specified refreshed rate.
func ToJSON(filename string, iter uint32, refreshRate uint64) error {
	// Verify if previous profile exists and whether or not to overwrite
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("Previous metric file with name %s exists. Overwrite? (Y/N) ", filename)
		var choice string
		fmt.Scanf("%s", &choice)

		choice = strings.ToLower(choice)
		if choice != "y" {
			return nil
		}
		os.Remove(filename)
	}

	// Open file pointer to file to be written
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Encoder to encode JSON data into file
	encoder := json.NewEncoder(logFile)
	stats := NewOverallStats()

	// Encode JSON object by object into file
	for i := uint32(0); i < iter; i++ {
		err := stats.updateData()
		if err != nil {
			fmt.Println("Error in iteration", i, "Error:", err)
		} else {
			err = encoder.Encode(&stats)
			if err != nil {
				fmt.Println("Error in iteration", i, "Error:", err)
			}
		}

		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}

	return nil
}
