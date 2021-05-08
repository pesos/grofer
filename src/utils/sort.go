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

package utils

import (
	"sort"
	"strconv"
)

// SortData helps sort table rows. It sorts the table based on values given in the sortIdx column and sorts ascending if sortAsc is true. sortCase is set to identify the set of 'less' functions to use to sort the selected column by.
func SortData(data [][]string, sortIdx int, sortAsc bool, sortCase string) {

	// Define less functions
	intSort := func(i, j int) bool {
		x, _ := strconv.Atoi(data[i][sortIdx])
		y, _ := strconv.Atoi(data[j][sortIdx])
		if sortAsc {
			return x < y
		}
		return x > y
	}

	strSort := func(i, j int) bool {
		if sortAsc {
			return data[i][sortIdx] < data[j][sortIdx]
		}
		return data[i][sortIdx] > data[j][sortIdx]
	}

	floatSort := func(i, j int) bool {
		x1 := data[i][sortIdx]
		y1 := data[j][sortIdx]
		x, _ := strconv.ParseFloat(x1[:len(x1)-1], 32)
		y, _ := strconv.ParseFloat(y1[:len(y1)-1], 32)
		if sortAsc {
			return x < y
		}
		return x > y
	}

	// Set function map
	sortFuncs := make(map[int]func(i, j int) bool)
	switch sortCase {
	case "PROCS":
		sortFuncs = map[int]func(i, j int) bool{
			0: intSort,
			1: strSort,
			2: floatSort,
			3: floatSort,
			4: strSort,
			5: strSort,
			6: strSort,
			7: intSort,
		}
	case "CONTAINER":
		sortFuncs = map[int]func(i, j int) bool{
			0: strSort,
			1: strSort,
			2: strSort,
			3: strSort,
			4: strSort,
			5: floatSort,
			6: floatSort,
		}

	default:
		sortFuncs[sortIdx] = strSort
	}

	// Sort data
	sort.Slice(data, sortFuncs[sortIdx])
}
