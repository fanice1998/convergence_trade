package zigzag

import (
	"github.com/common"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// KLine represents a single K-line data point
type KLine struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}

// ZigzagPoint represents a single point in the Zigzag indicator
type ZigzagPoint struct {
	Index int
	Price float64
}

// Zigzag represents the Zigzag indicator
type Zigzag struct {
	Length int
	Points []ZigzagPoint
}

// CalculateZigzag calculates the Zigzag indicator
func CalculateZigzag(data []KLine, length int) Zigzag {
	var points []ZigzagPoint
	var lastHigh, lastLow float64
	var lastHighIndex, lastLowIndex int

	for i, k := range data {
		if i == 0 {
			lastHigh = k.High
			lastLow = k.Low
			lastHighIndex = i
			lastLowIndex = i
			continue
		}

		if k.High > lastHigh {
			lastHigh = k.High
			lastHighIndex = i
		}

		if k.Low < lastLow {
			lastLow = k.Low
			lastLowIndex = i
		}

		if i-lastHighIndex >= length {
			points = append(points, ZigzagPoint{Index: lastHighIndex, Price: lastHigh})
			lastHigh = k.High
			lastHighIndex = i
		}

		if i-lastLowIndex >= length {
			points = append(points, ZigzagPoint{Index: lastLowIndex, Price: lastLow})
			lastLow = k.Low
			lastLowIndex = i
		}
	}

	return Zigzag{Length: length, Points: points}
}

// PlotZigzag plots the Zigzag indicator on a chart
func PlotZigzag(data []KLine, zigzag Zigzag) {
	p := plot.New()

	p.Title.Text = "Zigzag Indicator"
	p.X.Label.Text = "Index"
	p.Y.Label.Text = "Price"

	// Plot K-line data
	kline := make(plotter.XYs, len(data))
	for i, k := range data {
		kline[i].X = float64(i)
		kline[i].Y = k.Close
	}
	line, err := plotter.NewLine(kline)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	// Plot Zigzag points
	zigzagPoints := make(plotter.XYs, len(zigzag.Points))
	for i, point := range zigzag.Points {
		zigzagPoints[i].X = float64(point.Index)
		zigzagPoints[i].Y = point.Price
	}
	scatter, err := plotter.NewScatter(zigzagPoints)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Points(3)
	scatter.GlyphStyle.Color = plotter.DefaultLineStyle.Color
	p.Add(scatter)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, "zigzag.png"); err != nil {
		panic(err)
	}
}

func TestZigzagLine(data []common.KlineData) []ZigzagPoint {
	// 讀取指定路徑下的資料夾，並取得此位置路徑所有檔案
	fdata := make([]KLine, 1000)
	for i := 0; i < 1000; i++ {
		fdata[i] = KLine{
			Open:  data[i].Data[0],
			Close: data[i].Data[1],
			Low:   data[i].Data[2],
			High:  data[i].Data[3],
		}
	}

	zigzag := CalculateZigzag(fdata, 21)

	PlotZigzag(fdata, zigzag)
	return zigzag.Points
}
