package common

import "strconv"

func IntToStringMap(m map[string]int) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = strconv.Itoa(v)
	}

	return result
}

func Int64ToStringMap(m map[string]int64) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = strconv.FormatInt(v, 10)
	}

	return result
}

func UpdateMap(m, n map[string]string) map[string]string {
	result := make(map[string]string, len(m))

	for k, v := range m {
		result[k] = v
	}

	for k, v := range n {
		result[k] = v
	}
	return result
}
