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

// FactoryOption is used to inject command specific configuration.
type FactoryOption func(MetricScraper)

// WithAllAs sets the all flag value for the ContainerCommand.
func WithAllAs(all bool) FactoryOption {
	return func(ms MetricScraper) {
		cms := ms.(*containerMetrics)
		cms.all = all
	}
}

// WithCPUInfoAs sets the cpuinfo flag value for the RootCommand.
func WithCPUInfoAs(cpuInfo bool) FactoryOption {
	return func(ms MetricScraper) {
		swm := ms.(*systemWideMetrics)
		swm.cpuInfo = cpuInfo
	}
}
