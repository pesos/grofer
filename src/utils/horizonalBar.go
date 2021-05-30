package utils

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

type HorizontalBarChart struct {
	*ui.Block
	BarColors   []ui.Color
	LabelStyles []ui.Style
	Data        []float64
	Labels      []string
	MaxVal      float64
	BarGap      int
	BarWidth    int
	ColResizer  func()
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
	self.ColResizer()
	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64FromSlice(self.Data)
	}
	barYCoordinate := self.Inner.Min.Y

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
