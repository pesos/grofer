package utils_test

import (
	"testing"

	"github.com/pesos/grofer/src/utils"
)

func TestSortData(t *testing.T) {
	tests := []struct {
		inputVal    [][]string
		sortIdx     int
		sortAsc     bool
		sortCase    string
		expectedVal [][]string
	}{
		{
			inputVal: [][]string{
				{"CID1", "IMG1", "NAME1", "Up 5 Seconds", "Running", "12.34%", "4.20%"},
				{"CID2", "IMG2", "NAME2", "Up 1 Second", "Running", "69.69%", "4.20%"},
			},
			sortIdx:  5,
			sortAsc:  false,
			sortCase: "CONT",
			expectedVal: [][]string{
				{"CID2", "IMG2", "NAME2", "Up 1 Second", "Running", "69.69%", "4.20%"},
				{"CID1", "IMG1", "NAME1", "Up 5 Seconds", "Running", "12.34%", "4.20%"},
			},
		},
		{
			inputVal: [][]string{
				{"CID1", "IMG1", "NAME1", "Up 5 Seconds", "Running", "12.34%", "4.20%"},
				{"CID2", "IMG2", "NAME2", "Up 1 Second", "Running", "69.69%", "4.20%"},
			},
			sortIdx:  3,
			sortAsc:  true,
			sortCase: "",
			expectedVal: [][]string{
				{"CID2", "IMG2", "NAME2", "Up 1 Second", "Running", "69.69%", "4.20%"},
				{"CID1", "IMG1", "NAME1", "Up 5 Seconds", "Running", "12.34%", "4.20%"},
			},
		},
	}

	for _, test := range tests {
		utils.SortData(test.inputVal, test.sortIdx, test.sortAsc, test.sortCase)
		utils.Equals(t, test.inputVal, test.expectedVal)
	}

}
