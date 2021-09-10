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

// CPUTableChart is a custom widget to display a CPU Table
type CPUTableChart struct {
	*ui.Block
	Data               []float64
	NumCores           int
	CellSize           int
	TopRow             int
	StatusColor        []ui.Color
	NumRows            int
	NumCols            int
	DefaultBorderColor ui.Color // indicates default border color
	ActiveBorderColor  ui.Color // indicates active border color
}

// NewCPUTableChart is a constructor for type CPUTableChart
func NewCPUTableChart() *CPUTableChart {
	return &CPUTableChart{
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
		DefaultBorderColor: ui.ColorCyan,
		ActiveBorderColor:  ui.ColorWhite,
	}
}

// Draw helps draw the CPU table Chart widget onto the UI buffer
func (c *CPUTableChart) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)
	w := c.Inner.Dx()
	h := c.Inner.Dy()
	numCols := int(w / (2 * c.CellSize))
	numRows := int(h / c.CellSize)
	c.NumCols = numCols
	c.NumRows = numRows
	xCoord := c.Inner.Min.X
	yCoord := c.Inner.Min.Y
	k := numCols * c.TopRow
	for i := 0; i < numRows && k < len(c.Data); i++ {
		for j := 0; j < numCols && k < len(c.Data); j++ {
			buf.Fill(
				ui.NewCell(' ', ui.NewStyle(ui.ColorClear, c.StatusColor[int(c.Data[k]/10)])),
				image.Rect(xCoord, yCoord, xCoord+2*c.CellSize, yCoord+c.CellSize),
			)
			buf.SetString(
				fmt.Sprintf("CPU%d", k),
				ui.NewStyle(c.StatusColor[int(c.Data[k]/10)], ui.ColorClear, ui.ModifierReverse),
				image.Pt(xCoord+2, yCoord+1),
			)
			buf.SetString(
				fmt.Sprintf("%.1f", c.Data[k]),
				ui.NewStyle(c.StatusColor[int(c.Data[k]/10)], ui.ColorClear, ui.ModifierReverse),
				image.Pt(xCoord+2, yCoord+2),
			)
			k++
			xCoord += 2 * c.CellSize
		}
		yCoord += c.CellSize
		xCoord = c.Inner.Min.X
	}
}

// ScrollUp moves the cursor one position upwards
func (c *CPUTableChart) ScrollUp() {
	if c.TopRow > 0 {
		c.TopRow--
	}
}

// ScrollDown moves the cursor one position downwards
func (c *CPUTableChart) ScrollDown() {
	if len(c.Data)-(c.TopRow+1)*c.NumCols > 0 {
		c.TopRow++
	}
}

// ScrollTop moves the cursor to the top
func (c *CPUTableChart) ScrollTop() {
	c.TopRow = 0
}

// ScrollBottom moves the cursor to the bottom
func (c *CPUTableChart) ScrollBottom() {
	c.TopRow = len(c.Data) / c.NumCols
}

// ScrollHalfPageUp moves the cursor half a page up
func (c *CPUTableChart) ScrollHalfPageUp() {
	c.TopRow = c.TopRow - (c.Inner.Dy())/(2*c.CellSize)
	if c.TopRow < 0 {
		c.TopRow = 0
	}
}

// ScrollHalfPageDown moves the cursor half a page down
func (c *CPUTableChart) ScrollHalfPageDown() {
	c.TopRow = c.TopRow + (c.Inner.Dy())/(2*c.CellSize)
	if c.TopRow > len(c.Data)/c.NumCols {
		c.TopRow = len(c.Data) / c.NumCols
	}
}

// ScrollPageUp moves the cursor a page up
func (c *CPUTableChart) ScrollPageUp() {
	c.TopRow = c.TopRow - (c.Inner.Dy())/(c.CellSize)
	if c.TopRow < 0 {
		c.TopRow = 0
	}
}

// ScrollPageDown moves the cursor a page down
func (c *CPUTableChart) ScrollPageDown() {
	c.TopRow = c.TopRow + (c.Inner.Dy())/(c.CellSize)
	if c.TopRow > len(c.Data)/c.NumCols {
		c.TopRow = len(c.Data) / c.NumCols
	}
}

// ScrollToIndex moves the cursor to a specified index
func (c *CPUTableChart) ScrollToIndex(idx int) {
	if idx >= 0 && idx <= len(c.Data)/c.NumCols {
		c.TopRow = idx
	}
}

// DisableCursor turns off the cursor and un highlights the table
func (c *CPUTableChart) DisableCursor() {
	c.BorderStyle.Fg = c.DefaultBorderColor
}

// EnableCursor turns on the cursor and highlights the table
func (c *CPUTableChart) EnableCursor() {
	c.BorderStyle.Fg = c.ActiveBorderColor
}

// ensure interface compliance.
var _ ui.Drawable = (*CPUTableChart)(nil)
