package config

func LoadProdConfig(cfg *Config) {
	cfg.Environment = "production"
	cfg.LogLevel = "info"
}
