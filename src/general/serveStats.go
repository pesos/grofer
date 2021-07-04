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
	"math"
	"strings"
	"time"

	"github.com/pesos/grofer/src/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
}

// GetCPURates fetches and returns the current cpu rate
func GetCPURates() ([]float64, error) {
	cpuRates, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, err
	}
	return cpuRates, nil
}

// ServeCPURates serves the cpu rates to the cpu channel
func ServeCPURates(ctx context.Context, cpuChannel chan utils.DataStats) error {
	cpuRates, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	data := utils.DataStats{
		CpuStats: cpuRates,
		FieldSet: "CPU",
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case cpuChannel <- data:
		return nil
	}
}

// ServeMemRates serves stats about the memory to the data channel
func ServeMemRates(ctx context.Context, dataChannel chan utils.DataStats) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	memRates := []float64{roundOff(memory.Total), roundOff(memory.Used), roundOff(memory.Available), roundOff(memory.Free), roundOff(memory.Cached)}

	data := utils.DataStats{
		MemStats: memRates,
		FieldSet: "MEM",
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

// ServeTemperatureRates feeds temperature values from input sensors into the data channel
// Credits to https://github.com/cjbassi/gotop
func ServeTemperatureRates(ctx context.Context, dataChannel chan utils.DataStats) error {
	sensors, err := host.SensorsTemperatures()
	if err != nil && !strings.Contains(err.Error(), "Number of warnings:") {
		return err
	}
	// 2D string stores Header and Rows
	tempRates := [][]string{{"Sensor", "Temp(°C)"}}
	for _, sensor := range sensors {
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			temp_label := sensor.SensorKey
			// Only read input sensors
			label := strings.TrimSuffix(sensor.SensorKey, "_input")
			label = strings.TrimSuffix(label, "_thermal")
			if temp_label != label {
				temp := fmt.Sprintf("%.1f °C", sensor.Temperature)
				row := []string{label, temp}
				tempRates = append(tempRates, row)
			}
		}
	}

	data := utils.DataStats{
		TempStats: tempRates,
		FieldSet:  "TEMP",
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

// ServeDiskRates serves the disk rate data to the data channel
func ServeDiskRates(ctx context.Context, dataChannel chan utils.DataStats) error {
	var partitions []disk.PartitionStat
	var err error
	partitions, err = disk.Partitions(false)
	if err != nil {
		return err
	}

	rows := [][]string{{"Mount", "Total", "Used %", "Used", "Free", "FS Type"}}
	for _, value := range partitions {
		usageVals, _ := disk.Usage(value.Mountpoint)

		if strings.HasPrefix(value.Device, "/dev/loop") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
			continue
		} else {

			path := usageVals.Path
			total := fmt.Sprintf("%.2f G", float64(usageVals.Total)/(1024*1024*1024))
			used := fmt.Sprintf("%.2f G", float64(usageVals.Used)/(1024*1024*1024))
			usedPercent := fmt.Sprintf("%.2f %s", usageVals.UsedPercent, "%")
			free := fmt.Sprintf("%.2f G", float64(usageVals.Free)/(1024*1024*1024))
			fs := usageVals.Fstype
			row := []string{path, total, usedPercent, used, free, fs}
			rows = append(rows, row)

		}
	}

	data := utils.DataStats{
		DiskStats: rows,
		FieldSet:  "DISK",
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

// ServeNetRates serves info about the network to the data channel
func ServeNetRates(ctx context.Context, dataChannel chan utils.DataStats) error {
	netStats, err := net.IOCounters(false)
	if err != nil {
		return err
	}
	IO := make(map[string][]float64)
	for _, IOStat := range netStats {
		nic := []float64{float64(IOStat.BytesSent) / (1024), float64(IOStat.BytesRecv) / (1024)}
		IO[IOStat.Name] = nic
	}

	data := utils.DataStats{
		NetStats: IO,
		FieldSet: "NET",
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}
