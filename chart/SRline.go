package chart

import (
	"fmt"
	"math"

	"github.com/muesli/clusters"
)

func findExtrema(prices []float64) ([]int, []int) {
	var maxima, minima []int

	for i := 1; i < len(prices)-1; i++ {
		if prices[i] > prices[i-1] && prices[i] > prices[i+1] {
			maxima = append(maxima, i)
		}
		if prices[i] < prices[i-1] && prices[i] < prices[i+1] {
			minima = append(minima, i)
		}
	}

	return maxima, minima
}

type point struct {
	index int
	price float64
}
type Points []point

func (p Points) Dims() int { return 1 }

func (p Points) Distance(i, j int) float64 {
	return math.Abs(p[i].price - p[j].price)
}

func (p Points) Point(i int) []float64 {
	return []float64{p[i].price}
}

func (p point) Coordinates() clusters.Coordinates {
	return clusters.Coordinates{p.price}
}

func DBSCAN(data Points, epsilon float64, minPoints int) [][]point {
	var clusters [][]point
	visited := make([]bool, len(data))

	for i := range data {
		if visited[i] {
			continue
		}
		visited[i] = true
		neighbors := regionQuery(data, i, epsilon)
		if len(neighbors) < minPoints {
			continue
		}

		var cluster []point
		cluster = append(cluster, data[i])
		for len(neighbors) > 0 {
			n := neighbors[0]
			neighbors = neighbors[1:]
			if !visited[n] {
				visited[n] = true
				neighbors2 := regionQuery(data, n, epsilon)
				if len(neighbors2) >= minPoints {
					neighbors = append(neighbors, neighbors2...)
				}
			}
			cluster = append(cluster, data[n])
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

func regionQuery(data Points, i int, eps float64) (averages []int) {
	var neighbors []int
	for j := range data {
		if math.Abs(data[i].price-data[j].price) <= eps {
			neighbors = append(neighbors, j)
		}
	}
	return neighbors
}

func TestSuportResistanceLine(prices []float64) (averages []float64) {
	// prices := []float64{100, 102, 101, 104, 103, 99, 98, 97, 100, 105, 104, 106, 107, 106}

	maxima, minima := findExtrema(prices)
	fmt.Println("Maxima indices: ", maxima)
	fmt.Println("Minima indices: ", minima)

	var extremaPoints Points
	for _, i := range maxima {
		extremaPoints = append(extremaPoints, point{i, prices[i]})
	}

	for _, i := range minima {
		extremaPoints = append(extremaPoints, point{i, prices[i]})
	}

	epsilon := 2.0
	minPoints := 2
	clusters := DBSCAN(extremaPoints, epsilon, minPoints)

	for _, cluster := range clusters {
		if len(cluster) > 0 {
			sum := 0.0
			for _, p := range cluster {
				sum += p.price
			}

			average := sum / float64(len(cluster))
			fmt.Printf("Support/Resistance level: %.2f\n", average)
			averages = append(averages, average)
		}
	}
	return averages
}
