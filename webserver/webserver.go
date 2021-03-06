package webserver

import (
	"github.com/BalkanTech/goilerplate/config"
	"net/http"
	"log"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/csrf"
)

func Start(config *config.Config, router *mux.Router) error{
	secure := true
	if config.Debug {
		secure = false // In Debug mode csrf.Secure needs to be set to false
	}
	csrf := csrf.Protect([]byte(config.GetCSRF()), csrf.Secure(secure))

	server := &http.Server {
		Addr: config.Server.Addr(),
		ReadTimeout: config.Server.ReadTimeout * time.Second,
		WriteTimeout: config.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
		Handler: csrf(router),
	}

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("| Goilerplate project v0.0.1                                |")
	fmt.Println("| Copyright (c) 2017, Balkan C & T OOD                      |")
	fmt.Println("-------------------------------------------------------------")
	log.Printf("Starting webserver on \"%s\"", server.Addr)
	return server.ListenAndServe()
}