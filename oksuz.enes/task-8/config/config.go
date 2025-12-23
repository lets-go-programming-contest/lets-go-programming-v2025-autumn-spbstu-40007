package config

type Config struct {
	Environment string
	LogLevel    string
}

func NewConfig() *Config {
	return &Config{
		Environment: "development",
		LogLevel:    "debug",
	}
}
