package main

import (
	"fmt"
	"math"
	"testing"

	// 导入其他必要的包，例如用于读取数据的包
	"github.com/chart"
)

// 使用您提供的KlineData结构体
type KlineData struct {
	Date string
	Data [4]float32 // 假设数组顺序为开盘价，收盘价，最低价，最高价
}

// 检查是否符合策略条件
func CheckStrategyCondition(klines []chart.KlineData, index int) (bool, float32, float32, float32) {
	if index < 4 { // 前四根K线不检查
		return false, 0, 0, 0
	}

	// 计算前四根K线的体积平均
	var sumVolume float32
	for i := index - 4; i < index; i++ {
		sumVolume += float32(math.Abs(float64(klines[i].Data[1] - klines[i].Data[0]))) // 假设Data[3]是体积
	}
	averageVolume := sumVolume / 4
	// 当前K线体积是否大于平均体积的三倍
	if math.Abs(float64(klines[index].Data[1]-klines[index].Data[0])) > 3*float64(averageVolume) {
		var stopLoss, targetProfit float32
		ratioR := 2 // R的倍数设为2

		// 寻找前四根K线的最高点和最低点
		highest := klines[index-4].Data[2] // 假设Data[2]是最高价
		lowest := klines[index-4].Data[1]  // 假设Data[1]是最低价
		for i := index - 3; i < index; i++ {
			if klines[i].Data[2] > highest {
				highest = klines[i].Data[2]
			}
			if klines[i].Data[1] < lowest {
				lowest = klines[i].Data[1]
			}
		}

		entryPrice := klines[index].Data[1]                // 假设以开盘价进场
		if klines[index].Data[1] > klines[index].Data[0] { // 如果收盘价大于开盘价
			stopLoss = lowest
		} else {
			stopLoss = highest
		}

		risk := entryPrice - stopLoss
		if klines[index].Data[1] > klines[index].Data[0] {
			targetProfit = entryPrice + float32(ratioR)*risk
		} else {
			targetProfit = entryPrice - float32(ratioR)*risk
		}

		return true, entryPrice, stopLoss, targetProfit
	}
	return false, 0, 0, 0
}

func TestMain(t *testing.T) {
	// 假定klines已填充了历史K线数据
	var klines []chart.KlineData = chart.Kd
	// ...加载数据到klines...

	var orderList = []Order{}

	tcash := 1000.0
	var newIndex int
	for i, k := range klines {
		ok, entryPrice, stopLoss, targetProfit := CheckStrategyCondition(klines, i)
		if ok {
			fmt.Printf("在第%d根K线满足入场条件，止损价为%.2f，目标利润价为%.2f\n", i+1, stopLoss, targetProfit)
			// 此处可以加入执行开仓、设置止损止盈等逻辑
			if newIndex != 0 {
				newIndex += 1
			} else {
				newIndex = 1
			}

			orderList = append(orderList, Order{Index: newIndex, entryPrice: entryPrice, stopLoss: stopLoss, targetProfit: targetProfit})

		}
		for j, d := range orderList {
			if d.stopLoss > d.targetProfit && k.Data[3] > d.stopLoss {
				tcash -= tcash * 0.01
				orderList = append(orderList[:j], orderList[j+1:]...)
				fmt.Println("止损")
				fmt.Println("ID", d.Index)
				fmt.Println("当前现金：", tcash)
				break
			} else if d.targetProfit > d.stopLoss && k.Data[2] < d.stopLoss {
				tcash -= tcash * 0.01
				orderList = append(orderList[:j], orderList[j+1:]...)
				fmt.Println("止损")
				fmt.Println("ID", d.Index)
				fmt.Println("当前现金：", tcash)
				break
			} else if d.targetProfit > d.stopLoss && k.Data[3] > d.targetProfit {
				tcash += tcash * 0.01
				orderList = append(orderList[:j], orderList[j+1:]...)
				fmt.Println("盈利")
				fmt.Println("ID", d.Index)
				fmt.Println("当前现金：", tcash)
				break
			} else if d.targetProfit < d.stopLoss && k.Data[2] < d.targetProfit {
				tcash += tcash * 0.01
				fmt.Println("盈利")
				fmt.Println("ID", d.Index)
				fmt.Println("当前现金：", tcash)
				orderList = append(orderList[:j], orderList[j+1:]...)
				break
			}
		}

	}
}

type Order struct {
	Index        int
	entryPrice   float32
	stopLoss     float32
	targetProfit float32
}
