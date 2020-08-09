package general

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"time"

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
}

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
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
	memRates := memStats{roundOff(memory.Total), roundOff(memory.Available), roundOff(memory.Used), roundOff(memory.Free)}
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
	endUpdateTime := uint64(time.Now().Unix())
	avg := uint64((startUpdateTime + endUpdateTime) / 2)
	data.Epoch = avg
	return nil
}

func getJSONData(iter int64, interval uint64) ([]OverallStats, error) {
	var data []OverallStats
	var i int64
	stats := &OverallStats{}
	for i = 0; i < iter; i++ {
		err := stats.updateData()
		if err != nil {
			return data, err
		}
		data = append(data, *stats)
	}
	return data, nil
}

// ExportJSON exports data to a JSON file for a specified number of iterations
// and a specified refreshed rate.
func ExportJSON(fileName string, iter int64, interval uint64) error {
	toWrite, _ := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
	defer toWrite.Close()
	encoder := json.NewEncoder(toWrite)
	data, err := getJSONData(iter, interval)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}
