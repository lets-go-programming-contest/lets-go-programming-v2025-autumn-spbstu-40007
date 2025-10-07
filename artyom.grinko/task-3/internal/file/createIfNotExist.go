package file

import (
	"os"
	"path/filepath"
)

func CreateIfNotExists(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if Exists(path) {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}
