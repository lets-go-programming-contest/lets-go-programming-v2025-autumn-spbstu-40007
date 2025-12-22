package config

import "fmt"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var Current Config

func PrintInfo() {
	fmt.Printf("%s %s", Current.Environment, Current.LogLevel)
}
