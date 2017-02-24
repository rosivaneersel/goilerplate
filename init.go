package main

import (
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/users"
	cfg "github.com/BalkanTech/goilerplate/config"
	db "github.com/BalkanTech/goilerplate/databases"
	"github.com/BalkanTech/goilerplate/view"
)

func initWithGorm(db *gorm.DB, r *mux.Router, a *alerts.Alerts, t string) {

	RootView := view.NewView("Goilerplate", "base", a, "templates/root.html")
	r.HandleFunc("/", RootView.DefaultHandler)

	UserManager := users.NewUserManagerGorm(db)
	UserManager.Init(false) // Initialize DB on first use and make migrations (if any)
	// ToDo: Set Template directory via config and make NewUserView use this setting
	UserViews := users.NewUserViews(UserManager, a, t, "user_register.html", "user_edit.html", "user_show.html", "login.html", "change_password.html", "user_index_admin.html")
	UserRoutes(r, UserViews)
}

func initWithMGO(db *db.MongoDBConnection, r *mux.Router, a *alerts.Alerts, t string) {
}

func Init(c *cfg.Config) (r *mux.Router, e error) {
	if !c.Database.IsValidType() {
		return nil, fmt.Errorf("Invalid database type in config file")
	}

	r = mux.NewRouter()
	a := &alerts.Alerts{}

	if c.Database.IsGorm() {
		DB, err := db.NewGormConnection(c)
		if err != nil {
			return nil, err
		}

		initWithGorm(DB, r, a, c.GetTemplatesPath())
	}

	if c.Database.IsMGO() {
		DB, err := db.NewMgoConnection(c)
		if err != nil {
			return nil, err
		}

		initWithMGO(DB, r, a, c.GetTemplatesPath())
	}

	r.PathPrefix(c.Static.GetURL()).Handler(http.StripPrefix(c.Static.GetURL(), http.FileServer(http.Dir(c.Static.GetPath()))))
	return
}