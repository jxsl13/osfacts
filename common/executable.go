package common

import "os"

func IsExecutable(filePath string) (bool, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	fm := fi.Mode().Perm()

	return fm&0111 == 0111, nil
}
