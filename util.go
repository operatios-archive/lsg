package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func isPathHidden(path string) bool {
	if path == "." {
		return false
	}

	components := splitPath(path)

	for i := 1; i <= len(components); i++ {
		f := filepath.Join(components[:i]...)

		if strings.Contains(f, "..") {
			abs, _ := filepath.Abs(f)
			f = abs
		}

		file, err := newFile(f)
		if err != nil {
			continue
		}

		if file.isHidden() {
			return true
		}
	}
	return false
}

func splitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func humanizeSize(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10)
	}

	fSize := float64(size)
	fSize /= 1024

	for _, unit := range []rune{'K', 'M', 'G', 'T', 'P', 'E', 'Z'} {
		if fSize < 9 {
			return fmt.Sprintf("%.1f%c", fSize, unit)
		} else if fSize < 1000 {
			return fmt.Sprintf("%.0f%c", fSize, unit)
		}
		fSize /= 1024
	}
	return fmt.Sprintf("%.1f%c", fSize, 'Y')
}
