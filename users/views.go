package users

import (
	"net/http"

	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/view"
	"github.com/BalkanTech/goilerplate/session"
	"github.com/gorilla/mux"
)

type UserViews struct {
	NewView     *view.View
	EditView    *view.View
	LoginView   *view.View
	DisplayView *view.View

	manager UserManager
	router  *mux.Router
	alerts  *alerts.Alerts
}

func (v *UserViews) NewViewHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *UserViews) CreateHandler(w http.ResponseWriter, r *http.Request) {
	//ToDo: Set endpoints via config
	username := r.FormValue("Username")
	password := r.FormValue("Password")
	password2 := r.FormValue("Password2")
	email := r.FormValue("Email")

	if password != password2 {
		v.alerts.New("Error", "alert-danger", "Passwords don't match")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	newUser := &User{Username: username, Email: email}
	newUser.SetPassword(password)

	err := v.manager.Create(newUser)
	if err != nil {
		v.alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	v.alerts.New("Success", "alert-info", "You have successfully registered your account. Please check your email to activate your account.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (v *UserViews) EditViewHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *UserViews) UpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *UserViews) LoginViewHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *UserViews) LoginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("Login")
	password := r.FormValue("Password")
	remember := r.FormValue("Remember")

	_ = remember // Ignore remember for now

	u, err := v.manager.Authenticate(login, password, AuthByUsernameOrEmail)
	if err != nil {
		v.alerts.New("Error", "alert-danger", "Invalid login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session.CreateSession(u.ID(), u.Username, w)
	if u.ChangePassword {
		v.alerts.New("Warning", "alert-warning", "You need to change your password.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	v.alerts.New("Success", "alert-success", "You have succesfully logged in.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (v *UserViews) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.DestroySession(w)
	v.alerts.New("Success", "alert-success", "You have succesfully logged out.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func NewUserViews(manager UserManager, alerts *alerts.Alerts, templates string, new string, edit string, display string, login string) *UserViews {
	views := &UserViews{manager: manager, alerts: alerts}
	views.NewView = view.NewView("Register", "base", alerts, templates+new)
	views.LoginView = view.NewView("Login", "base", alerts, templates+login)
	// ToDo: Add edit + display

	return views
}
