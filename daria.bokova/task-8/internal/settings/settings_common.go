package settings

import (
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env     string `yaml:"env"`
	Logging string `yaml:"logging"`
}

func ParseConfig(data []byte) (*AppConfig, error) {
	var cfg AppConfig
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Объявляем переменную-функцию
var GetConfig func() []byte
