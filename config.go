package main

import (
	"os"

	log "github.com/junwangustc/ustclog"
	"github.com/naoina/toml"
)

type Config struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Ip       string `toml:"ip"`
	Port     uint   `toml:"port"`
	Dbname   string `toml:"dbname"`
}

func NewConfig() *Config {
	c := &Config{}
	return c
}
func ParseConfig(path string) (cfg *Config, err error) {
	if path == "" {
		log.Fatalln("no configuration provided, using default settings")
	}
	log.Printf("Using configuration at: %s\n", path)
	config := NewConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	return config, toml.NewDecoder(f).Decode(&config)
}
