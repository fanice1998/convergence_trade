package main

import (
	"github.com/chart"
	"github.com/data"
)

func main() {
	data.DownloadData("kline", "SUIUSDT", "1h")
	ch := chart.KlineExamples{}
	ch.Chart()

}
