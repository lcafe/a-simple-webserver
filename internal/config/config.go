package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AdminToken     string            `json:"admin_token"`
	HTTPPort       string            `json:"http_port"`
	HTTPSPort      string            `json:"https_port"`
	UseHTTPS       bool              `json:"use_https"`
	CertFile       string            `json:"cert_file"`
	KeyFile        string            `json:"key_file"`
	VirtualHosts   map[string]string `json:"virtual_hosts"`
	FileServer     string            `json:"file_server"`
	BandwidthLimit int               `json:"bandwidth_limit"`
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
