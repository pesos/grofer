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

	rw "github.com/mattn/go-runewidth"

	. "github.com/gizak/termui/v3"
)

type BarChart struct {
	Block
	BarColors    []Color
	LabelStyles  []Style
	NumStyles    []Style // only Fg and Modifier are used
	NumFormatter func(float64) string
	Data         []float64
	Labels       []string
	BarWidth     int
	BarGap       int
	MaxVal       float64
}

func NewBarChart() *BarChart {
	return &BarChart{
		Block:        *NewBlock(),
		BarColors:    Theme.BarChart.Bars,
		NumStyles:    Theme.BarChart.Nums,
		LabelStyles:  Theme.BarChart.Labels,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       1,
		BarWidth:     3,
	}
}

func (self *BarChart) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = GetMaxFloat64FromSlice(self.Data)
		if maxVal == 0 {
			maxVal = 1
		}
	}

	barXCoordinate := self.Inner.Min.X

	for i, data := range self.Data {
		// draw bar
		height := int((data / maxVal) * float64(self.Inner.Dy()-1))
		for x := barXCoordinate; x < MinInt(barXCoordinate+self.BarWidth, self.Inner.Max.X); x++ {
			for y := self.Inner.Max.Y - 2; y > (self.Inner.Max.Y-2)-height; y-- {
				c := NewCell(' ', NewStyle(ColorClear, SelectColor(self.BarColors, i)))
				buf.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		if i < len(self.Labels) {
			labelXCoordinate := barXCoordinate +
				int((float64(self.BarWidth) / 2)) -
				int((float64(rw.StringWidth(self.Labels[i])) / 2))
			buf.SetString(
				self.Labels[i],
				SelectStyle(self.LabelStyles, i),
				image.Pt(labelXCoordinate, self.Inner.Max.Y-1),
			)
		}

		// draw number
		numberXCoordinate := barXCoordinate + int((float64(self.BarWidth) / 2))
		if numberXCoordinate <= self.Inner.Max.X {
			buf.SetString(
				self.NumFormatter(data),
				NewStyle(
					SelectStyle(self.NumStyles, i+1).Fg,
					SelectColor(self.BarColors, i),
					SelectStyle(self.NumStyles, i+1).Modifier,
				),
				image.Pt(numberXCoordinate, self.Inner.Max.Y-2),
			)
		}

		barXCoordinate += (self.BarWidth + self.BarGap)
	}
}
