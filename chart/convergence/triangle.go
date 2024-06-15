package convergence

import (
	"fmt"

	"github.com/common"
)

// KLine represents a single K-line data point
type KLine struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}

// TrendLine represents a trend line with a slope and intercept
type TrendLine struct {
	Slope     float64
	Intercept float64
}

// CalculateTrendLine calculates the trend line using linear regression
func CalculateTrendLine(data []common.KlineData, isUpper bool) TrendLine {
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(data))

	for i, k := range data {
		x := float64(i)
		y := k.Data[2]
		if isUpper {
			y = k.Data[3]
		}
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	return TrendLine{Slope: slope, Intercept: intercept}
}

// DetectConvergingTriangle detects if there is a converging triangle pattern
func DetectConvergingTriangle(data []common.KlineData) bool {
	upperTrendLine := CalculateTrendLine(data, true)
	lowerTrendLine := CalculateTrendLine(data, false)

	for i := range data {
		x := float64(i)
		upperY := upperTrendLine.Slope*x + upperTrendLine.Intercept
		lowerY := lowerTrendLine.Slope*x + lowerTrendLine.Intercept
		if upperY <= lowerY {
			return true
		}
	}
	return false
}

// DetectMultipleConvergingTriangles detects multiple converging triangle patterns in a large dataset
func DetectMultipleConvergingTriangles(data []common.KlineData, windowSize int) []int {
	var detectedIndices []int

	for i := 0; i <= len(data)-windowSize; i++ {
		window := data[i : i+windowSize]
		if DetectConvergingTriangle(window) {
			detectedIndices = append(detectedIndices, i)
		}
	}

	return detectedIndices
}

func TestConvergenceLine(prices []common.KlineData) []int {
	// Example data with 1000 K-line data points
	// data := make([]KLine, 1000)
	// for i := 0; i < 1000; i++ {
	//     data[i] = KLine{
	//         Open:  float64(100 + i%10),
	//         Close: float64(105 + i%10),
	//         High:  float64(110 + i%10),
	//         Low:   float64(95 + i%10),
	//     }
	// }

	windowSize := 50
	detectedIndices := DetectMultipleConvergingTriangles(prices, windowSize)

	fmt.Printf("Detected converging triangles at indices: %v\n", detectedIndices)
	return detectedIndices
}
