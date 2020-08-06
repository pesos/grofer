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

package utils_test

import (
	"testing"
	"time"

	"github.com/pesos/grofer/src/utils"
)

func TestRoundValues(t *testing.T) {
	tests := []struct {
		input               []float64
		expectedRoundedVals []float64
		expectedUnit        string
	}{
		{
			[]float64{999, 895},
			[]float64{999, 895},
			" ",
		},
		{
			[]float64{100000, 1000},
			[]float64{100, 1},
			" per thousand ",
		},
		{
			[]float64{10000000, 1000},
			[]float64{10, 0},
			" per million ",
		},
		{
			[]float64{100000000, 100000000000},
			[]float64{0.1, 100},
			" per trillion ",
		},
	}

	for _, test := range tests {
		testRoundedVals, testUnit := utils.RoundValues(test.input[0], test.input[1])
		utils.Equals(t, test.expectedRoundedVals, testRoundedVals)
		utils.Equals(t, test.expectedUnit, testUnit)
	}
}

func TestGetInMB(t *testing.T) {
	tests := []struct {
		inputVal    uint64
		precision   int
		expectedVal float64
	}{
		{1234567, 1, 1.2},
		{123456789, 2, 123.46},
		{123456789, 6, 123.456789},
		{0, 2, 0},
	}

	for _, test := range tests {
		testVal := utils.GetInMB(test.inputVal, test.precision)
		utils.Equals(t, testVal, test.expectedVal)
	}
}

func TestGetDateFromUnix(t *testing.T) {
	t1 := time.Unix(10000000, 0)
	date1 := t1.Format(time.RFC822)

	t2 := time.Unix(0, 0)
	date2 := t2.Format(time.RFC822)

	t3 := time.Unix(1596652055, 0)
	date3 := t3.Format(time.RFC822)

	t4 := time.Unix(9999999999, 0)
	date4 := t4.Format(time.RFC822)

	tests := []struct {
		inputVal    int64
		expectedVal string
	}{
		{10000000, date1},
		{0, date2},
		{1596652055, date3},
		{9999999999, date4},
	}

	for _, test := range tests {
		testVal := utils.GetDateFromUnix(test.inputVal)
		utils.Equals(t, testVal, test.expectedVal)
	}
}
