package container

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerMetrics struct {
	totalCPU     float64
	totalMem     float64
	totalNet     netStat
	totalBlk     blkStat
	perContainer []perContainerMetrics
}

type perContainerMetrics struct {
	containerID string
	image       string
	name        string
	status      string
	state       string
	cpu         float64
	mem         float64
	net         netStat
	blk         blkStat
}

type netStat struct {
	rx float64
	tx float64
}

type blkStat struct {
	read  uint64
	write uint64
}

func GetOverallMetrics() {
	metrics := ContainerMetrics{}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Get list of containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	metrcisChan := make(chan perContainerMetrics, len(containers))

	// get per container metrics
	for _, container := range containers {
		go getMetrics(cli, ctx, container, metrcisChan)
	}

	var totalCPU, totalMem float64
	totalNet := netStat{}
	totalBlk := blkStat{}

	// Aggregate metrics and compute total metrics
	for range containers {
		metric := <-metrcisChan

		totalCPU += metric.cpu

		totalMem += metric.mem

		totalNet.rx += metric.net.rx
		totalNet.tx += metric.net.tx

		totalBlk.read += metric.blk.read
		totalBlk.write += metric.blk.write

		metrics.perContainer = append(metrics.perContainer, metric)
	}

	metrics.totalCPU = totalCPU
	metrics.totalMem = totalMem
	metrics.totalNet = totalNet
	metrics.totalBlk = totalBlk

	fmt.Println(metrics)
}

func getMetrics(cli *client.Client, ctx context.Context, c types.Container, ch chan perContainerMetrics) {

	stats, _ := cli.ContainerStatsOneShot(ctx, c.ID)
	data := types.StatsJSON{}
	err := json.NewDecoder(stats.Body).Decode(&data)
	if err != nil {
		ch <- perContainerMetrics{}
		return
	}

	// Calculate CPU percent
	cpuPercent := 0.0

	cpuDelta := float64(data.CPUStats.CPUUsage.TotalUsage) - float64(data.PreCPUStats.CPUUsage.TotalUsage)

	systemDelta := float64(data.CPUStats.SystemUsage) - float64(data.PreCPUStats.SystemUsage)

	if cpuDelta > 0.0 && systemDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(data.CPUStats.CPUUsage.PercpuUsage)) * 100.0
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

	// Calculate Network
	var rx, tx float64

	for _, v := range data.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	// Calculate Memory
	memPercent := float64(data.MemoryStats.Usage) / float64(data.MemoryStats.Limit) * 100

	metrics := perContainerMetrics{
		containerID: c.ID[:10],
		image:       c.Image,
		name:        strings.Join(c.Names, ","),
		status:      c.Status,
		state:       c.State,
		cpu:         cpuPercent,
		mem:         memPercent,
		net:         netStat{rx: rx, tx: tx},
		blk:         blkStat{read: blkRead, write: blkWrite},
	}

	// Send back metrics
	ch <- metrics
}
