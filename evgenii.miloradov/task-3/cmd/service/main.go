

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/evgenii.miloradov/task-3/internal/config"
	"github.com/evgenii.miloradov/task-3/internal/data"
)

func main() {
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configFile == "" {
		fmt.Fprintln(os.Stderr, "Configuration file path is required")
		os.Exit(1)
	}

	appConfig, configErr := config.Load(*configFile)
	if configErr != nil {
		fmt.Fprintln(os.Stderr, configErr)
		os.Exit(1)
	}

	currencyData, dataErr := data.LoadFromFile(appConfig.InputFile)
	if dataErr != nil {
		fmt.Fprintln(os.Stderr, "Failed to load currency data:", dataErr)
		os.Exit(1)
	}

	sort.Slice(currencyData.Currencies, func(i, j int) bool {
		return currencyData.Currencies[i].Value > currencyData.Currencies[j].Value
	})

	if saveErr := currencyData.ExportToFile(appConfig.OutputFile); saveErr != nil {
		fmt.Fprintln(os.Stderr, "Failed to save output file:", saveErr)
		os.Exit(1)
	}
}
