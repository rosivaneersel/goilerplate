package session

import (
	"net/http"
	"github.com/gorilla/securecookie"
)

type ActiveUser struct {
	ID       string
	Username string
	EmailMD5 string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func CreateSession(id string, username string, w http.ResponseWriter) {
	v := map[string]interface{}{
		"id":       id,
		"username": username,
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
	return ActiveUser{ID: cValue["id"].(string), Username: cValue["username"].(string), EmailMD5: cValue["emailmd5"].(string)}, nil
}