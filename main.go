package main

import (
	"flag"
	"log"
	"github.com/BalkanTech/goilerplate/webserver"
	cfg "github.com/BalkanTech/goilerplate/config"
)

func main() {
	configFile := flag.String("config", "config.json", "Configuration file")
	flag.Parse()

	config := cfg.Config{File: *configFile}
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r, err := Init(&config)
	if err != nil {
		log.Fatal(err)
	}

	err = webserver.Start(&config, r)
	if err != nil {
		log.Fatal(err)
	}
}
