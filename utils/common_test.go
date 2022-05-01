package utils

import "testing"

func TestIsExist(t *testing.T) {
	patterns := []struct {
		input    string
		expected bool
	}{
		{input: ".", expected: true},
		{input: "..", expected: true},
		{input: "../", expected: true},
		{input: "../utils", expected: true},
		{input: "../utils/common.go", expected: true},
		{input: "notExistsDir", expected: false},
		{input: "notExistsFile.txt", expected: false},
	}

	for i, pattern := range patterns {
		actual := isExist(pattern.input)
		if actual != pattern.expected {
			t.Errorf("patterns[%d] want to %v but %v", i, pattern.expected, actual)
		}
	}
}

func TestPoint2Pixel(t *testing.T) {
	patterns := []struct {
		input, expected float64
	}{
		{input: 0, expected: 0},
		{input: 10, expected: 13.33},
		{input: 20, expected: 26.67},
		{input: 50, expected: 66.67},
		{input: 75, expected: 100},
		{input: 100, expected: 133.33},
	}

	for i, pattern := range patterns {
		actual := point2Pixel(pattern.input)
		if actual != pattern.expected {
			t.Errorf("patterns[%d] want to %v but %v", i, pattern.expected, actual)
		}
	}
}

func TestRound(t *testing.T) {
	patterns := []struct {
		val, digits, expected float64
	}{
		{val: 4, digits: -1, expected: 0},
		{val: 5, digits: -1, expected: 10},
		{val: 6, digits: -1, expected: 10},
		{val: 0.4, digits: 0, expected: 0},
		{val: 0.5, digits: 0, expected: 1},
		{val: 0.6, digits: 0, expected: 1},
		{val: 0.04, digits: 1, expected: 0},
		{val: 0.05, digits: 1, expected: 0.1},
		{val: 0.06, digits: 1, expected: 0.1},
	}

	for i, pattern := range patterns {
		actual := round(pattern.val, pattern.digits)
		if actual != pattern.expected {
			t.Errorf("patterns[%d] want to %v but %v", i, pattern.expected, actual)
		}
	}
}

func TestRoundUp(t *testing.T) {
	patterns := []struct {
		val, digits, expected float64
	}{
		{val: 4, digits: -1, expected: 10},
		{val: 5, digits: -1, expected: 10},
		{val: 6, digits: -1, expected: 10},
		{val: 0.4, digits: 0, expected: 1},
		{val: 0.5, digits: 0, expected: 1},
		{val: 0.6, digits: 0, expected: 1},
		{val: 0.04, digits: 1, expected: 0.1},
		{val: 0.05, digits: 1, expected: 0.1},
		{val: 0.06, digits: 1, expected: 0.1},

		{val: 0.1, digits: 0, expected: 1},
		{val: 0, digits: 0, expected: 0},
		{val: 1, digits: 0, expected: 1},
	}

	for i, pattern := range patterns {
		actual := roundUp(pattern.val, pattern.digits)
		if actual != pattern.expected {
			t.Errorf("patterns[%d] want to %v but %v", i, pattern.expected, actual)
		}
	}
}

func TestRoundUpInt(t *testing.T) {
	patterns := []struct {
		val, expected float64
	}{
		{val: 0, expected: 0},
		{val: 0.1, expected: 1},
		{val: 1, expected: 1},
	}

	for i, pattern := range patterns {
		actual := roundUpInt(pattern.val)
		if actual != pattern.expected {
			t.Errorf("patterns[%d] want to %v but %v", i, pattern.expected, actual)
		}
	}
}
