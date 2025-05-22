package conf

import (
	"encoding/json"
	"errors"
)

type WorkerConfig struct {
	ProjectName        string          `json:"name"`
	Version            string          `json:"version"`
	Extensions         []string        `json:"extensions"`
	Services           []string        `json:"services"`
	CompatibilityFlags []string        `json:"compatibility_flags"`
	Vars               json.RawMessage `json:"vars"`
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
		ProjectName:        "default",
		Version:            "0.0.1",
		Extensions:         []string{},
		Services:           []string{},
		CompatibilityFlags: []string{},
	}
}
