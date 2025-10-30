module github.com/Smolyaninoff/GoLang

go 1.22.7

require (
	github.com/goccy/go-yaml v1.18.0
	golang.org/x/text v0.14.0
)

replace github.com/Smolyaninoff/GoLang/internal/config => ./internal/config

replace github.com/Smolyaninoff/GoLang/internal/currencies => ./internal/currencies
