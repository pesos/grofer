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
	"errors"
	"strconv"

	proc "github.com/shirou/gopsutil/process"

	"github.com/docker/docker/client"
	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/container"
	"github.com/pesos/grofer/pkg/metrics/process"
)

// MetricScraperFactory constructs a MetricScaper for a command
// and returns it.
type MetricScraperFactory struct {
	// command is the command for which a MetricScraper
	// is created. This defaults to core.MainCommand.
	command core.Command
	// singularEntityMetrics indicate whether metrics
	// that need to be scraped are for a singular entity
	// or not, for ex - metrics for some process ID or
	// some container ID.
	singularEntityMetrics bool
	// entity is an identifier that can be used to scrape
	// metrics for it, ex CID, PID. Conversion from string
	// to the appropriate form of this entity should be
	// handled by an implementation of the MetricScraper
	// interface.
	entity string
	// scrapeIntervalMillisecond is the frequency in ms at
	// which metrics will be scraped.
	scrapeIntervalMillisecond uint64
}

// NewMetricScraperFactory is a constructor for the MetricScraperFactory type.
// By default, this will be for the core.MainCommand command.
func NewMetricScraperFactory() *MetricScraperFactory {
	return &MetricScraperFactory{}
}

// ForCommand sets the command for which a MetricScraper needs to be constructed.
func (msf *MetricScraperFactory) ForCommand(command core.Command) *MetricScraperFactory {
	msf.command = command
	return msf
}

// ForSingularEntity sets the factory to construct an entity specific MetricScraper.
func (msf *MetricScraperFactory) ForSingularEntity(entity string) *MetricScraperFactory {
	msf.singularEntityMetrics = true
	msf.entity = entity
	return msf
}

// WithScrapeInterval sets the scrape interval for the factory.
func (msf *MetricScraperFactory) WithScrapeInterval(interval uint64) *MetricScraperFactory {
	msf.scrapeIntervalMillisecond = interval
	return msf
}

// Construct constructs the MetricScraper for a particular Command and returns it.
func (msf *MetricScraperFactory) Construct() (MetricScraper, error) {
	switch msf.command {
	case core.RootCommand:
		return msf.constructSystemWideMetricScraper()
	case core.ContainerCommand:
		return msf.constructContainerMetricScraper()
	case core.ProcCommand:
		return msf.constructProcessMetricScraper()
	}
	return nil, errors.New("command not recognized")
}

func (msf *MetricScraperFactory) constructSystemWideMetricScraper() (MetricScraper, error) {
	return &systemWideMetrics{
		refreshRate: msf.scrapeIntervalMillisecond,
	}, nil
}

func (msf *MetricScraperFactory) constructContainerMetricScraper() (MetricScraper, error) {
	if msf.singularEntityMetrics {
		return msf.newSingluarContainerMetrics()
	}
	return msf.newContainerMetrics()
}

func (msf *MetricScraperFactory) newContainerMetrics() (*containerMetrics, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	cms := &containerMetrics{
		client:      cli,
		refreshRate: msf.scrapeIntervalMillisecond,
		metricBus:   make(chan container.ContainerMetrics),
	}

	return cms, nil
}

func (msf *MetricScraperFactory) newSingluarContainerMetrics() (*singularContainerMetrics, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	scms := &singularContainerMetrics{
		client:      cli,
		refreshRate: msf.scrapeIntervalMillisecond,
		cid:         msf.entity,
		metricBus:   make(chan container.PerContainerMetrics, 1),
	}

	return scms, nil
}

func (msf *MetricScraperFactory) constructProcessMetricScraper() (MetricScraper, error) {
	if msf.singularEntityMetrics {
		return msf.newSingluarProcessMetrics()
	}
	return msf.newProcessMetrics()
}

func (msf *MetricScraperFactory) newProcessMetrics() (*processMetrics, error) {
	pm := &processMetrics{
		refreshRate: msf.scrapeIntervalMillisecond,
		metricBus:   make(chan []*proc.Process, 1),
	}

	return pm, nil
}

func (msf *MetricScraperFactory) newSingluarProcessMetrics() (*singularProcessMetrics, error) {
	pid, err := strconv.ParseInt(msf.entity, 10, 32)
	if err != nil {
		return nil, err
	}
	spm := &singularProcessMetrics{
		refreshRate: msf.scrapeIntervalMillisecond,
		metricBus:   make(chan *process.Process, 1),
		pid:         int32(pid),
	}

	return spm, nil
}
