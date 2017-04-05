package main

import (
	"github.com/BalkanTech/goilerplate/users"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, v *users.UserViews) {
	r.HandleFunc("/register", v.CreateHandler).Methods("POST")
	r.HandleFunc("/register", v.NewView.DefaultHandler)

	r.HandleFunc("/login", v.LoginHandler).Methods("POST")
	r.HandleFunc("/login", v.LoginView.DefaultHandler)

	r.HandleFunc("/logout", users.RequireLogin(v, v.LogoutHandler))

	r.HandleFunc("/change_password", users.RequireLogin(v, v.ChangePasswordHandler)).Methods("POST")
	r.HandleFunc("/change_password", users.RequireLogin(v, v.ChangePasswordView.DefaultHandler))

	r.HandleFunc("/account", users.RequireLogin(v, v.DisplayViewHandler))
	r.HandleFunc("/account/activate/{code}", v.ActivationHandler)
	r.HandleFunc("/account/edit", users.RequireLogin(v, v.UpdateHandler)).Methods("POST")
	r.HandleFunc("/account/edit", users.RequireLogin(v, v.EditViewHandler))
	r.HandleFunc("/account/delete", users.RequireLogin(v, v.DeleteHandler))
	r.HandleFunc("/account/undelete/{id}/{code}", users.RequireLogin(v, v.UndeleteHandler))


	//ToDo: Delete

	r.HandleFunc("/admin/user", users.RequireAdmin(v, v.AdminIndexHandler))
	r.HandleFunc("/admin/user/add", users.RequireAdmin(v, v.NewView.DefaultHandler))
	r.HandleFunc("/admin/user/view/{id}", users.RequireAdmin(v, v.DisplayViewHandler))
	r.HandleFunc("/admin/user/edit", users.RequireAdmin(v, v.UpdateHandler)).Methods("POST")
	r.HandleFunc("/admin/user/edit/{id}", users.RequireAdmin(v, v.EditViewHandler))
	r.HandleFunc("/admin/user/delete/{id}", users.RequireAdmin(v, v.DeleteHandler))
	r.HandleFunc("/admin/user/undelete/{id}", users.RequireAdmin(v, v.UndeleteHandler))
}