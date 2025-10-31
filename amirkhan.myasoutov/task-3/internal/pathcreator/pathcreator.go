package pathcreator

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirPermissions = 0o755

func EnsureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, dirPermissions)
	if err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	return nil
}
