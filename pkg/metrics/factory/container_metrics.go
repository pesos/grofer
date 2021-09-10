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

	containerGraph "github.com/pesos/grofer/pkg/sink/tui/container"

	"github.com/docker/docker/client"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/container"
	"github.com/pesos/grofer/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type containerMetrics struct {
	client      *client.Client
	all         bool
	refreshRate uint64
	sink        core.Sink // defaults to TUI.
	metricBus   chan container.OverallMetrics
}

// Serve serves metrics for all containers running on the system.
func (cms *containerMetrics) Serve(opts ...Option) error {
	// apply command specific options.
	for _, opt := range opts {
		opt(cms)
	}
	eg, ctx := errgroup.WithContext(context.Background())

	// start producing metrics.
	eg.Go(func() error {
		return utils.TickUntilDone(ctx, cms.refreshRate, func() error {
			metrics, err := container.GetOverallMetrics(ctx, cms.client, cms.all)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case cms.metricBus <- metrics:
			}

			return nil
		})
	})

	// Start consuming metrics.
	switch cms.sink {
	case core.TUI:
		eg.Go(func() error {
			return containerGraph.OverallVisuals(ctx, cms.client, cms.all, cms.metricBus, cms.refreshRate)
		})
	}

	return eg.Wait()
}

// SetSink sets the Sink for the produced metrics.
func (cms *containerMetrics) SetSink(sink core.Sink) {
	cms.sink = sink
}

// ensure interface compliance.
var _ MetricScraper = (*containerMetrics)(nil)

type singularContainerMetrics struct {
	client      *client.Client
	refreshRate uint64
	cid         string
	sink        core.Sink // defaults to TUI.
	metricBus   chan container.PerContainerMetrics
}

// Serve serves metrics for a particular container running on the system.
func (scms *singularContainerMetrics) Serve(opts ...Option) error {
	// apply command specific options.
	for _, opt := range opts {
		opt(scms)
	}
	eg, ctx := errgroup.WithContext(context.Background())

	// start producing metrics.
	eg.Go(func() error {
		return utils.TickUntilDone(ctx, scms.refreshRate, func() error {
			metrics, err := container.GetContainerMetrics(ctx, scms.client, scms.cid)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case scms.metricBus <- metrics:
			}

			return nil
		})
	})

	// Start consuming metrics.
	switch scms.sink {
	case core.TUI:
		eg.Go(func() error {
			return containerGraph.PerContainerVisuals(ctx, scms.metricBus, scms.refreshRate)
		})
	}

	return eg.Wait()
}

// SetSink sets the Sink for the produced metrics.
func (scms *singularContainerMetrics) SetSink(sink core.Sink) {
	scms.sink = sink
}

// ensure interface compliance.
var _ MetricScraper = (*singularContainerMetrics)(nil)
