package common

import (
	"bufio"
	"errors"
	"os"
)

var (
	ErrFileEmpty = errors.New("the provided file is empty")
)

func GetFileLines(filePath string) ([]string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return SplitLines(string(b)), nil
}

func GetFileFirstLine(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", ErrFileEmpty
}
