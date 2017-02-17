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
	ChangePasswordView    *view.View
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
	a, _ := session.GetUser(r)

	u, err := v.manager.GetByID(a.ID)
	if err != nil {
		v.alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	v.EditView.Data = map [string]interface{}{"User": u}
	v.EditView.ExecuteTemplate(w, r)
}

func (v *UserViews) UpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *UserViews) DisplayViewHandler(w http.ResponseWriter, r *http.Request) {
	a, _ := session.GetUser(r)

	u, err := v.manager.GetByID(a.ID)
	if err != nil {
		v.alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	v.DisplayView.Data = map [string]interface{}{"User": u}
	v.DisplayView.ExecuteTemplate(w, r)
}

func (v *UserViews) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := session.GetUser(r)
	password := r.FormValue("Password")
	newpassword := r.FormValue("NewPassword")
	newpassword2 := r.FormValue("NewPassword2")

	if newpassword != newpassword2 {
		v.alerts.New("Error", "alert-danger", "New password and confirmation don't match")
		//ToDo: Set change password URL via config
		http.Redirect(w, r, "/change_password", http.StatusFound)
		return
	}

	user, err := v.manager.Authenticate(u.Username, password, AuthByUsername)
	if err != nil {
		v.alerts.New("Error", "alert-danger", "Invalid password")
		http.Redirect(w, r, "/change_password", http.StatusFound)
		return
	}

	user.SetPassword(newpassword)
	session.DestroySession(w)
	v.manager.Update(user)
	v.alerts.New("Success", "alert-success", "Your password has been updated. Please login again with your new password.")
	http.Redirect(w, r, "/login", http.StatusFound)
	return
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

func NewUserViews(manager UserManager, alerts *alerts.Alerts, templates string, new string, edit string, display string, login string, changepw string) *UserViews {
	views := &UserViews{manager: manager, alerts: alerts}
	views.NewView = view.NewView("Register", "base", alerts, templates+new)
	views.EditView = view.NewView("Edit", "base", alerts, templates+edit)
	views.LoginView = view.NewView("Login", "base", alerts, templates+login)
	views.ChangePasswordView = view.NewView("Change password", "base", alerts, templates+changepw)
	views.DisplayView = view.NewView("Account", "base", alerts, templates+display)
	// ToDo: Add edit + display

	return views
}
