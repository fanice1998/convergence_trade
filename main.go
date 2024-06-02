package main

import (
	"fmt"

	"github.com/chart"
)

func main() {
	// data.DownloadData("kline", "SUIUSDT", "1h")
	// ch := chart.KlineExamples{}
	// ch.Chart()
	kd := chart.Kd
	fmt.Println(kd[0])
}
