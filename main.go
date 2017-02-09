package main

import (
	"flag"
	"log"

	cfg "github.com/BalkanTech/goilerplate/config"
	db "github.com/BalkanTech/goilerplate/databases"
	"github.com/BalkanTech/goilerplate/users"
	"github.com/BalkanTech/goilerplate/view"
	"github.com/BalkanTech/goilerplate/webserver"
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

	rootView := view.NewView("Goilerplate", "base", "templates/root.html")
	router.HandleFunc("/", rootView.DefaultHandler)

	UserController = users.NewUserManagerGorm(DB)

	loginView := view.NewView("Login", "base", "templates/login.html")
	router.HandleFunc("/login", loginView.DefaultHandler).Methods("GET")

	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler)

	registerView := view. NewView("Register", "base", "templates/register.html")
	router.HandleFunc("/register", registerView.DefaultHandler)

	err = webserver.Start(*configFile)
	if err != nil {
		log.Fatal(err)
	}
}
