package main

import (
	"io"
	"os"
)

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

func fileTruncated(file *os.File) (bool, error) {
	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}

	return info.Size() < offset, nil
}
