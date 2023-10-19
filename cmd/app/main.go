package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigChan, syscall.SIGINT)

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

	go func() {
		<-sigChan
		cancel()
	}()

	if err := core.Start(ctx); err != nil {
		return
	}
	fmt.Println("Finish")
}
