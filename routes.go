package main

import (
	"github.com/BalkanTech/goilerplate/users"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, v *users.UserViews) {
	r.HandleFunc("/register", v.NewView.DefaultHandler)
	r.HandleFunc("/register", v.CreateHandler).Methods("POST")
	// ToDo: edit
	// ToDo: display
	// ToDo: delete
	r.HandleFunc("/login", v.LoginView.DefaultHandler)
	r.HandleFunc("/login", v.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", v.LogoutHandler)
}
