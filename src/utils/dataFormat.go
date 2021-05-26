/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	Q = math.Pow(10, 15)
)

func roundOffNearestTen(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

func RoundValues(num1, num2 float64, inBytes bool) ([]float64, string) {
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
		units = " per billion "

	case n < Q:
		nums = append(nums, roundOffNearestTen(num1, T))
		nums = append(nums, roundOffNearestTen(num2, T))
		units = " per trillion "

	case n >= Q:
		nums = append(nums, roundOffNearestTen(num1, Q))
		nums = append(nums, roundOffNearestTen(num2, Q))
		units = " per quadrillion "
	}

	if inBytes {
		switch units {
		case " ":
			units = " B "
		case " per thousand ":
			units = " kB "
		case " per million ":
			units = " mB "
		case " per billion ":
			units = " gB "
		case " per trillion ":
			units = " tB "
		case " per quadrillion ":
			units = " pB "
		}
	}

	return nums, units

}

func roundDownFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Trim trims a float to the specified number of precision decimal digits.
func Trim(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundDownFloat(num*output)) / output
}

// GetInMB converts bytes to MB
func GetInMB(bytes uint64, precision int) float64 {
	temp := float64(bytes) / 1000000
	return Trim(temp, precision)
}

// GetDateFromUnix gets a date and time in RFC822 format from a unix epoch
func GetDateFromUnix(createTime int64) string {
	t := time.Unix(createTime/1000, 0)
	date := t.Format(time.RFC822)
	return date
}

func RoundFloat(num float64, base string, precision int) float64 {
	x := num
	div := math.Pow10(precision)
	switch base {
	case "K":
		x /= 1024
	case "M":
		x /= (1024 * 1024)
	case "G":
		x /= (1024 * 1024 * 1024)
	}
	return math.Round(x*div) / div
}

func RoundUint(num uint64, base string, precision int) float64 {
	x := float64(num)
	return RoundFloat(x, base, precision)
}
