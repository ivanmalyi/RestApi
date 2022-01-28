package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/ivanmalyi/RestApi/internal/app/appserver"
	"log"
)

var configPath string

func init()  {
	flag.StringVar(&configPath, "config-path", "config/appserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := appserver.NewConfig()
	_,err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	server := appserver.New(config)
	err = server.Start();
	if err!=nil {
		log.Fatal(err)
	}
}
