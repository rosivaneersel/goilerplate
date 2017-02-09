package main

import (
	"github.com/BalkanTech/goilerplate/webserver"
	"flag"
	"log"
	"github.com/BalkanTech/goilerplate/view"
	db "github.com/BalkanTech/goilerplate/databases"
	cfg "github.com/BalkanTech/goilerplate/config"
	"github.com/BalkanTech/goilerplate/users"
	"github.com/jinzhu/gorm"
)

var DB = &gorm.DB{}
var UserController = &users.UserGorm{}

func main() {
	configFile := flag.String("config", "config.json", "Configuration file")
	flag.Parse()

	config := cfg.Config{File: *configFile}
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB, err = db.NewGormConnection(&config)
	if err != nil {
		log.Fatal(err)
	}

	router := webserver.Router

	rootView := view.NewView("templates/index.html")
	router.HandleFunc("/", rootView.DefaultHandler)

	UserController = users.NewUserManagerGorm(DB)

	loginView := view.NewView("templates/login.html")
	router.HandleFunc("/login", loginView.DefaultHandler).Methods("GET")


	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler)

	err = webserver.Start(*configFile)
	if err != nil {
		log.Fatal(err)
	}
}
