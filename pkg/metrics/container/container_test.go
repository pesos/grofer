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
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/pesos/grofer/pkg/utils"
)

func TestGetPerCPUPercents(t *testing.T) {
	tests := []struct {
		inputStats     types.StatsJSON
		expectedOutput []string
	}{
		{
			inputStats: types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							PercpuUsage: []uint64{0},
						},
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							PercpuUsage: []uint64{0, 0},
						},
					},
				},
			},
			expectedOutput: []string{"NA", "NA"},
		},
	}

	for _, test := range tests {
		testVal := getPerCPUPercents(&test.inputStats)
		utils.Equals(t, testVal, test.expectedOutput)
	}
}

func TestGetCPUPercent(t *testing.T) {
	tests := []struct {
		inputStats     types.StatsJSON
		expectedOutput float64
	}{
		{
			inputStats: types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							PercpuUsage: []uint64{0},
						},
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							PercpuUsage: []uint64{0, 0},
						},
					},
				},
			},
			expectedOutput: 0,
		},
	}

	for _, test := range tests {
		testVal := getCPUPercent(&test.inputStats)
		utils.Equals(t, testVal, test.expectedOutput)
	}
}
