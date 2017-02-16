package main

import (
	"github.com/BalkanTech/goilerplate/users"
	"github.com/gorilla/mux"
	"github.com/BalkanTech/goilerplate/view"
	"github.com/BalkanTech/goilerplate/alerts"
)

func UserRoutes(r *mux.Router, v *users.UserViews) {
	r.HandleFunc("/register", v.CreateHandler).Methods("POST")
	r.HandleFunc("/register", v.NewView.DefaultHandler)

	// ToDo: edit
	// ToDo: display
	// ToDo: delete
	r.HandleFunc("/login", v.LoginHandler).Methods("POST")
	r.HandleFunc("/login", v.LoginView.DefaultHandler)

	r.HandleFunc("/logout", v.LogoutHandler)

	a := &alerts.Alerts{}

	profileShowView := view.NewView("Profile", "base", a, "templates/profile_show.html")
	r.HandleFunc("/profile", profileShowView.DefaultHandler)
}
