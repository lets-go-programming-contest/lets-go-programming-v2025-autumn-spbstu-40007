package jsonencoder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/treadwave/task-3/internal/sort"
	"github.com/treadwave/task-3/internal/structs"
)

func JSONEncoder(valutes []structs.Valute, outpath string) error {
	dir := filepath.Dir(outpath)

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("error creating dir: %w", err)
	}

	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error opening or creating file : %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
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
