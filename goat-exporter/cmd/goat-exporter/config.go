package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"port"`
	RPCNode string `mapstructure:"goat_rpc_node"`
}

func initConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("port", "8080")
	_ = viper.BindEnv("goat_rpc_node")
}

func loadConfig() (Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if strings.TrimSpace(cfg.RPCNode) == "" {
		return fmt.Errorf("goat rpc node cannot be empty; set GOAT_RPC_NODE, e.g. https://rpc.goat.network")
	}
	if _, err := url.ParseRequestURI(cfg.RPCNode); err != nil {
		return fmt.Errorf("goat rpc node invalid url: %w", err)
	}
	return nil
}
