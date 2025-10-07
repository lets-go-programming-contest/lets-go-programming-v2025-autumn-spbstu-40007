package file

import (
	"os"
	"path/filepath"
)

func CreateIfNotExists(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil { //nolint:noinlineerr
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
