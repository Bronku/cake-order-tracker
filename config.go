package main

import (
	_ "embed"
	"time"

	"github.com/BurntSushi/toml"
)

//go:embed config.default.toml
var defaultConfig string

var config struct {
	Database struct {
		File string
	}

	Server struct {
		Addr              string
		ReadHeaderTimeout time.Duration
		RequestTimeout    time.Duration
	}
}

func loadConfig() {
	_, err := toml.Decode(defaultConfig, &config)
	if err != nil {
		panic(err)
	}
	_, _ = toml.DecodeFile("config.toml", &config)
}
