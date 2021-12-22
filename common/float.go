package common

import "strconv"

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
