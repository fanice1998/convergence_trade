package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/chart"
)

func TestDownloadData(t *testing.T) {
	dirPath := "./ETHUSDT/1h"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var kd []chart.KlineData

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
				open, _ := strconv.ParseFloat(d[1], 32)
				high, _ := strconv.ParseFloat(d[2], 32)
				low, _ := strconv.ParseFloat(d[3], 32)
				close, _ := strconv.ParseFloat(d[4], 32)
				kd = append(kd, chart.KlineData{
					Date: d[0],
					Data: [4]float32{
						float32(open),
						float32(high),
						float32(low),
						float32(close),
					},
				})
			}
		}
	}
}
