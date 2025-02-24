package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Cfg struct {
	DB     DbConfig
	Server ServerCfg
	File   FileCfg
}

type DbConfig struct {
	Connection string
}

type ServerCfg struct {
	Port string
}

type FileCfg struct {
	Path string
}

func GetConfig(path string) Cfg {
	var cfg Cfg
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	return cfg
}
