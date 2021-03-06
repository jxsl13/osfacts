package common

import (
	"errors"
	"sort"
	"strconv"
)

var (
	ErrMapEmpty = errors.New("map empty")
)

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

func MapContainsKeys(m map[string]string, keys ...string) bool {
	for _, key := range keys {
		_, found := m[key]
		if !found {
			return false
		}
	}
	return true
}

func IntMapContainsKeys(m map[string]int, keys ...string) bool {
	for _, key := range keys {
		_, found := m[key]
		if !found {
			return false
		}
	}
	return true
}

func IntMapToStringValues(m map[string]int) []string {
	result := make([]string, 0, len(m))
	for _, v := range m {
		result = append(result, strconv.Itoa(v))
	}
	return result
}

func IntMapToIntValues(m map[string]int) []int {
	result := make([]int, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(result)))
	return result
}

func FirstIntItem(m map[string]int) (int, error) {
	for _, v := range m {
		return v, nil
	}

	return 0, ErrMapEmpty
}

func FirstIntItemWithDefault(m map[string]int, def int) int {
	for _, v := range m {
		return v
	}
	return def
}

func FirstStringItem(m map[string]string) (string, error) {
	for _, v := range m {
		return v, nil
	}

	return "", ErrMapEmpty
}

func FirstStringItemWithDefault(m map[string]string, def string) string {
	for _, v := range m {
		return v
	}

	return def
}
