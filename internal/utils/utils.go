package utils

import (
	"strconv"
	"strings"
	"time"
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

// ParseToTimePtr to parse str to time pointer.
func ParseToTimePtr(layout, str string) *time.Time {
	tmp, err := time.Parse(layout, str)
	if err != nil {
		return nil
	}
	return &tmp
}

// ParseToBoolPtr to parse str to bool pointer.
func ParseToBoolPtr(str string) *bool {
	if str == "" {
		return nil
	}
	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil
	}
	return &b
}
