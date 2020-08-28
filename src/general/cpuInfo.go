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
	"os"
	"os/exec"
	"time"

	gjson "github.com/tidwall/gjson"
)

type CPULoad struct {
	Usr    int64 `json:"usr"`
	Nice   int64 `json:"nice"`
	Sys    int64 `json:"sys"`
	Iowait int64 `json:"iowait"`
	Irq    int64 `json:"irq"`
	Soft   int64 `json:"soft"`
	Steal  int64 `json:"steal"`
	Guest  int64 `json:"guest"`
	Gnice  int64 `json:"gnice"`
	Idle   int64 `json:"idle"`
}

func NewCPULoad() *CPULoad {
	return &CPULoad{}
}

func (c *CPULoad) updateCPULoad() error {
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
	c.Usr = stats["usr"].Int()
	c.Nice = stats["nice"].Int()
	c.Sys = stats["sys"].Int()
	c.Iowait = stats["iowait"].Int()
	c.Irq = stats["irq"].Int()
	c.Soft = stats["soft"].Int()
	c.Steal = stats["steal"].Int()
	c.Guest = stats["guest"].Int()
	c.Gnice = stats["gnice"].Int()
	c.Idle = stats["idle"].Int()

	return nil
}

// GetCPULoad updated the CPULoad struct and serves the data to the data channel.
func GetCPULoad(cpuLoad *CPULoad,
	dataChannel chan *CPULoad,
	endChannel chan os.Signal,
	refreshRate int32) error {
	for {
		select {
		case <-endChannel: // Stop execution if end signal received
			return nil

		default: // Get Memory and CPU rates per core periodically
			err := cpuLoad.updateCPULoad()
			if err != nil {
				return err
			}
			dataChannel <- cpuLoad
			time.Sleep(time.Duration(refreshRate) * time.Millisecond)
		}
	}
}
