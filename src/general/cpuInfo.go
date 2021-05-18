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
	"os/exec"
	"strconv"

	"github.com/pesos/grofer/src/utils"
	gjson "github.com/tidwall/gjson"
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

// UpdateCPULoad updates fields of the type CPULoad
func (c *CPULoad) UpdateCPULoad() error {
	mpstat := "mpstat"
	arg0 := "-o"
	arg1 := "JSON"
	cmd := exec.Command(mpstat, arg0, arg1)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	statsExtract := gjson.Get(string(stdout), "sysstat.hosts.0.statistics.0.cpu-load.0")
	stats := statsExtract.Map()
	c.Usr = int(stats["usr"].Int())
	c.Nice = int(stats["nice"].Int())
	c.Sys = int(stats["sys"].Int())
	c.Iowait = int(stats["iowait"].Int())
	c.Irq = int(stats["irq"].Int())
	c.Soft = int(stats["soft"].Int())
	c.Steal = int(stats["steal"].Int())
	c.Guest = int(stats["guest"].Int())
	c.Gnice = int(stats["gnice"].Int())
	c.Idle = int(stats["idle"].Int())

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
