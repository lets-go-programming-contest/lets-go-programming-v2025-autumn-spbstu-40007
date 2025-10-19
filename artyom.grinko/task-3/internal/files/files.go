package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func wrapErr(err error) error {
	return fmt.Errorf("files: %w", err)
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func CreateIfNotExists(path string) error {
	dir := filepath.Dir(path)

	// Magic number? Are you serious?
	if err := os.MkdirAll(dir, 0o750); err != nil { //nolint:mnd
		return wrapErr(err)
	}

	if Exists(path) {
		if err := os.Remove(path); err != nil {
			return wrapErr(err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return wrapErr(err)
	}

	if err = file.Close(); err != nil {
		return wrapErr(err)
	}

	return nil
}
