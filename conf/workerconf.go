package conf

import (
	"encoding/json"
	"errors"
)

type WorkerConfig struct {
	ProjectName     string   `json:"name"`
	Version         string   `json:"version"`
	EnabledServices []string `json:"enabled_services"`
}

func ParseWorkerConfig(s string) (*WorkerConfig, error) {
	if s == "" {
		return nil, errors.New("input string is empty")
	}
	var config WorkerConfig
	err := json.Unmarshal([]byte(s), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func DefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		ProjectName:     "default",
		Version:         "0.0.1",
		EnabledServices: []string{"ai"},
	}
}
