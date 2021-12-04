package common

import "path/filepath"

func RealPath(filePath string) (string, error) {
	sPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(sPath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
