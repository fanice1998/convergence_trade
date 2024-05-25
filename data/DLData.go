package data

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Binance kline data 響應的 xml 格式
type ListBucketResult struct {
	XMLName        xml.Name `xml:"ListBucketResult"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Name           string   `xml:"Name"`
	Prefix         string   `xml:"Prefix"`
	Marker         string   `xml:"Marker"`
	MaxKeys        string   `xml:"MaxKeys"`
	Delimiter      string   `xml:"Delimiter"`
	IsTruncated    string   `xml:"IsTruncated"`
	CommonPrefixes []struct {
		Text   string `xml:",chardata"`
		Prefix string `xml:"Prefix"`
	} `xml:"CommonPrefixes"`
	Contents []struct {
		Text         string `xml:",chardata"`
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         string `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	} `xml:"Contents"`
}

type DLData struct {
	Name string
	URL  string
}

// xml 網址 https://s3-ap-northeast-1.amazonaws.com/data.binance.vision?delimiter=/&prefix=data/futures/um/monthly/klines/APEUSDT/12h/
// 下載網址 https://data.binance.vision/?prefix=data/futures/um/monthly/klines/APEUSDT/12h/

func (d *DLData) Request() ([]byte, error) {
	response, err := http.Get(d.URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 請求 https://s3-ap-northeast-1.amazonaws.com/data.binance.vision 取得 xml
func Test(prefix ...string) {
	// 如果 prefix 為空，則代表使用預設值 "data/futures/um/monthly/klines/ENSUSDT/1h"
	if len(prefix) == 0 {
		prefix = append(prefix, "data/futures/um/monthly/klines/ENSUSDT/1h")
	} else if prefix[0] == "kline" && len(prefix) == 1 {
		prefix[0] = "data/futures/um/monthly/klines"
	} else if len(prefix) == 2 {
		prefix[0] = "data/futures/um/monthly/klines"
		prefix[0] = fmt.Sprintf("%s/%s", prefix[0], prefix[1])
	} else if len(prefix) == 3 {
		prefix[0] = "data/futures/um/monthly/klines"
		prefix[0] = fmt.Sprintf("%s/%s/%s", prefix[0], prefix[1], prefix[2])
	} else {
		fmt.Println("輸入格是錯誤!")
		fmt.Println("請輸入格式: kline, symbol, time")
		return
	}

	d := DLData{
		Name: "binance",
		URL:  fmt.Sprintf("https://s3-ap-northeast-1.amazonaws.com/data.binance.vision?delimiter=/&prefix=%s/", prefix[0]),
	}
	fmt.Println(d.URL)
	body, err := d.Request()
	if err != nil {
		fmt.Println(err)
	}

	// 將響應的 xml 解析
	var result ListBucketResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("輸入格是錯誤!")
		fmt.Println("請輸入格式: kline, symbol, time")
		fmt.Println("Error unmarshalling XML: ", err)
		return
	}

	// 如果 result.CommonPrefixes 不為空，則代表有資料
	// 如果 result.CommonPrefixes 為空，則代表無資料，並且顯示 result.Contents 的資料
	if result.CommonPrefixes != nil {
		for _, prefix := range result.CommonPrefixes {
			fmt.Println(prefix.Prefix)
		}
	} else {
		for _, content := range result.Contents {
			fileName := strings.Split(content.Key, "/")
			// fmt.Println(content.Key)
			// fmt.Println(fileName[len(fileName)-1])
			downloadData(strings.Join(append([]string{"https://data.binance.vision"}, content.Key), "/"), fileName[len(fileName)-1])
		}
	}

}

// 下載檔案
func downloadData(DLURL, fileName string) error {
	// 請求下載網址
	resp, err := http.Get(DLURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 建立多層目錄名稱
	// example: ETHUSDT/1h/ETHUSDT-1h-2020-01.zip => ETHUSDT/1h/
	dirpath := strings.Join(strings.Split(fileName, "-")[0:2], "/")
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
		return err
	}

	// 建立檔案
	out, err := os.Create(strings.Join([]string{dirpath, fileName}, "/"))
	if err != nil {
		return err
	}
	defer out.Close()

	// 下載檔案
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
