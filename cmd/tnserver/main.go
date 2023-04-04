package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/artemxgod/transaction-server/internal/app/tnserver"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/tnserver.toml", "path to the cfg file")
}

func main() {
	flag.Parse()

	cfg := tnserver.NewConfig()
	toml.DecodeFile(configPath, cfg)

	if err := tnserver.Start(cfg); err != nil {
		log.Fatal(err)
	}
}

