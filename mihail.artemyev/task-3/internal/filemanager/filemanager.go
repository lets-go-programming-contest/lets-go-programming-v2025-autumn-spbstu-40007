package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

const DirectoryPermissions = 0o755

func EnsureDirectoryExists(filePath string) error {
	parentDirectory := filepath.Dir(filePath)

	directoryCreationError := os.MkdirAll(parentDirectory, DirectoryPermissions)
	if directoryCreationError != nil {
		return fmt.Errorf(
			"creating directory %q: %w",
			parentDirectory,
			directoryCreationError,
		)
	}

	return nil
}
