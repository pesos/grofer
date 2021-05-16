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
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pesos/grofer/src/utils"
)

// CPULoad type contains info about load on CPU from various sources
// as well as general stats about the CPU.
type CPULoad struct {
	CPURates [][]string `json:"-"`
	Usr      int        `json:"usr"`
	Nice     int        `json:"nice"`
	Sys      int        `json:"sys"`
	Iowait   int        `json:"iowait"`
	Soft     int        `json:"soft"`
	Steal    int        `json:"steal"`
	Guest    int        `json:"guest"`
	Gnice    int        `json:"gnice"`
	Idle     int        `json:"idle"`
	Irq      int        `json:"irq"`
}

// NewCPULoad is a constructor for the CPULoad type.
func NewCPULoad() *CPULoad {
	return &CPULoad{}
}

// ReadCPULoad reads /proc/stat and returns the average load
func ReadCPULoad() ([10]float64, error) {
	data, error := os.ReadFile("/proc/stat")
	if error != nil {
		return [10]float64{}, error
	}
	string_data := string(data)
	lines := strings.Split(string_data, "\n")
	vals := strings.Split(lines[0], " ")
	var avg [10]float64
	sum := 0
	for i, x := range vals {
		if i < 2 {
			continue
		}
		curr, err := strconv.Atoi(x)
		if err != nil {
			return [10]float64{}, error
		} else {
			sum += curr
		}
	}

	for i, x := range vals {
		if i < 2 {
			continue
		}
		curr, err := strconv.Atoi(x)
		if err != nil {
			fmt.Print(err)
			return [10]float64{}, err
		} else {
			avg[i-2] = 100 * float64(curr) / float64(sum)
		}
	}
	return avg, error
}

// UpdateCPULoad updates fields of the type CPULoad
func (c *CPULoad) UpdateCPULoad() error {
	stats, err := ReadCPULoad()
	if err != nil {
		return err
	}
	c.Usr = int(stats[0])
	c.Nice = int(stats[1])
	c.Sys = int(stats[2])
	c.Idle = int(stats[3])
	c.Iowait = int(stats[4])
	c.Irq = int(stats[5])
	c.Soft = int(stats[6])
	c.Steal = int(stats[7])
	c.Guest = int(stats[8])
	c.Gnice = int(stats[9])

	cpuRates, err := GetCPURates()
	if err != nil {
		return err
	}

	rate := []string{}
	cpus := []string{}
	for i, cpuRate := range cpuRates {
		cpus = append(cpus, "CPU "+strconv.Itoa(i))
		rate = append(rate, fmt.Sprintf("%.2f%%", cpuRate))
	}
	rates := [][]string{cpus, rate}

	c.CPURates = rates

	return nil
}

// GetCPULoad updated the CPULoad struct and serves the data to the data channel.
func GetCPULoad(ctx context.Context, cpuLoad *CPULoad, dataChannel chan *CPULoad, refreshRate uint64) error {
	return utils.TickUntilDone(ctx, int64(refreshRate), func() error {
		err := cpuLoad.UpdateCPULoad()
		if err != nil {
			return err
		}
		dataChannel <- cpuLoad

		return nil
	})
}
