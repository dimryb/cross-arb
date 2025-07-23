package config

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	CrossArbConfig struct {
		Log Log `yaml:"log"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
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
