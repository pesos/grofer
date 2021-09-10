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
	"fmt"
	"math"
	"time"
)

var (
	kilo = math.Pow(10, 3)
	mega = math.Pow(10, 6)
	giga = math.Pow(10, 9)
	tera = math.Pow(10, 12)
	peta = math.Pow(10, 15)
)

// SecondsToHuman converts a given UNIX epoch in seconds to human readble format
func SecondsToHuman(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = fmt.Sprintf("%dy %dm %dw %dd %dH %dM %dS",
			int(years),
			int(months),
			int(weeks),
			int(days),
			int(hours),
			int(minutes),
			int(seconds),
		)
	} else if months > 0 {
		result = fmt.Sprintf("%dm %dw %dd %dH %dM %dS",
			int(months),
			int(weeks),
			int(days),
			int(hours),
			int(minutes),
			int(seconds),
		)
	} else if weeks > 0 {
		result = fmt.Sprintf("%dw %dd %dH %dM %dS",
			int(weeks),
			int(days),
			int(hours),
			int(minutes),
			int(seconds),
		)
	} else if days > 0 {
		result = fmt.Sprintf("%dd %dH %dM %dS",
			int(days),
			int(hours),
			int(minutes),
			int(seconds),
		)
	} else if hours > 0 {
		result = fmt.Sprintf("%dH %dM %dS",
			int(hours),
			int(minutes),
			int(seconds),
		)
	} else if minutes > 0 {
		result = fmt.Sprintf("%dM %dS",
			int(minutes),
			int(seconds),
		)
	} else {
		result = fmt.Sprintf("%dS",
			int(seconds),
		)
	}

	return
}

func roundOffNearestTen(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

// RoundValues rounds off values to nearest K, G, M, etc. Returns the rounded values and the unit. If inBytes is set to true, units are returned as B, kB, gB, etc.
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
	case n < kilo:
		nums = append(nums, num1)
		nums = append(nums, num2)
		units = " "

	case n < mega:
		nums = append(nums, roundOffNearestTen(num1, kilo))
		nums = append(nums, roundOffNearestTen(num2, kilo))
		units = " per thousand "

	case n < giga:
		nums = append(nums, roundOffNearestTen(num1, mega))
		nums = append(nums, roundOffNearestTen(num2, mega))
		units = " per million "

	case n < tera:
		nums = append(nums, roundOffNearestTen(num1, giga))
		nums = append(nums, roundOffNearestTen(num2, giga))
		units = " per billion "

	case n < peta:
		nums = append(nums, roundOffNearestTen(num1, tera))
		nums = append(nums, roundOffNearestTen(num2, tera))
		units = " per trillion "

	case n >= peta:
		nums = append(nums, roundOffNearestTen(num1, peta))
		nums = append(nums, roundOffNearestTen(num2, peta))
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

// GetDateFromUnix gets a date and time in RFC822 format from a unix epoch in millisecond
func GetDateFromUnix(createTime int64) string {
	t := time.Unix(createTime/1000, 0)
	date := t.Format(time.RFC822)
	return date
}

// RoundFloat rounds off a float to a given base and precision
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

// RoundUint rounds a float by making use of RoundFloat
func RoundUint(num uint64, base string, precision int) float64 {
	x := float64(num)
	return RoundFloat(x, base, precision)
}
