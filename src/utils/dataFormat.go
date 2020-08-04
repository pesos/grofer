package utils

import (
	"math"
	"time"
)

var (
	K = math.Pow(10, 3)
	M = math.Pow(10, 6)
	G = math.Pow(10, 9)
	T = math.Pow(10, 12)
)

func roundOffNearestTen(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

func RoundValues(num1, num2 float64) ([]float64, string) {
	nums := []float64{}
	var units string
	var n float64
	if num1 > num2 {
		n = num1
	} else {
		n = num2
	}

	switch {
	case n < K:
		nums = append(nums, num1)
		nums = append(nums, num2)
		units = " "

	case n < M:
		nums = append(nums, roundOffNearestTen(num1, K))
		nums = append(nums, roundOffNearestTen(num2, K))
		units = " per thousand "

	case n < G:
		nums = append(nums, roundOffNearestTen(num1, M))
		nums = append(nums, roundOffNearestTen(num2, M))
		units = " per million "

	case n < T:
		nums = append(nums, roundOffNearestTen(num1, G))
		nums = append(nums, roundOffNearestTen(num2, G))
		units = " per trillion "
	}

	return nums, units

}

func roundDownFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func trim(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundDownFloat(num*output)) / output
}

// GetInMB converts bytes to MB
func GetInMB(bytes uint64, precision int) float64 {
	temp := float64(bytes) / 1000000
	return trim(temp, precision)
}

// GetDateFromUnix gets a date and time in RFC822 format from a unix epoch
func GetDateFromUnix(createTime int64) string {
	t := time.Unix(createTime, 0)
	date := t.Format(time.RFC822)
	return date
}
