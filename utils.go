package libsysinfo

import (
	"strconv"
)

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err.Error())
	}

	return i
}

func atof64(s string) float64 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err.Error())
	}

	return f
}
