package xmldecoder

import (
 "encoding/xml"
 "os"

 "github.com/UwUshkin/task-3/internal/data" 
)

func DecodeCBRXML(filePath string) (*data.ValCurs, error) {
 xmlFile, err := os.Open(filePath)
 if err != nil {
  return nil, err
 }
 defer xmlFile.Close()

 decoder := xml.NewDecoder(xmlFile)
 var valCurs data.ValCurs

 if err := decoder.Decode(&valCurs); err != nil {
  return nil, err
 }

 return &valCurs, nil
}
