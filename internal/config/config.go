package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HTTPPort    string `json:"http_port"`
	VirtualHost string `json:"virtual_host"`
	FileServer  string `json:"file_server"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
