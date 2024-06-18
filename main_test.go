package main

import (
	"os"
	"testing"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Test_main(t *testing.T) {
	// 创建一个新的K线图
	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "蜡烛图和折线图重叠"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "X"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Y"}),
	)

	// 添加K线图数据
	kline.AddSeries("KLine", []opts.KlineData{
		{Value: []float64{20, 30, 10, 35}},
		{Value: []float64{40, 35, 30, 55}},
		{Value: []float64{33, 38, 33, 40}},
	})

	// 创建一个新的折线图
	line := charts.NewLine()
	lineData := []opts.LineData{
		{Value: 30},
		{Value: nil}, // 假设这里的数据被清洗掉了，使用nil表示
		{Value: 20},
	}
	line.AddSeries("Line", lineData).SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
	)

	// 将折线图和K线图重叠
	kline.Overlap(line)

	// 创建页面
	page := components.NewPage()
	page.AddCharts(kline)

	// 将图表渲染成HTML文件
	f, _ := os.Create("kline_line_overlap_with_nil.html")
	page.Render(f)
}
