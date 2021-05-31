// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

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

func (self *SparklineGroup) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	sparklineHeight := self.Inner.Dy() / len(self.Sparklines)

	for i, sl := range self.Sparklines {
		heightOffset := (sparklineHeight * (i + 1))
		barHeight := sparklineHeight
		if i == len(self.Sparklines)-1 {
			heightOffset = self.Inner.Dy()
			barHeight = self.Inner.Dy() - (sparklineHeight * i)
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
		for j := 0; j < len(sl.Data) && j < self.Inner.Dx(); j++ {
			data := sl.Data[j]
			height := int((data / maxVal) * float64(barHeight))
			sparkChar := ui.IRREGULAR_BLOCKS[15]
			for k := 0; k < height; k++ {
				yBarCoord := self.Inner.Min.Y - 1 + heightOffset - k
				if sl.Reverse {
					yBarCoord = self.Inner.Min.Y - 1 + heightOffset + k
				}
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, yBarCoord),
				)

			}
			if height == 0 {
				if sl.Reverse {
					sparkChar = ui.IRREGULAR_BLOCKS[3]
				} else {
					sparkChar = ui.IRREGULAR_BLOCKS[12]
				}
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset),
				)
			}
		}

		if sl.Title != "" {
			// draw title
			yTitleCoord := self.Inner.Min.Y - 1 + heightOffset - barHeight
			if sl.Reverse {
				yTitleCoord = self.Inner.Min.Y - 1 + heightOffset + sparklineHeight
			}
			buf.SetString(
				ui.TrimString(sl.Title, self.Inner.Dx()),
				sl.TitleStyle,
				image.Pt(self.Inner.Min.X, yTitleCoord),
			)
		}
	}
}
