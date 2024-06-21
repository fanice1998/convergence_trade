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

// 計算 SMA Y軸資料
func calculateSMA(days int, data []common.KlineData) []interface{} {
	if days <= 0 || days > len(data) {
		return nil
	}

	// 宣告一個新的陣列，長度為資料的長度
	sma := make([]interface{}, len(data))

	for i := 0; i < len(data); i++ {

		var sum float64

		// 計算 SMA當前值
		// 若小於指定天數，則不計算
		if days > i {
			sma[i] = nil
			continue
		} else {
			for j := i - days; j < i; j++ {
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", data[j].Data[3]), 64)
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

type KlineDataChart struct {
	Name          string
	KlineData     []common.KlineData
	SMA           []int `json:"omitempty"`
	markKLineOpts []charts.SeriesOpts
	xAxis         []string
	yAxis         []opts.KlineData
}

// 初始化圖表名稱
func (k *KlineDataChart) SetName(name string) {
	k.Name = name
}

// 初始化設定 Data
func (k *KlineDataChart) SetData(kd []common.KlineData) {
	if len(kd) > 0 {
		k.KlineData = kd
	} else {
		panic("KlineData is empty")
	}
}

// 初始化設定 sma
func (k *KlineDataChart) SetSMA(sma []int) {
	for _, v := range sma {
		if v <= 0 {
			panic("SMA is less than 0")
		}
	}
	k.SMA = sma
}

// 設定圖表標記水平線
func (k *KlineDataChart) markLineChart(yline float64) {
	k.markKLineOpts = append(
		k.markKLineOpts,
		charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
			Name:  "markLine",
			YAxis: yline,
		}),
	)
}

// markLine to convergence 圖表標記線
func (k *KlineDataChart) markLine(i []interface{}) {
	k.markKLineOpts = append(k.markKLineOpts, charts.WithMarkLineNameCoordItemOpts(
		opts.MarkLineNameCoordItem{
			// Name:        "test",
			Coordinate0: []interface{}{i[0], i[1]},
			Coordinate1: []interface{}{i[2], i[3]},
		},
	))

}

// smaChart 建立 SMA 圖表
func (k *KlineDataChart) smaChart(days int) *charts.Line {
	line := charts.NewLine()
	x := k.xAxis
	y := make([]opts.LineData, len(k.KlineData))

	sma := calculateSMA(days, k.KlineData)

	for i := 0; i < len(k.KlineData); i++ {
		y[i] = opts.LineData{Value: sma[i]}
	}

	line.SetXAxis(x).AddSeries(fmt.Sprintf("SMA_%d", days), y, charts.WithSeriesAnimation(false))

	return line
}

// 主要視圖
func (k *KlineDataChart) mainChart() *charts.Kline {
	kline := charts.NewKLine()

	// 視圖比例
	// 是以百分比定位而不是以索引為定位
	// endPercent := len(k.KlineData)
	// startPercent := 0
	var startPercent, endPercent float32
	endPercent = float32(len(k.KlineData))
	if endPercent > 100 {
		startPercent = ((endPercent - 100) * 100) / endPercent
	}

	// kline 圖表設定選項
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: k.Name,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      startPercent,
			End:        endPercent,
			XAxisIndex: []int{0},
		}),
	)

	return kline
}

// 圖表運行主函數
func (k *KlineDataChart) Chart() *charts.Kline {

	// 建立 X, Y 軸資料
	k.xAxis = make([]string, 0)
	k.yAxis = make([]opts.KlineData, 0)
	for i := 0; i < len(k.KlineData); i++ {
		k.xAxis = append(k.xAxis, k.KlineData[i].Date)
		k.yAxis = append(k.yAxis, opts.KlineData{Value: k.KlineData[i].Data})
	}

	// 預設的 markLineOpts
	// 這裡的 markLineOpts 會覆蓋掉 SetSeriesOptions 的設定
	k.markKLineOpts = []charts.SeriesOpts{
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "green",
			Color0:       "red",
			BorderColor:  "darkgreen",
			BorderColor0: "darkred",
		}),
		//  之後應該可以用來做收斂圖形
		charts.WithMarkLineNameCoordItemOpts(opts.MarkLineNameCoordItem{
			Name:        "test",
			Coordinate0: []interface{}{"1577836800000", float64(130)},
			Coordinate1: []interface{}{"1580511600000", float64(220)},
		}),
	}

	// 建立主圖表物件
	main := k.mainChart()

	// 設定 markLineOpts
	SuportResistanceLines := TestSuportResistanceLine(
		func(d []common.KlineData) (closes []float64) {
			for _, v := range d {
				closes = append(closes, v.Data[1])
			}
			return closes
		}(k.KlineData),
	)
	if len(SuportResistanceLines) > 0 {
		for _, v := range SuportResistanceLines {
			k.markLineChart(v)
		}
	}

	// 繪製 e-chart
	main.SetXAxis(k.xAxis).AddSeries(k.Name, k.yAxis).
		SetSeriesOptions(
			k.markKLineOpts...,
		)

	// 插入 SMA 圖表
	if len(k.SMA) != 0 {
		for _, v := range k.SMA {
			main.Overlap(k.smaChart(v))
		}
	}

	return main
}

type KlineExamples struct{}

func (KlineExamples) Chart(kd []common.KlineData) {
	kline := KlineDataChart{}
	kline.SetName("ETHUSDT")
	kline.SetData(kd)
	kline.SetSMA([]int{7, 30, 90})

	page := components.NewPage()
	page.AddCharts(
		// klineDataZoomInside(kd),
		kline.Chart(),
		// klineStyle(kd),
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
