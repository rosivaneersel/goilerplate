package users

import (
	"net/http"
	"github.com/BalkanTech/goilerplate/session"
	"log"
)

var RequireLoginRedirectTo string = "/login"

func RequireLogin(v *UserViews, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetUser(r)
		if err != nil {
			v.Alerts.New("Warning", "alert-warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}
}

func RequireAdmin(v *UserViews, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := session.GetUser(r)
		if err != nil {
			v.Alerts.New("Warning", "alert-warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
			return
		}

		u, err := v.Manager.GetByID(user.ID)
		if !u.IsAdmin || err != nil {
			log.Printf("%s: Non-admin access attempt of: %s (%s)", r.URL.Path, u.Username, u.Email)
			v.Alerts.New("Warning", "alert-warning", "Insufficient privileges.")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		next(w, r)
	}
}
