package pathcreator

import (
	"os"
	"path/filepath"
)

func EnsureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}
