package main

import "os"

func fileRotated(file *os.File, path string) (bool, error) {
	current, err := file.Stat()
	if err != nil {
		return false, err
	}

	latest, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return !os.SameFile(current, latest), nil
}
