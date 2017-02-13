package users

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/BalkanTech/goilerplate/alerts"
	"log"
)

type ActiveUser struct {
	ID       string
	Username string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func CreateSession(u *User, w http.ResponseWriter) {
	v := map[string]interface{}{
		"id":       u.ID(),
		"username": u.Username,
	}

	if encoded, err := cookieHandler.Encode("session", v); err == nil {
		c := &http.Cookie{
			Name:  "goilerplate",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
	}
}

func DestroySession(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:   "goilerplate",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}

func GetUser(r *http.Request) (user ActiveUser, err error) {
	c, err := r.Cookie("goilerplate")

	if err != nil {
		return user, err
	}

	cValue := make(map[string]interface{})
	if err = cookieHandler.Decode("session", c.Value, &cValue); err != nil {
		return user, err
	}
	return ActiveUser{ID: cValue["id"].(string), Username: cValue["username"].(string)}, nil
}

var RequireLoginRedirectTo string = "/login"

func RequireLogin(a *alerts.Alerts, users UserManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUser(r)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}

		_, err = users.GetByID(user.ID)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}
		next(w, r)
	}
}

func RequireAdmin(a *alerts.Alerts, users UserManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUser(r)
		if err != nil {
			a.New("Warning", "warning", "You need to login to access that page.")
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
			return
		}

		u, err := users.GetByID(user.ID)
		if !u.IsAdmin || err != nil {
			log.Printf("%s: Non-admin access attempt", r.URL.Path)
			http.Error(w, "Access denied", http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
