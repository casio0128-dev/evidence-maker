package utils

import (
	"image"
	"math"
	"os"
)

func isExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func point2Pixel(point float64) float64 {
	return round(point/0.75, 2)
}

func round(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundInt(num*shift) / shift
}

func roundInt(num float64) float64 {
	t := math.Trunc(num)
	if math.Abs(num-t) >= 0.5 {
		return t + math.Copysign(1, num)
	}
	return t
}

func roundUp(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundUpInt(num*shift) / shift
}

func roundUpInt(num float64) float64 {
	t := math.Trunc(num)
	return t + math.Copysign(1, num)
}

func getDirNames(path string, skipCondition func(entry os.DirEntry) bool) ([]string, error) {
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

func getImageSize(path string) (height, width int, err error) {
	pict, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}

	img, _, err := image.Decode(pict)
	if err != nil {
		return 0, 0, err
	}

	bound := img.Bounds()
	height = bound.Max.Y
	width = bound.Max.X

	return height, width, nil
}
