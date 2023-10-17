package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/snakoner/go-crypto-robot/internal/core"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/app.toml", "path to config file")
}

func main() {
	flag.Parse()

	config, err := core.NewConfig()
	if err != nil {
		return
	}

	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		return
	}

	core, err := core.New(config)
	if err != nil {
		return
	}

	if err := core.Start(); err != nil {
		return
	}
}
