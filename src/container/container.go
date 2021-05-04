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

package container

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	ui "github.com/gizak/termui/v3"
	"github.com/pesos/grofer/src/general"
)

type PerContainerMetrics struct {
	// Metrics common for Overall and per container
	ID     string
	Image  string
	Name   string
	Status string
	State  string
	Cpu    float64
	Mem    float64
	Net    netStat
	Blk    blkStat
	// Metrics specific to per container
	Pid     string
	NetInfo []netInfo
	PerCPU  []string
	PortMap []portMap
	Mounts  []mountInfo
	Procs   []procInfo
}

type netStat struct {
	Rx float64
	Tx float64
}

type blkStat struct {
	Read  uint64
	Write uint64
}

type netInfo struct {
	Name    string
	Driver  string
	Ip      string
	Ingress bool
}

type mountInfo struct {
	Src  string
	Dst  string
	Mode string
}

type portMap struct {
	Host      int
	Container int
	Protocol  string
}

type procInfo struct {
	UID string
	PID string
	CMD string
}

// GetContainerMetrics provides per container metrics in the form of PerContainerMetrics Structs
func GetContainerMetrics(cid string) (PerContainerMetrics, error) {

	metrics := PerContainerMetrics{}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return metrics, err
	}

	// Get container using a filter
	args := filters.NewArgs(
		filters.KeyValuePair{
			Key:   "id",
			Value: cid,
		},
	)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: args})
	if err != nil {
		return metrics, err
	}

	if len(containers) > 1 {
		return metrics, fmt.Errorf("multiple containers with same ID exist")
	} else if len(containers) < 1 {
		return metrics, general.ErrInvalidContainer
	}

	c := containers[0]

	// Get PID
	inspectData, err := cli.ContainerInspect(ctx, cid)
	if err != nil {
		return metrics, nil
	}

	// Get Container Stats
	stats, _ := cli.ContainerStats(ctx, cid, false)
	data := types.StatsJSON{}
	err = json.NewDecoder(stats.Body).Decode(&data)
	if err != nil {
		return metrics, err
	}

	// Calculate CPU  and per CPU percent
	cpuPercent := 0.0
	numCPUs := len(data.CPUStats.CPUUsage.PercpuUsage)

	cpuDelta := float64(data.CPUStats.CPUUsage.TotalUsage) - float64(data.PreCPUStats.CPUUsage.TotalUsage)

	systemDelta := float64(data.CPUStats.SystemUsage) - float64(data.PreCPUStats.SystemUsage)

	if cpuDelta > 0.0 && systemDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(numCPUs) * 100.0
	}

	// Get Per CPU utilizations
	preLen := len(data.PreCPUStats.CPUUsage.PercpuUsage)
	postLen := len(data.CPUStats.CPUUsage.PercpuUsage)
	numCPUs = ui.MaxInt(preLen, postLen)

	perCpuPercents := make([]string, numCPUs)

	// If first run, skip percpu metrics
	if preLen != postLen {
		for i := range perCpuPercents {
			perCpuPercents[i] = "NA"
		}
	} else {
		for i, usage := range data.CPUStats.CPUUsage.PercpuUsage {
			perCpuPercent := 0.0

			cpuDelta := float64(usage) - float64(data.PreCPUStats.CPUUsage.PercpuUsage[i])

			if cpuDelta > 0.0 && systemDelta > 0.0 {
				perCpuPercent = (cpuDelta / systemDelta) * float64(numCPUs) * 100.0
			}
			perCpuPercents[i] = fmt.Sprintf("%.2f%%", perCpuPercent)
		}
	}

	// Calculate Memory usage
	memPercent := float64(data.MemoryStats.Usage) / float64(data.MemoryStats.Limit) * 100

	// Calculate Network Metrics
	var rx, tx float64
	for _, v := range data.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	// Calculate blk IO
	var blkRead, blkWrite uint64
	for _, bioEntry := range data.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}

	// Get Network Settings
	netData := []netInfo{}
	for _, network := range c.NetworkSettings.Networks {
		id := network.NetworkID

		net, err := cli.NetworkInspect(ctx, id, types.NetworkInspectOptions{})
		if err != nil {
			continue
		}

		n := netInfo{
			Name:    net.Name,
			Driver:  net.Driver,
			Ip:      network.IPAddress,
			Ingress: net.Ingress,
		}

		netData = append(netData, n)
	}

	// Get Port mappings
	portData := []portMap{}
	for _, port := range c.Ports {
		p := portMap{
			Host:      int(port.PublicPort),
			Container: int(port.PrivatePort),
			Protocol:  port.Type,
		}

		portData = append(portData, p)
	}

	// Get mounted volumes
	mountData := []mountInfo{}
	for _, mount := range c.Mounts {
		m := mountInfo{
			Src:  mount.Source,
			Dst:  mount.Destination,
			Mode: mount.Mode,
		}

		mountData = append(mountData, m)
	}

	// Get processes in container
	procs, err := cli.ContainerTop(ctx, cid, []string{})
	if err != nil {
		return metrics, nil
	}

	procData := []procInfo{}
	for _, proc := range procs.Processes {
		p := procInfo{
			UID: proc[0],
			PID: proc[1],
			CMD: proc[7],
		}

		procData = append(procData, p)
	}

	// Populate metrics
	metrics = PerContainerMetrics{
		ID:      c.ID[:10],
		Image:   c.Image,
		Name:    strings.TrimLeft(strings.Join(c.Names, ","), "/"),
		Status:  c.Status,
		State:   c.State,
		Cpu:     cpuPercent,
		Mem:     memPercent,
		Net:     netStat{Rx: rx, Tx: tx},
		Blk:     blkStat{Read: blkRead, Write: blkWrite},
		Pid:     fmt.Sprintf("%d", inspectData.State.Pid),
		NetInfo: netData,
		PerCPU:  perCpuPercents,
		PortMap: portData,
		Mounts:  mountData,
		Procs:   procData,
	}

	return metrics, nil
}
