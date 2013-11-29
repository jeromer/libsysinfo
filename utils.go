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
