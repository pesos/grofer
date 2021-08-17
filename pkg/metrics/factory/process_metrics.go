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
	"github.com/pesos/grofer/pkg/metrics/process"
	processGraph "github.com/pesos/grofer/pkg/sink/tui/process"
	"github.com/pesos/grofer/pkg/utils"
	proc "github.com/shirou/gopsutil/process"
	"golang.org/x/sync/errgroup"
)

type processMetrics struct {
	refreshRate uint64
	sink        core.Sink // defaults to TUI.
	metricBus   chan []*proc.Process
}

// Serve serves metrics of all processes running in the system.
func (pm *processMetrics) Serve(opts ...FactoryOption) error {
	// apply command specific options.
	for _, opt := range opts {
		opt(pm)
	}
	eg, ctx := errgroup.WithContext(context.Background())

	// start producing metrics.
	eg.Go(func() error {
		alteredRefreshRate := uint64(4 * pm.refreshRate / 5)
		return utils.TickUntilDone(ctx, alteredRefreshRate, func() error {
			procs, err := proc.Processes()
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case pm.metricBus <- procs:
			}

			return err
		})
	})

	// Start consuming metrics.
	switch pm.sink {
	case core.TUI:
		eg.Go(func() error {
			return processGraph.AllProcVisuals(ctx, pm.metricBus, pm.refreshRate)
		})
	}

	return eg.Wait()
}

// SetSink sets the sink that consumes the produced metrics.
func (apm *processMetrics) SetSink(sink core.Sink) {
	apm.sink = sink
}

// ensure interface compliance.
var _ MetricScraper = (*processMetrics)(nil)

type singularProcessMetrics struct {
	pid         int32
	refreshRate uint64
	sink        core.Sink // defaults to TUI.
	metricBus   chan *process.Process
}

// Serve serves metrics of a particular process.
func (spm *singularProcessMetrics) Serve(opts ...FactoryOption) error {
	// apply command specific options.
	for _, opt := range opts {
		opt(spm)
	}
	eg, ctx := errgroup.WithContext(context.Background())

	process, err := process.NewProcess(spm.pid)
	if err != nil {
		return core.ErrInvalidPID
	}

	// start producing metrics.
	eg.Go(func() error {
		alteredRefreshRate := uint64(4 * spm.refreshRate / 5)
		return utils.TickUntilDone(ctx, alteredRefreshRate, func() error {
			process.UpdateProcInfo()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case spm.metricBus <- process:
			}

			return err
		})
	})

	// start consuming metrics.
	switch spm.sink {
	case core.TUI:
		eg.Go(func() error {
			return processGraph.ProcVisuals(ctx, spm.metricBus, spm.refreshRate)
		})
	}

	return eg.Wait()
}

// SetSink sets the sink that consumes the produced metrics.
func (spm *singularProcessMetrics) SetSink(sink core.Sink) {
	spm.sink = sink
}

// ensure interface compliance.
var _ MetricScraper = (*singularContainerMetrics)(nil)
