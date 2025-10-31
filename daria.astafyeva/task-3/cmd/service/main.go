package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"

	"github.com/itsdasha/task-3/internal/config"
	"github.com/itsdasha/task-3/internal/data"
	"github.com/itsdasha/task-3/internal/output"
)

type CurrencyProcessor struct {
	Raw    []data.Currency
	Result []data.OutputCurrency
}

func main() {
	var cfgPath, fmtType string
	flag.StringVar(&cfgPath, "config", "config.yaml", "Path to YAML config")
	flag.StringVar(&fmtType, "output-format", "json", "Output format: json, yaml, xml")
	flag.Parse()

	settings := config.LoadSettings(cfgPath)

	processor := new(CurrencyProcessor)
	processor.Raw = data.LoadCurrencies(settings.SourceFile)

	sort.Slice(processor.Raw, func(i, j int) bool {
		return processor.Raw[i].Rate > processor.Raw[j].Rate
	})

	processor.Convert()
	output.Save(processor.Result, settings.TargetFile, fmtType)

	fmt.Printf("Processed %d currencies -> '%s' [%s]\n",
		len(processor.Result), settings.TargetFile, fmtType)
}

func (p *CurrencyProcessor) Convert() {
	for _, currency := range p.Raw {
		var num int
		if currency.NumCode != "" {
			val, err := strconv.Atoi(currency.NumCode)
			if err != nil {
				panic(fmt.Errorf("invalid NumCode '%s': %w", currency.NumCode, err))
			}
			num = val
		}

		p.Result = append(p.Result, data.OutputCurrency{
			Num:    num,
			Char:   currency.CharCode,
			Amount: currency.Rate,
		})
	}
}
