package chart

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/common"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/html2img"
)

// 可內部捲動的圖表
func klineDataZoomInside(kd []common.KlineData) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].Date)
		y = append(y, opts.KlineData{Value: kd[i].Data})
	}

	// 圖像比例
	// 起點為startCount，終點為endCount
	// 圖像比例 = (總數量 - 終點) * 100 / 總數量
	// 因沒辦法指定索引，故只能用百分比方式當索引
	// var startCount float32 = 0.0
	// var endCount float32 = 100.0
	// startCount = (float32(len(kd)) - endCount) * 100 / float32(len(kd))

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(inside)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type: "inside",
			// Start:      float32(startCount),
			// End:        float32(endCount),
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	// 繪製樣式
	// markLineOpts := make([]charts.SeriesOpts, 0)
	markKLineOpts := []charts.SeriesOpts{
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "green",
			Color0:       "red",
			BorderColor:  "darkgreen",
			BorderColor0: "darkred",
		}),
		// charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
		// 	Name: "max",
		// 	Type: "max",
		// }), charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
		// 	Name: "min",
		// 	Type: "min",
		// }), charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
		// 	Name:  "test",
		// 	YAxis: 1500,
		// }),
		charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
			Name:     "SMA",
			YAxis:    1500,
			ValueDim: "highest",
		}),
		//  之後應該可以用來做收斂圖形
		charts.WithMarkLineNameCoordItemOpts(opts.MarkLineNameCoordItem{
			Name:        "test",
			Coordinate0: []interface{}{"1577836800000", float64(300)},
			Coordinate1: []interface{}{"1580511600000", float64(2000)},
		}),
		// charts.WithMarkLineStyleOpts(opts.MarkLineStyle{
		// 	Label: &opts.Label{
		// 		Show: true,
		// 	},
		// }),
	}
	fmt.Println("calculateSMA length: ", len(calculateSMA(4, y)))
	fmt.Println("KlineData length: ", len(y))

	// 繪製 e-chart
	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			markKLineOpts...,
		)

	// calculateSMA(20, y)
	// fmt.Println(calculateSMA(20, y))

	line := charts.NewLine()
	lx := make([]string, 0)
	ly := make([]opts.LineData, 0)
	sma := calculateSMA(30, y)
	fmt.Println("SMA ", sma[30])
	for i := 0; i < len(kd); i++ {
		lx = append(lx, kd[i].Date)
		ly = append(ly, opts.LineData{Value: sma[i]})
	}
	line.SetXAxis(lx).AddSeries("line", ly)

	kline.Overlap(line)
	return kline
}

// 計算 SMA Y軸資料
func calculateSMA(days int, data []opts.KlineData) []float64 {
	if days <= 0 || days > len(data) {
		return nil
	}

	sma := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		// sum := float64(0.0)
		sum := 0.0
		if days > i {
			sma[i] = 0.0
			continue
		} else {
			for j := i - days; j < i; j++ {
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", data[j].Value.([4]float64)[3]), 64)
				sum += t
			}
		}
		sma[i] = sum / float64(days)
	}

	return sma
}

// 測試圖表樣式
func klineStyle(kd []common.KlineData) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].Date)
		y = append(y, opts.KlineData{Value: kd[i].Data})
	}

	totalCount := len(kd)
	startPercent := 0
	if totalCount > 100 {
		startPercent = ((totalCount - 100) * 100) / totalCount
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "different style",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start: float32(startPercent),
			End:   float32(startPercent + 100),
			// Start:      0,
			// End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "highest value",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "lowest value",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show: true,
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#ec0000",
				Color0:       "#00da3c",
				BorderColor:  "#8A0000",
				BorderColor0: "#008F28",
			}),
		)
	return kline
}

type KlineExamples struct{}

func (KlineExamples) Chart(kd []common.KlineData) {
	page := components.NewPage()
	page.AddCharts(
		klineDataZoomInside(kd),
		klineStyle(kd),
	)

	err := os.MkdirAll("./examples/html", 0777)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("./examples/html/kline.html")
	// f, err := os.Create("kline.html")
	if err != nil {
		panic(err)

	}
	defer f.Close()

	page.Render(io.MultiWriter(f))

	// 將 html 渲染後得結果儲存成圖片
	pwd, _ := os.Getwd()
	fileURL := "file://" + filepath.Join(pwd, "./examples/html/kline.html")
	html2img.SaveImage(fileURL)
}
