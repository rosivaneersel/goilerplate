package main

import (
	"github.com/BalkanTech/goilerplate/users"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, v *users.UserViews) {
	r.HandleFunc("/register", v.CreateHandler).Methods("POST")
	r.HandleFunc("/register", v.NewView.DefaultHandler)

	// ToDo: display
	// ToDo: delete
	r.HandleFunc("/login", v.LoginHandler).Methods("POST")
	r.HandleFunc("/login", v.LoginView.DefaultHandler)

	r.HandleFunc("/logout", v.LogoutHandler)

	r.HandleFunc("/change_password", v.ChangePasswordHandler).Methods("POST")
	r.HandleFunc("/change_password", v.ChangePasswordView.DefaultHandler)

	r.HandleFunc("/account", v.DisplayViewHandler)

	r.HandleFunc("/account/edit", v.UpdateHandler).Methods("POST")
	r.HandleFunc("/account/edit", v.EditViewHandler)
}
