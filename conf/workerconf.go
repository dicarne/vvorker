package conf

import (
	"encoding/json"
	"errors"
)

type ExtensionConfig struct {
	Binding string `json:"binding"`
	Name    string `json:"name"`
}

type AiConfig struct {
	Model   string `json:"model"`
	ApiKey  string `json:"api_key"`
	BaseUrl string `json:"base_url"`
	Binding string `json:"binding"`
}

type SQLDBConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Binding    string `json:"binding"`
	ResourceID string `json:"resource_id"`
}

type OSSConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Binding         string `json:"binding"`
	Bucket          string `json:"bucket"`
	UseSSL          bool   `json:"use_ssl"`
	Region          string `json:"region"`
	ResourceID      string `json:"resource_id"`
	SessionToken    string `json:"session_token"`
}

type KV struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Binding    string `json:"binding"`
	ResourceID string `json:"resource_id"`
}

type Assets struct {
	Binding   string `json:"binding"`
	Directory string `json:"directory"`
}

type WorkerConfig struct {
	ProjectName        string            `json:"name"`
	Version            string            `json:"version"`
	Extensions         []ExtensionConfig `json:"extensions"`
	Services           []string          `json:"services"`
	CompatibilityFlags []string          `json:"compatibility_flags"`
	Vars               json.RawMessage   `json:"vars"`
	Ai                 []AiConfig        `json:"ai"`
	PgSql              []SQLDBConfig     `json:"pgsql"`
	OSS                []OSSConfig       `json:"oss"`
	KV                 []KV              `json:"kv"`
	Assets             []Assets          `json:"assets"`
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
		Extensions:         []ExtensionConfig{},
		Services:           []string{},
		CompatibilityFlags: []string{},
	}
}
