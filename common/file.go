package common

import "os"

func GetFileLines(filePath string) ([]string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return SplitLines(string(b)), nil
}
