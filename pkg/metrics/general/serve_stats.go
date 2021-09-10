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
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/utils"
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

// ServeInfo provides information about the system such as OS info, uptime, boot time, etc.
func ServeInfo(ctx context.Context, cpuChannel chan AggregatedMetrics) error {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return err
	}

	hostInfo := [][]string{
		{"Hostname", info.Hostname},
		{"Up Time", utils.SecondsToHuman(int(info.Uptime))},
		{"Boot Time", utils.GetDateFromUnix(int64(info.BootTime * 1000))},
		{"Processes", fmt.Sprintf("%d", info.Procs)},
		{"OS/Platform", fmt.Sprintf("%s/%s %s", info.OS, info.Platform, info.PlatformVersion)},
		{"Kernel/Arch", fmt.Sprintf("%s/%s", info.KernelVersion, info.KernelArch)},
	}

	data := AggregatedMetrics{
		FieldSet: "INFO",
		HostInfo: hostInfo,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case cpuChannel <- data:
		return nil
	}
}

// ServeBattery serves battery percentage information
func ServeBattery(ctx context.Context, cpuChannel chan AggregatedMetrics) error {

	_, err1 := os.Stat("/sys/class/power_supply/BAT0/charge_now")
	_, err2 := os.Stat("/sys/class/power_supply/BAT0/charge_full")

	data := AggregatedMetrics{
		FieldSet:       "BATTERY",
		BatteryPercent: 0,
	}

	if err1 == nil && err2 == nil {
		currentBS, _ := ioutil.ReadFile("/sys/class/power_supply/BAT0/charge_now")
		fullBS, _ := ioutil.ReadFile("/sys/class/power_supply/BAT0/charge_full")

		current, err1 := strconv.ParseFloat(strings.Trim(string(currentBS), "\t\n "), 64)
		if err1 != nil {
			return err1
		}

		full, err2 := strconv.ParseFloat(strings.Trim(string(fullBS), "\t\n "), 64)
		if err2 != nil {
			return err2
		}

		if full == 0 {
			full = 1
		}

		data.BatteryPercent = int((current / full) * 100)

	} else if os.IsNotExist(err1) || os.IsNotExist(err2) {
		return core.ErrBatteryNotFound
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case cpuChannel <- data:
		return nil
	}

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
func ServeCPURates(ctx context.Context, cpuChannel chan AggregatedMetrics) error {
	cpuRates, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	data := AggregatedMetrics{
		CPUStats: cpuRates,
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
func ServeMemRates(ctx context.Context, dataChannel chan AggregatedMetrics) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	memRates := []float64{roundOff(memory.Total), roundOff(memory.Used), roundOff(memory.Available), roundOff(memory.Free), roundOff(memory.Cached)}

	data := AggregatedMetrics{
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
func ServeTemperatureRates(ctx context.Context, dataChannel chan AggregatedMetrics) error {
	sensors, err := host.SensorsTemperatures()
	if err != nil && !strings.Contains(err.Error(), "Number of warnings:") {
		return err
	}
	// 2D string stores Header and Rows
	tempRates := [][]string{{"Sensor", "Temp(°C)"}}
	for _, sensor := range sensors {
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			tempLabel := sensor.SensorKey
			// Only read input sensors
			label := strings.TrimSuffix(sensor.SensorKey, "_input")
			label = strings.TrimSuffix(label, "_thermal")
			if tempLabel != label {
				temp := fmt.Sprintf("%.1f °C", sensor.Temperature)
				row := []string{label, temp}
				tempRates = append(tempRates, row)
			}
		}
	}

	data := AggregatedMetrics{
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
func ServeDiskRates(ctx context.Context, dataChannel chan AggregatedMetrics) error {
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

	data := AggregatedMetrics{
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
func ServeNetRates(ctx context.Context, dataChannel chan AggregatedMetrics) error {
	netStats, err := net.IOCounters(false)
	if err != nil {
		return err
	}
	IO := make(map[string][]float64)
	for _, IOStat := range netStats {
		nic := []float64{float64(IOStat.BytesSent) / (1024), float64(IOStat.BytesRecv) / (1024)}
		IO[IOStat.Name] = nic
	}

	data := AggregatedMetrics{
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
