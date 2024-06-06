package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/vebcreatex7/diploma_magister/internal/config"
	"os"
)

func initConfig(path string) (config.API, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return config.API{}, fmt.Errorf("openning file: %w", err)
	}

	var cfg config.API

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return config.API{}, fmt.Errorf("unmarshaling config: %w", err)
	}

	return cfg, nil
}
