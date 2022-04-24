package utils

import (
	"math"
	"os"
)

func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Point2Pixel(point float64) float64 {
	return point / 0.75
}

func RoundUp(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundUpInt(num*shift) / shift
}

func roundUpInt(num float64) float64 {
	t := math.Trunc(num)
	return t + math.Copysign(1, num)
}

func GetDirNames(path string, skipCondition func(entry os.DirEntry) bool) ([]string, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, dir := range dirs {
		if skipCondition(dir) {
			continue
		}
		result = append(result, dir.Name())
	}
	return result, nil
}
