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
	"fmt"
	"image"

	ui "github.com/gizak/termui/v3"
)

// Custom widget to print a CPU Table

type CpuTableChart struct {
	*ui.Block
	Data        []float64
	NumCores    int
	CellSize    int
	TopRow      int
	StatusColor []ui.Color
	NumRows     int
	NumCols     int
}

func NewCpuTableChart() *CpuTableChart {
	return &CpuTableChart{
		Block:    ui.NewBlock(),
		Data:     []float64{0},
		NumCores: 0,
		CellSize: 4,
		TopRow:   0,
		NumRows:  0,
		NumCols:  0,
		StatusColor: []ui.Color{
			ui.Color(46),
			ui.Color(82),
			ui.Color(154),
			ui.Color(191),
			ui.Color(190),
			ui.Color(226),
			ui.Color(220),
			ui.Color(214),
			ui.Color(202),
			ui.Color(196),
			ui.Color(160),
		},
	}
}

func (chart *CpuTableChart) Draw(buf *ui.Buffer) {
	chart.Block.Draw(buf)
	w := chart.Inner.Dx()
	h := chart.Inner.Dy()
	numCols := int(w / (2 * chart.CellSize))
	numRows := int(h / chart.CellSize)
	chart.NumCols = numCols
	chart.NumRows = numRows
	xCoord := chart.Inner.Min.X
	yCoord := chart.Inner.Min.Y
	k := numCols * chart.TopRow
	for i := 0; i < numRows && k < len(chart.Data); i++ {
		for j := 0; j < numCols && k < len(chart.Data); j++ {
			buf.Fill(
				ui.NewCell(' ', ui.NewStyle(ui.ColorClear, chart.StatusColor[int(chart.Data[k]/10)])),
				image.Rect(xCoord, yCoord, xCoord+2*chart.CellSize, yCoord+chart.CellSize),
			)
			buf.SetString(
				fmt.Sprintf("CPU%d", k),
				ui.NewStyle(chart.StatusColor[int(chart.Data[k]/10)], ui.ColorClear, ui.ModifierReverse),
				image.Pt(xCoord+2, yCoord+1),
			)
			buf.SetString(
				fmt.Sprintf("%.1f", chart.Data[k]),
				ui.NewStyle(chart.StatusColor[int(chart.Data[k]/10)], ui.ColorClear, ui.ModifierReverse),
				image.Pt(xCoord+2, yCoord+2),
			)
			k++
			xCoord += 2 * chart.CellSize
		}
		yCoord += chart.CellSize
		xCoord = chart.Inner.Min.X
	}
}

func (c *CpuTableChart) ScrollUp() {
	if c.TopRow > 0 {
		c.TopRow--
	}
}

func (c *CpuTableChart) ScrollDown() {
	if len(c.Data)-(c.TopRow+1)*c.NumCols > 0 {
		c.TopRow++
	}
}

// ensure interface compliance.
var _ ui.Drawable = (*CpuTableChart)(nil)
