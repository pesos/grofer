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
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Holds metrics for all exisitng conatiners
type ContainerMetrics struct {
	TotalCPU     float64
	TotalMem     float64
	TotalNet     netStat
	TotalBlk     blkStat
	PerContainer []PerContainerMetrics
}

// GetOverallMetrics provides metrics about all running containers in the form of ContainerMetrics structs
func GetOverallMetrics(ctx context.Context, cli *client.Client, cliMutex *sync.Mutex) (ContainerMetrics, error) {
	metrics := ContainerMetrics{}

	// Get list of containers
	cliMutex.Lock()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	cliMutex.Unlock()
	if err != nil {
		return metrics, err
	}

	metrcisChan := make(chan PerContainerMetrics, len(containers))

	// get per container metrics
	cliMutex.Lock()
	for _, container := range containers {
		go getMetrics(cli, ctx, container, metrcisChan)
	}

	var totalCPU, totalMem float64
	totalNet := netStat{}
	totalBlk := blkStat{}

	// Aggregate metrics and compute total metrics
	for range containers {
		metric := <-metrcisChan

		totalCPU += metric.Cpu

		totalMem += metric.Mem

		totalNet.Rx += metric.Net.Rx
		totalNet.Tx += metric.Net.Tx

		totalBlk.Read += metric.Blk.Read
		totalBlk.Write += metric.Blk.Write

		metrics.PerContainer = append(metrics.PerContainer, metric)
	}
	cliMutex.Unlock()

	metrics.TotalCPU = totalCPU
	metrics.TotalMem = totalMem
	metrics.TotalNet = totalNet
	metrics.TotalBlk = totalBlk

	return metrics, nil
}

func getMetrics(cli *client.Client, ctx context.Context, c types.Container, ch chan PerContainerMetrics) {

	// Send back metrics
	metrics := PerContainerMetrics{}
	defer func() {
		ch <- metrics
	}()

	stats, err := cli.ContainerStats(ctx, c.ID, false)
	if err != nil {
		return
	}

	data := types.StatsJSON{}
	err = json.NewDecoder(stats.Body).Decode(&data)
	if err != nil {
		return
	}
	stats.Body.Close()

	// Calculate CPU percent
	cpuPercent := getCPUPercent(&data)

	// Calculate Memory
	memPercent := float64(data.MemoryStats.Usage) / float64(data.MemoryStats.Limit) * 100

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

	// Calculate Network
	var rx, tx float64

	for _, v := range data.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	metrics = PerContainerMetrics{
		ID:     c.ID[:10],
		Image:  c.Image,
		Name:   strings.TrimLeft(strings.Join(c.Names, ","), "/"),
		Status: c.Status,
		State:  c.State,
		Cpu:    cpuPercent,
		Mem:    memPercent,
		Net:    netStat{Rx: rx, Tx: tx},
		Blk:    blkStat{Read: blkRead, Write: blkWrite},
	}
}
