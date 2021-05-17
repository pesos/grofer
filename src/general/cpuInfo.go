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
func readCPULoad(c *CPULoad) error {
	data, error := os.ReadFile("/proc/stat")
	if error != nil {
		return error
	}
	stringData := string(data)
	lines := strings.Split(stringData, "\n")
	// First line should store in the format:
	// "cpu <usr> <nice> <system> <idle> <iowait> <irq> <softirq> <steal> <guest> <guest_nice>"
	vals := strings.Fields(lines[0])
	var avg [10]float64
	sum := 0
	for i, x := range vals {
		// Start from index 1 as first element is the cpu/cpu<no>
		if i < 1 {
			continue
		}
		curr, err := strconv.Atoi(x)
		if err != nil {
			return error
		} else {
			avg[i-1] = float64(curr)
			sum += curr
		}
	}
	for i, x := range avg {
		avg[i] = 100 * x / float64(sum)
	}

	c.Usr = int(avg[0])
	c.Nice = int(avg[1])
	c.Sys = int(avg[2])
	c.Idle = int(avg[3])
	c.Iowait = int(avg[4])
	c.Irq = int(avg[5])
	c.Soft = int(avg[6])
	c.Steal = int(avg[7])
	c.Guest = int(avg[8])
	c.Gnice = int(avg[9])

	return error
}

// UpdateCPULoad updates fields of the type CPULoad
func (c *CPULoad) UpdateCPULoad() error {
	err := readCPULoad(c)
	if err != nil {
		return err
	}
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
