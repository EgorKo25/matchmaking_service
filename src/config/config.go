package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const configPath = "./config/config.yaml"

func NewAPIGWConfig() (*APIGWConfig, error) {
	cfg := &APIGWConfig{}
	if err := cfg.loadConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

type APIGWConfig struct {
	*ServerConfig     `yaml:"server"`
	*MatchmakerConfig `yaml:"matchmaker"`
}

type ServerConfig struct {
	TlsCert string `yaml:"tls-cert"`
	KeyFile string `yaml:"key-file"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type MatchmakerConfig struct {
	GroupSize             int           `yaml:"group-size"`
	AcceptableWaitingTime time.Duration `yaml:"acceptable-waiting-time"`
	DeltaLatency          float64       `yaml:"delta-latency"`
	DeltaSkill            float64       `yaml:"delta-skill"`
}

func (a *APIGWConfig) loadConfig() error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&a); err != nil {
		return err
	}
	return err
}
