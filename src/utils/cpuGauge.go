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
	BarHeight   int
	BarWidth    int
	BarGap      int
	ColResizer  func()
}

func NewCpuGauge() *CpuGauge {
	return &CpuGauge{
		Block:       ui.NewBlock(),
		Values:      []float64{},
		StatusColor: []ui.Color{ui.Color(46), ui.Color(82), ui.Color(154), ui.Color(191), ui.Color(190), ui.Color(226), ui.Color(220), ui.Color(214), ui.Color(202), ui.Color(196), ui.Color(160)},
		BarGap:      0,
		BarWidth:    4,
		BarHeight:   2,
		ColResizer:  func() {},
	}
}

func (c *CpuGauge) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)
	c.ColResizer()
	width := c.Inner.Dx()
	for i, val := range c.Values {
		x := c.Inner.Min.X + (width-c.BarWidth*10)/2
		y := c.Inner.Min.Y + (c.BarGap+c.BarHeight)*i
		for j := 0; j < int(roundOffNearestTen(val, 10)); j++ {
			_x := x + c.BarWidth*j
			buf.Fill(
				ui.NewCell(' ', ui.NewStyle(ui.ColorClear, c.StatusColor[j])),
				image.Rect(_x, y, _x+c.BarWidth, y+c.BarHeight),
			)
		}
		buf.SetString(
			fmt.Sprintf("CPU%d", i),
			ui.NewStyle(ui.ColorClear, ui.ColorClear),
			image.Pt(x-5, y),
		)
		buf.SetString(
			c.Labels[i],
			ui.NewStyle(ui.ColorClear, ui.ColorClear),
			image.Pt(x+1+10*c.BarWidth, y),
		)
	}
}
