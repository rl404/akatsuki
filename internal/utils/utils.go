package utils

import (
	"strconv"
	"strings"
)

// SplitDate will split string format to year,month,day.
// 2020-01-02 => 2020,1,2
// 2020-01 => 2020,1,0
// 2020 => 2020,0,0
func SplitDate(date string) (year int, month int, day int, err error) {
	if date == "" {
		return
	}

	split := strings.Split(date, "-")

	if len(split) >= 1 {
		year, err = strconv.Atoi(split[0])
		if err != nil {
			return
		}
	}

	if len(split) >= 2 {
		month, err = strconv.Atoi(split[1])
		if err != nil {
			return
		}
	}

	if len(split) >= 3 {
		day, err = strconv.Atoi(split[2])
		if err != nil {
			return
		}
	}

	return
}
