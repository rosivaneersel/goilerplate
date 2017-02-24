package main

import (
	"log"
	"github.com/BalkanTech/goilerplate/webserver"
	cfg "github.com/BalkanTech/goilerplate/config"
	"github.com/gorilla/mux"
)

var config *cfg.Config
var router *mux.Router

func main() {
	err := webserver.Start(config, router)
	if err != nil {
		log.Fatal(err)
	}
}