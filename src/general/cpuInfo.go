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
	"bufio"
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

// ReadCPULoad reads /proc/stat and returns the total load on all CPU cores
func (c *CPULoad) readCPULoad() error {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	// Read first line containing load values
	data, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		return err
	}
	// Split the first line into an array and omit the 1st index value as it only contains cpu/cpu<no>
	vals := strings.Fields(string(data))[1:]
	var avg [10]float64
	sum := 0
	// Convert and store the load values into a floating point array
	for i, x := range vals {
		curr, err := strconv.Atoi(x)
		if err != nil {
			return err
		} else {
			avg[i] = float64(curr)
			sum += curr
		}
	}
	// Calculate average values
	for i, x := range avg {
		avg[i] = 100 * x / float64(sum)
	}
	// Store values in CPULoad struct
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

	return err
}

// UpdateCPULoad updates fields of the type CPULoad
func (c *CPULoad) UpdateCPULoad() error {
	err := c.readCPULoad()
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

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- cpuLoad:
			return nil
		}
	})
}
