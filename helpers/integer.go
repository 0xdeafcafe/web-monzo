package helpers

import (
	"fmt"
	"math"
	"strings"
)

// ToIntegerSegment ..
func ToIntegerSegment(i int64, incDecimal bool) string {
	decimalStr := "."
	if !incDecimal {
		decimalStr = ""
	}

	precise := float64(i) / 100
	val := math.Abs(precise)
	return fmt.Sprintf("%d%s", int64(val), decimalStr)
}

// ToFractionalSegment ..
func ToFractionalSegment(i int64, incDecimal bool) string {
	decimalStr := "."
	if !incDecimal {
		decimalStr = ""
	}

	precise := float64(i) / 100
	trunc := precise - math.Trunc(precise)
	val := strings.Split(fmt.Sprintf("%.2f", trunc), ".")[1]
	return fmt.Sprintf("%s%s", decimalStr, val)
}
