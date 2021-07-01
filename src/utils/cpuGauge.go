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

type CpuGauge struct {
	*ui.Block
	Labels      []string
	Values      []float64
	StatusColor []ui.Color
	//BarHeight   int
	BarWidth   int
	BarGap     int
	ColResizer func()
}

func NewCpuGauge() *CpuGauge {
	return &CpuGauge{
		Block:       ui.NewBlock(),
		Values:      []float64{},
		StatusColor: []ui.Color{ui.Color(46), ui.Color(82), ui.Color(154), ui.Color(191), ui.Color(190), ui.Color(226), ui.Color(220), ui.Color(214), ui.Color(202), ui.Color(196), ui.Color(160)},
		BarGap:      0,
		BarWidth:    2,
		//BarHeight:   2,
		ColResizer: func() {},
	}
}

func (c *CpuGauge) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)
	c.ColResizer()
	for i, val := range c.Values {
		width := int(val) * c.Inner.Dx() / 100
		height := c.Inner.Min.Y + i*(c.BarWidth+c.BarGap)
		buf.Fill(
			ui.NewCell(' ', ui.NewStyle(ui.ColorClear, c.StatusColor[int(val)/10])),
			image.Rect(c.Inner.Min.X, height, c.Inner.Min.X+width, height+c.BarWidth),
		)
		label := fmt.Sprintf("CPU%d %0.2f%%", i, val)
		label_pos := c.Inner.Min.X + c.Inner.Dx()/2 - len(label)/2
		for i, x := range label {
			bg := buf.GetCell(image.Pt(label_pos+i, height)).Style.Bg
			if bg == ui.ColorClear {
				buf.SetCell(ui.NewCell(x, ui.NewStyle(c.StatusColor[int(val)/10], ui.ColorClear)), image.Pt(label_pos+i, height))
			} else {
				buf.SetCell(ui.NewCell(x, ui.NewStyle(c.StatusColor[int(val)/10], ui.ColorClear, ui.ModifierReverse)), image.Pt(label_pos+i, height))
			}
		}
	}
}
