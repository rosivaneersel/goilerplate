package webserver

import (
	"github.com/BalkanTech/goilerplate/config"
	"net/http"
	"log"
	"time"
	"fmt"
	"github.com/gorilla/mux"
)

func Start(config *config.Config, router *mux.Router) error{
	server := &http.Server {
		Addr: config.Server.Addr(),
		ReadTimeout: config.Server.ReadTimeout * time.Second,
		WriteTimeout: config.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
		Handler: router,
	}

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("| Goilerplate project v0.0.1                                |")
	fmt.Println("| Copyright (c) 2017, Balkan C & T OOD                      |")
	fmt.Println("-------------------------------------------------------------")
	log.Printf("Starting webserver on \"%s\"", server.Addr)
	return server.ListenAndServe()
}