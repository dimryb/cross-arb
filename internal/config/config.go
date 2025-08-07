package config

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	CrossArbConfig struct {
		Log       Log                 `yaml:"log"`
		Exchanges map[string]Exchange `yaml:"exchanges"`
		Symbols   []string            `yaml:"symbols"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	Exchange struct {
		APIKey            string  `yaml:"apiKey" env:"API_KEY"`
		SecretKey         string  `yaml:"secretKey" env:"SECRET_KEY"`
		BaseURL           string  `yaml:"baseUrl"`
		Enabled           bool    `yaml:"enabled"`
		OrderLimit        int     `yaml:"orderLimit"`
		MaxPriceDiff      float64 `yaml:"maxPriceDiff"`
		MinQtyImprovement float64 `yaml:"minQtyImprovement"`
	}
)

func Load(configPath string, target any) error {
	err := cleanenv.ReadConfig(path.Join("./", configPath), target)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(target)
	if err != nil {
		return fmt.Errorf("error updating env: %w", err)
	}

	return nil
}

func NewConfig(configPath string) (*CrossArbConfig, error) {
	cfg := &CrossArbConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}
