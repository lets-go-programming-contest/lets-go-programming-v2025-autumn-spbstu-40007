package jsonencoder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treadwave/task-3/internal/sort"
	"github.com/treadwave/task-3/internal/structs"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func JSONEncoder(valutes []structs.Valute, outpath string) error {
	dir := filepath.Dir(outpath)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("error creating dir: %w", err)
	}

	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("error opening or creating file : %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("warning: error closing file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)

	sortedValutes := sort.Sort(valutes)

	err = encoder.Encode(sortedValutes)
	if err != nil {
		return fmt.Errorf("error encoding xml: %w", err)
	}

	return nil
}
