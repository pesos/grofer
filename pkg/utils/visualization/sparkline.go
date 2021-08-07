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

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []float64
	Title      string
	TitleStyle ui.Style
	LineColor  ui.Color
	MaxVal     float64
	MaxHeight  int // TODO
	Reverse    bool
}

// SparklineGroup is a renderable widget which groups together the given sparklines.
type SparklineGroup struct {
	*ui.Block
	Sparklines []*Sparkline
}

// NewSparkline returns a unrenderable single sparkline that needs to be added to a SparklineGroup
func NewSparkline() *Sparkline {
	return &Sparkline{
		TitleStyle: ui.Theme.Sparkline.Title,
		LineColor:  ui.Theme.Sparkline.Line,
		Reverse:    false,
	}
}

func NewSparklineGroup(sls ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{
		Block:      ui.NewBlock(),
		Sparklines: sls,
	}
}

func (s *SparklineGroup) Draw(buf *ui.Buffer) {
	s.Block.Draw(buf)

	sparklineHeight := s.Inner.Dy() / len(s.Sparklines)

	for i, sl := range s.Sparklines {
		heightOffset := (sparklineHeight * (i + 1))
		barHeight := sparklineHeight
		if i == len(s.Sparklines)-1 {
			heightOffset = s.Inner.Dy()
			barHeight = s.Inner.Dy() - (sparklineHeight * i)
		}
		if sl.Title != "" {
			barHeight--
		}

		maxVal := sl.MaxVal
		if maxVal == 0 {
			maxVal, _ = ui.GetMaxFloat64FromSlice(sl.Data)
		}

		if sl.Reverse {
			heightOffset -= sparklineHeight
		}

		// draw line
		index := len(sl.Data) - 1
		for j := s.Inner.Dx() - 1; index >= 0 && j >= 0; j-- {
			data := sl.Data[index]
			index--
			height := int((data / maxVal) * float64(barHeight))
			sparkChar := ui.BARS[len(ui.BARS)-1]
			for k := 0; k < height; k++ {
				yBarCoord := s.Inner.Min.Y - 1 + heightOffset - k
				if sl.Reverse {
					yBarCoord = s.Inner.Min.Y - 1 + heightOffset + k
				}
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(sl.LineColor)),
					image.Pt(j+s.Inner.Min.X, yBarCoord),
				)
			}
			if height == 0 {
				if sl.Reverse {
					sparkChar = ui.IRREGULAR_BLOCKS[3]
				} else {
					sparkChar = ui.BARS[1]
				}
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(sl.LineColor)),
					image.Pt(j+s.Inner.Min.X, s.Inner.Min.Y-1+heightOffset),
				)
			}
		}

		if sl.Title != "" {
			// draw title
			yTitleCoord := s.Inner.Min.Y - 1 + heightOffset - barHeight
			if sl.Reverse {
				yTitleCoord = s.Inner.Min.Y - 1 + heightOffset + sparklineHeight
			}
			buf.SetString(
				ui.TrimString(sl.Title, s.Inner.Dx()),
				sl.TitleStyle,
				image.Pt(s.Inner.Min.X, yTitleCoord),
			)
		}
	}
}

// ensure interface compliance.
var _ ui.Drawable = (*SparklineGroup)(nil)
