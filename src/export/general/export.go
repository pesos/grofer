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

type terminate struct {
	finished bool
	err      error
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

func checkFileExistence(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
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

func getJSONData(iter uint32, refreshRate uint64, exportChan chan OverallStats, done chan terminate) {
	var i uint32
	stats := NewOverallStats()
	exit := terminate{
		finished: false,
		err:      nil,
	}
	for i = 0; i < iter; i++ {
		err := stats.updateData()
		if err != nil {
			exit.err = err
			exit.finished = true
			done <- exit
			break
		}
		exportChan <- *stats
		done <- exit
		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}
	exit.finished = true
	done <- exit
}

// ExportJSON exports data to a JSON file for a specified number of iterations
// and a specified refreshed rate.
func ExportJSON(fileName string, iter uint32, refreshRate uint64, bufferFraction float32) error {
	if exists := checkFileExistence(fileName); exists {
		err := os.Remove(fileName)
		if err != nil {
			return fmt.Errorf("file %s already exists and cannot be overwritten\n", fileName)
		}
		fmt.Printf("WARNING: file %s already exists and will be overwritten\n", fileName)
	}
	toWrite, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer toWrite.Close()

	// add an opening [ to signify opening of array of json objects
	_, err = toWrite.Write([]byte("[\n"))
	if err != nil {
		return err
	}

	bufferSize := int(float32(iter) * bufferFraction)
	if bufferSize == 0 {
		bufferSize = 1
	}
	exportChannel := make(chan OverallStats, bufferSize)
	doneChannel := make(chan terminate)
	defer close(exportChannel)
	defer close(doneChannel)

	encoder := json.NewEncoder(toWrite)
	encoder.SetIndent("", "  ")
	go getJSONData(iter, refreshRate, exportChannel, doneChannel)

	// encode and write the first json object.
	initialObject := <-exportChannel
	err = encoder.Encode(initialObject)
	if err != nil {
		return err
	}

	// buffer a user-defined fraction of total data in memory before making an IO
	// buffering is done instead of writing on a per record basis as the number of
	// hardware I/Os might pose a significant overhead for large values of iter.
	for {
		select {
		case status := <-doneChannel:
			if status.finished {
				// add a ] to close the array of JSON objects.
				_, err = toWrite.Write([]byte("]"))
				if err != nil {
					return err
				}
				return status.err
			}
			if len(exportChannel) == bufferSize {
				for len(exportChannel) > 0 {
					exportObj := <-exportChannel
					_, err := toWrite.Write([]byte(","))
					if err != nil {
						return err
					}
					err = encoder.Encode(exportObj)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
