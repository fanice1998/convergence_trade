package common

type KlineData struct {
	Date string
	// source data [open, high, low, close]
	// go-echart kline data [open, close, low, high]
	Data [4]float64
}
