package users

import (
	"net/http"
	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/session"
)

var RequireLoginRedirectTo string = "/login"

func RequireLogin(a *alerts.Alerts, users UserManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetUser(r)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}

		/*_, err = users.GetByID(user.ID)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}*/
		next(w, r)
	}
}

func RequireAdmin(a *alerts.Alerts, users UserManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetUser(r)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
			return
		}

		/*u, err := users.GetByID(user.ID)
		if !u.IsAdmin || err != nil {
			log.Printf("%s: Non-admin access attempt", r.URL.Path)
			http.Error(w, "Access denied", http.StatusBadRequest)
			return
		}*/
		next(w, r)
	}
}
