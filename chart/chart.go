package chart

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type KlineData struct {
	Date string
	Data [4]float32
}

// var kd = func() []KlineData {
// 	var kData []KlineData
// 	os.OpenFile("./Data/ETHUSDT/1h/")

// }

var kd = func() []KlineData {
	dirPath := "./ETHUSDT/1h"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var kd []KlineData

	// 迴圈讀取每個檔案
	for _, file := range files {
		// check file is csv
		if filepath.Ext(file.Name()) != ".csv" {
			continue
		}
		csvFile, err := os.Open(filepath.Join(dirPath, file.Name()))
		if err != nil {
			panic(err)
		}
		defer csvFile.Close()

		csvReader := csv.NewReader(csvFile)
		records, err := csvReader.ReadAll()
		if err != nil {
			panic(err)
		}

		for _, d := range records {
			if d[0] == "open_time" {
				continue
			} else {
				// source data [open, high, low, close]
				// go-echarts kline data [open, close, low, high]
				open, _ := strconv.ParseFloat(d[1], 32)
				high, _ := strconv.ParseFloat(d[2], 32)
				low, _ := strconv.ParseFloat(d[3], 32)
				close, _ := strconv.ParseFloat(d[4], 32)
				kd = append(kd, KlineData{
					Date: d[0],
					Data: [4]float32{
						float32(open),
						float32(close),
						float32(low),
						float32(high),
					},
				})
			}
		}
	}
	return kd
}()

func klineDataZoomInside(markLineOpts []charts.SeriesOpts) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].Date)
		y = append(y, opts.KlineData{Value: kd[i].Data})
	}
	var startCount float32 = 0.0
	var endCount float32 = 100.0
	startCount = (float32(len(kd)) - endCount) * 100 / float32(len(kd))
	fmt.Println(startCount, endCount)

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
			Type:       "inside",
			Start:      float32(startCount),
			End:        float32(endCount),
			XAxisIndex: []int{0},
		}),
	)

	// 繪製樣式
	// markLineOpts := make([]charts.SeriesOpts, 0)
	markLineOpts = append(markLineOpts,
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "green",
			Color0:       "red",
			BorderColor:  "darkgreen",
			BorderColor0: "darkred",
		}),
		charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "max",
			Type: "max",
		}), charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "min",
			Type: "min",
		}), charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
			Name:  "test",
			YAxis: 1500,
		}), charts.WithMarkLineStyleOpts(opts.MarkLineStyle{
			Label: &opts.Label{
				Show: true,
			},
		}))

	// 繪製 e-chart
	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			markLineOpts...,
		)
	return kline
}

func klineStyle() *charts.Kline {
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
			Start:      float32(startPercent),
			End:        float32(startPercent + 100),
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

func (KlineExamples) Chart() {
	page := components.NewPage()
	page.AddCharts(
		klineDataZoomInside(),
		klineStyle(),
	)

	err := os.MkdirAll("./examples/html", 0777)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("./examples/html/kline.html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}
