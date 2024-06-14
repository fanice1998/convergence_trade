package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chart"
	"github.com/common"
)

func readData(path string) (kd []common.KlineData) {
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

// 執行程式
func main() {
	// data.DownloadData("kline", "ETHUSDT", "1h")

	var kd = readData("./ETHUSDT/1h")

	ch := chart.KlineExamples{}
	ch.Chart(kd)

}
