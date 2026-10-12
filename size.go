// Package humane formats file sizes and relative dates the way macOS Finder does.
package humane

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// sizeUnits are Finder's capitalized, 1000-based unit labels above byte scale.
var sizeUnits = []string{"KB", "MB", "GB", "TB", "PB", "EB"}

// decimalPlaces lines up with sizeUnits: KB 0, MB 1, GB and above 2.
var decimalPlaces = []int{0, 1, 2, 2, 2, 2}

// HumanSize formats a byte count the way Finder does -- see docs/COMMENTS.md.
func HumanSize(bytes int64) string {
	if bytes == 0 {
		return "Zero KB"
	}
	if bytes < 1000 {
		if bytes == 1 {
			return "1 byte"
		}
		return fmt.Sprintf("%d bytes", bytes)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(1000))
	if exp > len(sizeUnits) {
		exp = len(sizeUnits)
	}
	value := roundForUnit(bytes, exp)
	if value >= 1000 && exp < len(sizeUnits) {
		exp++
		value = roundForUnit(bytes, exp)
	}

	return fmt.Sprintf("%s %s", trimZeros(value, decimalPlaces[exp-1]), sizeUnits[exp-1])
}

func roundForUnit(bytes int64, exp int) float64 {
	scale := math.Pow(10, float64(decimalPlaces[exp-1]))
	return math.Round(float64(bytes)/math.Pow(1000, float64(exp))*scale) / scale
}

func trimZeros(value float64, decimals int) string {
	s := strconv.FormatFloat(value, 'f', decimals, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
