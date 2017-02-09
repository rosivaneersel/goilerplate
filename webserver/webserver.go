package webserver

import (
	"github.com/BalkanTech/goilerplate/config"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"fmt"
)

var Router = mux.NewRouter()

func Start(configFile string) error{
	c := &config.Config{File: configFile}
	c.Load()

	server := &http.Server {
		Addr: c.Server.Addr(),
		ReadTimeout: c.Server.ReadTimeout * time.Second,
		WriteTimeout: c.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: c.Server.MaxHeaderBytes,
		Handler: Router,
	}

	//ToDo: Set static directory via config
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("| Goilerplate project v0.0.1                                |")
	fmt.Println("| Copyright (c) 2017, Balkan C & T OOD                      |")
	fmt.Println("-------------------------------------------------------------")
	log.Printf("Starting webserver on \"%s\"", server.Addr)
	return server.ListenAndServe()
}