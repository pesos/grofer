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
	"log"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// GlobalStats gets stats about the mem and the CPUs and prints it.
func GlobalStats(endChannel chan os.Signal,
	memChannel chan []float64,
	cpuChannel chan []float64,
	diskChannel chan [][]string,
	netChannel chan map[string][]float64,
	wg *sync.WaitGroup) {

	for {
		select {
		case <-endChannel: // Stop execution if end signal received
			wg.Done()
			return

		default: // Get Memory and CPU rates per core periodically

			cpuUsageRates, err := cpu.Percent(time.Second, true)
			if err != nil {
				log.Fatal(err)
			}

			memoryStat, err := mem.VirtualMemory()
			if err != nil {
				log.Fatal(err)
			}

			partitions, err := disk.Partitions(false)
			if err != nil {
				log.Fatal(err)
			}

			netIO, err := net.IOCounters(false)
			if err != nil {
				log.Fatal(err)
			}

			go PrintCPURates(cpuUsageRates, cpuChannel)
			go PrintMemRates(memoryStat, memChannel)
			go PrintDiskRates(partitions, diskChannel)
			PrintNetRates(netIO, netChannel)
		}
	}
}
