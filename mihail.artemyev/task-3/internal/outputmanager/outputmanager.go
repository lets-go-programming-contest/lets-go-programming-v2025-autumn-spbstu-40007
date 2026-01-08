package outputmanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mart22052006/task-3/internal/currencydata"
	"github.com/Mart22052006/task-3/internal/filemanager"
)

const DefaultOutputFilePermissions = 0o600

func WriteJSONOutput(
	outputFilePath string,
	currencyExchangeList []currencydata.CurrencyExchange,
) error {
	directoryCreationError := filemanager.EnsureDirectoryExists(outputFilePath)
	if directoryCreationError != nil {
		return fmt.Errorf("creating output directory: %w", directoryCreationError)
	}

	jsonEncodedData, encodingError := json.MarshalIndent(currencyExchangeList, "", "  ")
	if encodingError != nil {
		return fmt.Errorf("marshalling data to JSON format: %w", encodingError)
	}

	fileWritingError := os.WriteFile(outputFilePath, jsonEncodedData, DefaultOutputFilePermissions)
	if fileWritingError != nil {
		return fmt.Errorf("writing JSON to file %q: %w", outputFilePath, fileWritingError)
	}

	return nil
}
