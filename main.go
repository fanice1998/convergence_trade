package main

import (
	"github.com/chart"
)

func main() {
	// data.DownloadData("kline", "ETHUSDT", "1h")
	ch := chart.KlineExamples{}
	ch.Chart()
}
