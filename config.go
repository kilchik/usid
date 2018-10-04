package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type Config struct {
	Listen             string `toml:"listen"`
	GoogleClientId     string `toml:"google_client_id"`
	GoogleClientSecret string `toml:"google_client_secret"`
}

func (c *Config) validate() error {
	if c.Listen == "" {
		return fmt.Errorf("no listen in config")
	}
	if c.GoogleClientId == "" {
		return fmt.Errorf("no google_client_id in config")
	}
	if c.GoogleClientSecret == "" {
		return fmt.Errorf("no google_client_secret")
	}
	return nil
}

func readConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		logE.Fatalf("open %q: %v", path, err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		logE.Fatalf("read config: %v", err)
	}
	conf := &Config{}
	if _, err := toml.Decode(string(content), conf); err != nil {
		logE.Fatalf("decode config: %v", err)
	}
	return conf, conf.validate()
}
