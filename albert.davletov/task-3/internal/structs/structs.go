package structs

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type Valutes struct {
	Valutes []Valute `xml:"Valute"`
}

type TempValute struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
}

type TempValutes struct {
	TempValutes []TempValute `xml:"Valute"`
}
