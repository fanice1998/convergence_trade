package convergence

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/common"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// KLine represents a single K-line data point
type ZKLine struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}

// ZigzagPoint represents a single point in the Zigzag indicator
type ZigzagPoint struct {
	Index int
	Price float64
}

// CalculateZigzag calculates the Zigzag indicator
func CalculateZigzag(data []ZKLine, threshold float64) []ZigzagPoint {
	var points []ZigzagPoint
	if len(data) == 0 {
		return points
	}

	lastHigh := data[0].High
	lastLow := data[0].Low
	lastHighIndex := 0
	lastLowIndex := 0
	points = append(points, ZigzagPoint{Index: 0, Price: data[0].Close})

	for i := 1; i < len(data); i++ {
		if data[i].High > lastHigh {
			lastHigh = data[i].High
			lastHighIndex = i
		}

		if data[i].Low < lastLow {
			lastLow = data[i].Low
			lastLowIndex = i
		}

		if (lastHigh-data[i].Low)/lastHigh >= threshold {
			points = append(points, ZigzagPoint{Index: lastHighIndex, Price: lastHigh})
			lastLow = data[i].Low
			lastLowIndex = i
		} else if (data[i].High-lastLow)/lastLow >= threshold {
			points = append(points, ZigzagPoint{Index: lastLowIndex, Price: lastLow})
			lastHigh = data[i].High
			lastHighIndex = i
		}
	}

	points = append(points, ZigzagPoint{Index: len(data) - 1, Price: data[len(data)-1].Close})
	return points
}

// PlotZigzag plots the Zigzag indicator on a chart
func PlotZigzag(data []ZKLine, zigzag []ZigzagPoint) {
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
	zigzagPoints := make(plotter.XYs, len(zigzag))
	for i, point := range zigzag {
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

func TestZigzag(t *testing.T) {
	// Example data with 100 K-line data points
	kd := ReadData("../../ETHUSDT/1h")
	data := make([]ZKLine, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = ZKLine{
			Open:  kd[i].Data[0],
			Close: kd[i].Data[1],
			High:  kd[i].Data[3],
			Low:   kd[i].Data[2],
		}
	}

	threshold := 0.05 // 5% threshold
	zigzag := CalculateZigzag(data, threshold)
	PlotZigzag(data, zigzag)
}

func ReadData(path string) (kd []common.KlineData) {
	// 讀取指定路徑下的資料夾，並取得此位置路徑所有檔案
	files, err := os.ReadDir(path)
	if err != nil {
		panic("Read directory Fail! message !")
	}

	for _, file := range files {
		// 如果檔案是CSV 則繼續，否則跳過
		if filepath.Ext(file.Name()) != ".csv" {
			continue
		}
		csvFile, err := os.Open(filepath.Join(path, file.Name()))
		if err != nil {
			fmt.Println("檔案讀取發生異常")
			panic(err)
		}
		defer csvFile.Close()

		csvReader := csv.NewReader(csvFile)
		records, err := csvReader.ReadAll()
		if err != nil {
			fmt.Println("CSV ReadALL() 讀取異常")
			panic(err)
		}

		for _, d := range records {
			if len(kd) == 1000 {
				break
			}
			if d[0] == "open_time" {
				continue
			} else {
				// source data [open, high, low, close]
				// go-echart kline data [open, close, low, high]
				open, _ := strconv.ParseFloat(d[1], 64)
				high, _ := strconv.ParseFloat(d[2], 64)
				low, _ := strconv.ParseFloat(d[3], 64)
				close, _ := strconv.ParseFloat(d[4], 64)
				kd = append(kd, common.KlineData{
					Date: d[0],
					Data: [4]float64{
						open,
						close,
						low,
						high,
					},
				})
			}
		}
	}
	return kd
}
