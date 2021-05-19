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

	ui "github.com/gizak/termui/v3"
)

type BarChart struct {
	NumFormatter func(float64) string
	Labels       []string
	BarColors    []ui.Color
	LabelStyles  []ui.Style
	NumStyles    []ui.Style
	Data         []float64
	ui.Block
	BarWidth int
	BarGap   int
	MaxVal   float64
}

func NewBarChart() *BarChart {
	return &BarChart{
		Block:        *ui.NewBlock(),
		BarColors:    ui.Theme.BarChart.Bars,
		NumStyles:    ui.Theme.BarChart.Nums,
		LabelStyles:  ui.Theme.BarChart.Labels,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       1,
		BarWidth:     8,
	}
}

func (b *BarChart) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	maxVal := b.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64FromSlice(b.Data)
		if maxVal == 0 {
			maxVal = 1
		}
	}

	barXCoordinate := b.Inner.Min.X

	for i, data := range b.Data {
		// draw bar
		height := int((data / maxVal) * float64(b.Inner.Dy()-1))
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+b.BarWidth, b.Inner.Max.X); x++ {
			for y := b.Inner.Max.Y - 2; y > (b.Inner.Max.Y-2)-height; y-- {
				c := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(b.BarColors, i)))
				buf.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		if i < len(b.Labels) {
			labelXCoordinate := barXCoordinate +
				int((float64(b.BarWidth) / 2)) -
				int((float64(rw.StringWidth(b.Labels[i])) / 2))
			buf.SetString(
				b.Labels[i],
				ui.SelectStyle(b.LabelStyles, i),
				image.Pt(labelXCoordinate, b.Inner.Max.Y-1),
			)
		}

		// draw number
		numberXCoordinate := barXCoordinate + int((float64(b.BarWidth) / 2))
		if numberXCoordinate <= b.Inner.Max.X {
			buf.SetString(
				b.NumFormatter(data),
				ui.NewStyle(
					ui.SelectStyle(b.NumStyles, i+1).Fg,
					ui.SelectColor(b.BarColors, i),
					ui.SelectStyle(b.NumStyles, i+1).Modifier,
				),
				image.Pt(numberXCoordinate, b.Inner.Max.Y-2),
			)
		}

		barXCoordinate += (b.BarWidth + 4)
	}
}
