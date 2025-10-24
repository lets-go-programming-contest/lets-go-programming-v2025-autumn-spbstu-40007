package processor

import (
 "encoding/json"
 "log"
 "os"
 "sort"

 "github.com/UwUshkin/task-3/internal/data" 
 "github.com/UwUshkin/task-3/internal/xmldecoder"
)

func ProcessAndSave(inputPath, outputPath string) error {
 log.Printf("Decoding XML from: %s", inputPath)
 valCurs, err := xmldecoder.DecodeCBRXML(inputPath)
 if err != nil {
  log.Printf("Error decoding XML: %v", err)
  return err
 }

 var outputCurrencies []data.OutputCurrency
 for _, valute := range valCurs.Valutes {
  outputCurrencies = append(outputCurrencies, valute.ConvertToOutput())
 }

 log.Println("Sorting currencies by Value (descending)")
 sort.Slice(outputCurrencies, func(i, j int) bool {
  return outputCurrencies[i].Value > outputCurrencies[j].Value
 })

 log.Printf("Encoding to JSON and saving to: %s", outputPath)
 jsonData, err := json.MarshalIndent(outputCurrencies, "", "  ") 
 if err != nil {
  log.Printf("Error marshaling to JSON: %v", err)
  return err
 }

 if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
  log.Printf("Error writing JSON file: %v", err)
  return err
 }

 log.Println("Processing complete. Data saved successfully.")
 return nil
}
