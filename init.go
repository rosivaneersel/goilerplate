package main

import (
	"github.com/gorilla/mux"
	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/users"
	db "github.com/BalkanTech/goilerplate/databases"
	cfg "github.com/BalkanTech/goilerplate/config"
	"github.com/BalkanTech/goilerplate/view"
	"log"
	"flag"
	"net/http"
	"os"
)

func init() {
	a := &alerts.Alerts{}
	var UserManager users.UserManager

	configFile := flag.String("config", "config.json", "Configuration file")
	flag.Parse()

	config = &cfg.Config{File: *configFile}
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if !config.Database.IsValidType() {
		log.Fatalf("Invalid database type in config file")
	}

	router = mux.NewRouter()
	RootView := view.NewView("Goilerplate", "base", a, "templates/root.html")
	router.HandleFunc("/", RootView.DefaultHandler)

	if _, err := os.Stat(config.Static.GetPath()); os.IsNotExist(err) {
		log.Fatalf("Static path \"%s\"", config.Static.GetPath())
	}
	router.PathPrefix(config.Static.GetURL()).Handler(http.StripPrefix(config.Static.GetURL(), http.FileServer(http.Dir(config.Static.GetPath()))))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	if config.Database.IsGorm() {
		DB, err := db.NewGormConnection(config)
		if err != nil {
			log.Fatalf(err.Error())
		}

		UserManager = users.NewUserManagerGorm(DB)
		UserManager.Init(false) // Initialize DB on first use and make migrations (if any)
	}

	if config.Database.IsMGO() {
		DB, err := db.NewMgoConnection(config)
		if err != nil {
			log.Fatal(err)
		}

		UserManager = users.NewUserManagerMgo(DB, "users")
		UserManager.Init(false) // Initialize DB on first use and make migrations (if any)
	}

	// Todo: Get template dir from config
	initUserTemplates("templates/", UserManager, a)
}


func initUserTemplates(t string, m users.UserManager, a *alerts.Alerts) {
	userTemplates := &users.Templates{
		BaseTemplates: t,
		New: "user_register.html",
		Edit: "user_edit.html",
		Display: "user_show.html",
		Login: "login.html",
		ChangePassword:"change_password.html",
		AdminIndex: "user_index_admin.html",
	}
	UserViews := users.Views(m, a, userTemplates)
	UserRoutes(router, UserViews)
}