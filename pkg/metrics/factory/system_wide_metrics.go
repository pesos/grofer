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

package factory

import (
	"context"

	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/general"
	overallGraph "github.com/pesos/grofer/pkg/sink/tui/general"
	"golang.org/x/sync/errgroup"
)

type systemWideMetrics struct {
	cpuInfo     bool
	batteryInfo bool
	sink        core.Sink // defaults to TUI.
	refreshRate uint64
}

// Serve serves system wide metrics.
func (swm *systemWideMetrics) Serve(opts ...Option) error {
	// apply command specific options.
	for _, opt := range opts {
		opt(swm)
	}

	if swm.cpuInfo {
		return swm.serveCPUInfo()
	}

	if swm.batteryInfo {
		return swm.serveBatteryInfo()
	}
	return swm.serveGenericMetrics()
}

// serveGenericMetrics serves generic metrics such as metrics related to
// network, memory, CPU etc.
func (swm *systemWideMetrics) serveGenericMetrics() error {
	eg, ctx := errgroup.WithContext(context.Background())
	metricBus := make(chan general.AggregatedMetrics, 1)

	// start producing metrics.
	eg.Go(func() error {
		alteredRefreshRate := uint64(4 * swm.refreshRate / 5)
		return general.GlobalStats(ctx, metricBus, alteredRefreshRate)
	})

	// start consuming metrics.
	switch swm.sink {
	case core.TUI:
		eg.Go(func() error {
			return overallGraph.RenderCharts(ctx, metricBus, swm.refreshRate)
		})
	}

	return eg.Wait()
}

// serveCPUInfo serves specific CPU metrics such as time spent servicing
// different type of IRQs.
func (swm *systemWideMetrics) serveCPUInfo() error {
	eg, ctx := errgroup.WithContext(context.Background())
	metricBus := make(chan *general.CPULoad, 1)

	// start producing metrics.
	eg.Go(func() error {
		cpuLoad := general.NewCPULoad()
		alteredRefreshRate := uint64(4 * swm.refreshRate / 5)
		return general.GetCPULoad(ctx, cpuLoad, metricBus, alteredRefreshRate)
	})

	// start consuming metrics.
	switch swm.sink {
	case core.TUI:
		eg.Go(func() error {
			return overallGraph.RenderCPUinfo(ctx, metricBus, swm.refreshRate)
		})
	}

	return eg.Wait()
}

// serveBatteryInfo serves battery stats such as the name of the manufacturer
// and other battery specific information.
func (swm *systemWideMetrics) serveBatteryInfo() error {
	eg, ctx := errgroup.WithContext(context.Background())
	metricBus := make(chan general.BatteryData, 1)

	// start producing metrics.
	eg.Go(func() error {
		batteryInfo := general.NewBatteryInfo()
		alteredRefreshRate := uint64(4 * swm.refreshRate / 5)
		return general.GetBatteryInfo(ctx, batteryInfo, metricBus, alteredRefreshRate)
	})

	// start consuming metrics.
	switch swm.sink {
	case core.TUI:
		eg.Go(func() error {
			return overallGraph.RenderBatteryinfo(ctx, metricBus, swm.refreshRate)
		})
	}

	return eg.Wait()
}

// SetSink sets the Sink for the produced metrics.
func (swm *systemWideMetrics) SetSink(sink core.Sink) {
	swm.sink = sink
}

// ensure interface compliance.
var _ MetricScraper = (*systemWideMetrics)(nil)
