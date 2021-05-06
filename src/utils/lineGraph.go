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
// This particularar widget is inspired and borrowed from the implementation of https://github.com/cjbassi/gotop
package utils

import (
	"image"
	"sort"

	drawille "github.com/cjbassi/gotop/src/termui/drawille-go"
	ui "github.com/gizak/termui/v3"
)

// LineGraph implements a line graph of data points.
type LineGraph struct {
	*ui.Block

	Data   map[string][]float64
	Labels map[string]string

	HorizontalScale int
	MaxVal          float64

	LineColors       map[string]ui.Color
	DefaultLineColor ui.Color
}

// NewLineGraph creates and returns a lineGraph instance
func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block: ui.NewBlock(),

		Data:   make(map[string][]float64),
		Labels: make(map[string]string),

		HorizontalScale: 5,

		LineColors: make(map[string]ui.Color),
	}
}

func (l *LineGraph) Draw(buf *ui.Buffer) {
	l.Block.Draw(buf)
	// we render each data point on to the canvas then copy over the braille to the buffer at the end
	// fyi braille characters have 2x4 dots for each character
	c := drawille.NewCanvas()
	// used to keep track of the braille colors until the end when we render the braille to the buffer
	colors := make([][]ui.Color, l.Inner.Dx()+2)
	for i := range colors {
		colors[i] = make([]ui.Color, l.Inner.Dy()+2)
	}

	// sort the series so that overlapping data will overlap the same way each time
	seriesList := make([]string, len(l.Data))
	i := 0
	l.MaxVal = 1
	for seriesName := range l.Data {
		for _, val := range l.Data[seriesName] {
			if val > l.MaxVal {
				l.MaxVal = val
			}
		}
		seriesList[i] = seriesName
		i++
	}
	sort.Strings(seriesList)

	// draw lines in reverse order so that the first color defined in the colorscheme is on top
	for i := len(seriesList) - 1; i >= 0; i-- {
		seriesName := seriesList[i]
		seriesData := l.Data[seriesName]
		seriesLineColor, ok := l.LineColors[seriesName]
		if !ok {
			seriesLineColor = l.DefaultLineColor
		}

		// coordinates of last point
		lastY, lastX := -1, -1
		// assign colors to `colors` and lines/points to the canvas
		for i := len(seriesData) - 1; i >= 0; i-- {
			x := ((l.Inner.Dx() + 1) * 2) - 1 - (((len(seriesData) - 1) - i) * l.HorizontalScale)
			y := ((l.Inner.Dy() + 1) * 4) - 1 - int((float64((l.Inner.Dy())*4)-1)*(seriesData[i]/float64(l.MaxVal)))
			if x < 0 {
				// render the line to the last point up to the wall
				if x > 0-l.HorizontalScale {
					for _, p := range drawille.Line(lastX, lastY, x, y) {
						if p.X > 0 {
							c.Set(p.X, p.Y)
							colors[p.X/2][p.Y/4] = seriesLineColor
						}
					}
				}
				break
			}
			if lastY == -1 { // if this is the first point
				c.Set(x, y)
				colors[x/2][y/4] = seriesLineColor
			} else {
				c.DrawLine(lastX, lastY, x, y)
				for _, p := range drawille.Line(lastX, lastY, x, y) {
					colors[p.X/2][p.Y/4] = seriesLineColor
				}
			}
			lastX, lastY = x, y
		}

		// copy braille and colors to buffer
		for y, line := range c.Rows(c.MinX(), c.MinY(), c.MaxX(), c.MaxY()) {
			for x, char := range line {
				x /= 3 // idk why but it works
				if x == 0 {
					continue
				}
				if char != 10240 { // empty braille character
					buf.SetCell(
						ui.NewCell(char, ui.NewStyle(colors[x][y])),
						image.Pt(l.Inner.Min.X+x-1, l.Inner.Min.Y+y-1),
					)
				}
			}
		}
	}

	// renders key/label ontop
	for i, seriesName := range seriesList {
		if i+2 > l.Inner.Dy() {
			continue
		}
		seriesLineColor, ok := l.LineColors[seriesName]
		if !ok {
			seriesLineColor = l.DefaultLineColor
		}

		// render key ontop, but let braille be drawn over space characters
		str := seriesName + " " + l.Labels[seriesName]
		for k, char := range str {
			if char != ' ' {
				buf.SetCell(
					ui.NewCell(char, ui.NewStyle(seriesLineColor)),
					image.Pt(l.Inner.Min.X+2+k, l.Inner.Min.Y+i+1),
				)
			}
		}

	}
}
