package main

import (
	"net/http"
	"github.com/BalkanTech/goilerplate/users"
	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/sessions"
	"github.com/BalkanTech/goilerplate/view"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	u, err := UserController.Authenticate(login, password, users.AuthByUsernameOrEmail)
	if err != nil {
		view.Alerts = append(view.Alerts, alerts.Alert{Title: "Error", Class: "danger", Message: "Invalid login"})
		http.Redirect(w, r, "/login", 302)
		return
	}
	sessions.CreateSession(u, w)
	if u.ChangePassword {
		view.Alerts = append(view.Alerts, alerts.Alert{Title: "Warning", Class: "warning", Message: "You need to change your password."})
		http.Redirect(w, r, "/", 302)
		return
	}
	view.Alerts = append(view.Alerts, alerts.Alert{Title: "Success", Class: "success", Message: "You have succesfully logged in."})
	http.Redirect(w, r, "/", 302)
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessions.DestroySession(w)
	view.Alerts = append(view.Alerts, alerts.Alert{Title: "Success", Class: "success", Message: "You have succesfully logged out."})
	http.Redirect(w, r, "/", 302)
	return
}