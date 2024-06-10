package main

import (
	"testing"

	// 导入其他必要的包，例如用于读取数据的包

	"github.com/strategy"
)

// 使用您提供的KlineData结构体
type KlineData struct {
	Date string
	Data [4]float32 // 假设数组顺序为开盘价，收盘价，最低价，最高价
}

func TestMain(m *testing.T) {
	trader := strategy.Strategy{}
	trader.Signal = append(trader.Signal, strategy.Signal{})

	// trader.Signal.EMA.Status = true

	trader.Run()
}
