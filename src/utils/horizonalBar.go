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
	"image"

	ui "github.com/gizak/termui/v3"
)

type HorizontalBarChart struct {
	*ui.Block
	BarColors   []ui.Color // Custom bar colors
	LabelStyles []ui.Style // Styles label styles
	Data        []float64
	Labels      []string
	MaxVal      float64
	BarGap      int    // Gap between bars
	BarWidth    int    // Width of each bar
	ColResizer  func() // Function to resize bar
}

func NewHorizontalBarChart() *HorizontalBarChart {
	return &HorizontalBarChart{
		Block:       ui.NewBlock(),
		BarColors:   ui.Theme.BarChart.Bars,
		LabelStyles: ui.Theme.BarChart.Labels,
		BarWidth:    1,
		BarGap:      0,
		Labels:      []string{},
		ColResizer:  func() {},
	}
}

func (self *HorizontalBarChart) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)
	// Call function to resize columns depending on term size
	self.ColResizer()
	// Calculate maximum value if not given
	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64FromSlice(self.Data)
	}
	barYCoordinate := self.Inner.Min.Y
	// Draw the horizontal bars and print the labels.

	for i, data := range self.Data {
		barWidth := int((data / maxVal) * float64(self.Inner.Dx()))
		buf.Fill(
			ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(self.BarColors, i))),
			image.Rect(self.Inner.Min.X, barYCoordinate, barWidth+self.Inner.Min.X, barYCoordinate+self.BarWidth),
		)
		buf.SetString(
			self.Labels[i],
			ui.SelectStyle(self.LabelStyles, i),
			image.Pt(self.Inner.Min.X, barYCoordinate+self.BarWidth),
		)
		barYCoordinate += self.BarWidth + self.BarGap + 1
	}
}
