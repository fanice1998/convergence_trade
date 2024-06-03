package chart

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type KlineData struct {
	Date string
	Data [4]float32
}

var Kd = func() []KlineData {
	dirPath := "./SUIUSDT/1h"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var Kd []KlineData

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
				Kd = append(Kd, KlineData{
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
	return Kd
}()

func klineDataZoomInside() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(Kd); i++ {
		x = append(x, Kd[i].Date)
		y = append(y, opts.KlineData{Value: Kd[i].Data})
	}

	// 圖像比例
	// 起點為startCount，終點為endCount
	// 圖像比例 = (總數量 - 終點) * 100 / 總數量
	// 因沒辦法指定索引，故只能用百分比方式當索引
	var startCount float32 = 0.0
	var endCount float32 = 100.0
	startCount = (float32(len(Kd)) - endCount) * 100 / float32(len(Kd))

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
	markLineOpts := []charts.SeriesOpts{
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
				Show: false,
			},
		}),
	}

	// 繪製 e-chart
	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			markLineOpts...,
		)

	// calculateSMA(20, y)
	fmt.Println(calculateSMA(20, y))

	return kline
}

func calculateSMA(days int, data []opts.KlineData) []float32 {
	if days <= 0 || days >= len(data) {
		return nil
	}

	fmt.Printf("data: %v", len(data))

	sma := make([]float32, len(data))
	for i := 0; i < len(data)-1; i++ {
		sum := float32(0.0)
		if days > i {
			sma[i] = sum
			continue
		} else {
			for j := i - days; j < i; j++ {
				sum += data[j].Value.([4]float32)[3]
			}
		}
		sma[i] = sum / float32(days)
	}

	return sma
}

func klineStyle() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(Kd); i++ {
		x = append(x, Kd[i].Date)
		y = append(y, opts.KlineData{Value: Kd[i].Data})
	}

	totalCount := len(Kd)
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
	// f, err := os.Create("kline.html")
	if err != nil {
		panic(err)

	}
	defer f.Close()
	page.Render(io.MultiWriter(f))

	// 將 html 渲染後得結果儲存成圖片
	pwd, _ := os.Getwd()
	fileURL := "file://" + filepath.Join(pwd, "./examples/html/kline.html")
	saveImage(fileURL)
}

// 透過 chromedp 將 html 渲染後得結果儲存成圖片
func saveImage(fileURL string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte

	err := chromedp.Run(ctx, fullScreenshot(fileURL, &buf))
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("./examples/html/kline.png", buf, 0644); err != nil {
		panic(err)
	}
}

func fullScreenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.FullScreenshot(res, 90),
	}
}
