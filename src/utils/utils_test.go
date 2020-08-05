package utils_test

import (
	"testing"

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
	tests := []struct {
		inputVal    int64
		expectedVal string
	}{
		{10000000, "26 Apr 70 23:16 IST"},
		{0, "01 Jan 70 05:30 IST"},
		{1596652055, "05 Aug 20 23:57 IST"},
		{9999999999, "20 Nov 86 23:16 IST"},
	}

	for _, test := range tests {
		testVal := utils.GetDateFromUnix(test.inputVal)
		utils.Equals(t, testVal, test.expectedVal)
	}
}
