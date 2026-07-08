// Package humane formats file sizes and relative dates the way macOS
// Finder does, modeled on Swift's ByteCountFormatter and
// RelativeDateTimeFormatter: small, configurable formatter types with a
// single Format method, rather than a bare function.
package humane

import (
	"fmt"
	"math"
)

// sizeUnits are Finder's capitalized, 1000-based unit labels -- not the
// SI-correct lowercase "kB", and not 1024-based "KiB"/"MiB". No published
// Go library ships this exact combination; see the README for why.
var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// SizeFormatter formats byte counts the way Finder does: 1000-based
// math, capitalized unit labels, rounded to 2 significant digits. The
// zero value is ready to use -- there's no configuration yet, since
// Finder's is the only style this package currently needs.
type SizeFormatter struct{}

// Format returns bytes as a Finder-style human-readable string.
//
//	SizeFormatter{}.Format(79992)  == "80 KB"
//	SizeFormatter{}.Format(225935) == "226 KB"
//	SizeFormatter{}.Format(1500000) == "1.5 MB"
func (f SizeFormatter) Format(bytes int64) string {
	if bytes < 1000 {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(1000))
	if exp >= len(sizeUnits) {
		exp = len(sizeUnits) - 1
	}
	rounded := math.Round(float64(bytes)/math.Pow(1000, float64(exp))*10) / 10
	if rounded < 10 {
		return fmt.Sprintf("%.1f %s", rounded, sizeUnits[exp])
	}
	return fmt.Sprintf("%.0f %s", rounded, sizeUnits[exp])
}
