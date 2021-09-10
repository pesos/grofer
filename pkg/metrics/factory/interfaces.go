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
	"github.com/pesos/grofer/pkg/core"
)

// MetricScraper scrapes metrics of some form and serves it based on
// the implementation.
type MetricScraper interface {
	// Serve serves the metrics to a 'sink', which can be a TUI or
	// logic that exports these served metrics to either a file or
	// maybe even served over an endpoint (some day).
	Serve(opts ...Option) error
	// SetSink sets the Sink that consumes the metrics produced
	// by the MetricScraper.
	SetSink(core.Sink)
}
