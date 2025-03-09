package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	HTTPPort          string            `json:"http_port"`
	ProxyPrefix       string            `json:"proxy_prefix"`
	MaxVirtualHosts   int               `json:"max_virtual_hosts"`
	VirtualHosts      map[string]string `json:"virtual_hosts"`
	FileServerRootUrl string            `json:"file_server_root_url"`
	FileServer        string            `json:"file_server"`
}

var (
	instance *Config
	once     sync.Once
)

func openConfigFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo de configuração: %w", err)
	}
	return f, nil
}

func parseConfig(f *os.File) (*Config, error) {
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao decodificar o arquivo de configuração: %w", err)
	}
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.HTTPPort == "" {
		return fmt.Errorf("campo 'http_port' não informado")
	}

	if cfg.FileServer == "" {
		return fmt.Errorf("campo 'file_server' não informado")
	}

	if cfg.MaxVirtualHosts > 0 && len(cfg.VirtualHosts) > cfg.MaxVirtualHosts {
		return fmt.Errorf("número de virtual_hosts (%d) excede o limite máximo (%d)",
			len(cfg.VirtualHosts), cfg.MaxVirtualHosts)
	}

	return nil
}

func loadConfig(filename string) (*Config, error) {
	f, err := openConfigFile(filename)
	if err != nil {
		return nil, err
	}
	return parseConfig(f)
}

func GetConfig(filename string) (*Config, error) {
	var err error
	once.Do(func() {
		instance, err = loadConfig(filename)
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
