package files

import (
    "os"
    "path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func CreateIfNotExists(path string) error {
	dir := filepath.Dir(path)

	// Magic number? Are you serious?
	if err := os.MkdirAll(dir, 0o750); err != nil { //nolint:noinlineerr,mnd
		return err //nolint:wrapcheck
	}

	if Exists(path) {
		if err := os.Remove(path); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err = file.Close(); err != nil { //nolint:noinlineerr
		return err //nolint:wrapcheck
	}

	return nil
}

