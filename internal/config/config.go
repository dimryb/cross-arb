package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	MexcExchange = "mexc"
	JupExchange  = "jupiter"
)

type (
	CrossArbConfig struct {
		Log       Log                 `yaml:"log"`
		Exchanges map[string]Exchange `yaml:"exchanges"`
		Symbols   []string            `yaml:"symbols"`
		Scanner   ScannerConfig       `yaml:"scanner"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	// TODO: разделить структуры под разные биржи.

	Exchange struct {
		APIKey            string                `yaml:"apiKey" env:"API_KEY"`
		SecretKey         string                `yaml:"secretKey" env:"SECRET_KEY"`
		BaseURL           string                `yaml:"baseUrl"`
		BaseURLAdapter    string                `yaml:"baseUrlAdapter"`
		Timeout           time.Duration         `yaml:"timeout" env:"TIMEOUT"`
		Enabled           bool                  `yaml:"enabled"`
		OrderLimit        int                   `yaml:"orderLimit"`
		MaxPriceDiff      float64               `yaml:"maxPriceDiff"`
		MinQtyImprovement float64               `yaml:"minQtyImprovement"`
		Pairs             map[string]PairConfig `yaml:"pairs"`
	}

	PairConfig struct {
		Base  string `yaml:"base"`
		Quote string `yaml:"quote"`
	}

	ScannerConfig struct {
		Interval         string   `yaml:"interval"`
		Pairs            []string `yaml:"pairs"`
		LogOpportunities bool     `yaml:"logOpportunities"`
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
