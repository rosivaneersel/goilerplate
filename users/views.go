package users

import (
	"net/http"

	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/view"
	"github.com/BalkanTech/goilerplate/session"
	"github.com/gorilla/mux"
	"io/ioutil"

	"path"
)

type UserViews struct {
	NewView     *view.View
	EditView    *view.View
	ChangePasswordView    *view.View
	LoginView   *view.View
	DisplayView *view.View

	Manager UserManager
	Router  *mux.Router
	Alerts  *alerts.Alerts
}

func (v *UserViews) CreateHandler(w http.ResponseWriter, r *http.Request) {
	//ToDo: Set endpoints via config
	username := r.FormValue("Username")
	password := r.FormValue("Password")
	password2 := r.FormValue("Password2")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")

	if password != password2 {
		v.Alerts.New("Error", "alert-danger", "Passwords don't match")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	newUser := &User{Username: username, Email: email}
	newUser.Profile.FirstName = firstName
	newUser.Profile.LastName = lastName
	newUser.SetPassword(password)

	err := v.Manager.Create(newUser)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	v.Alerts.New("Success", "alert-info", "You have successfully registered your account. Please check your email to activate your account.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (v *UserViews) EditViewHandler(w http.ResponseWriter, r *http.Request) {
	a, _ := session.GetUser(r)

	u, err := v.Manager.GetByID(a.ID)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	v.EditView.Data = map [string]interface{}{"User": u}
	v.EditView.ExecuteTemplate(w, r)
}

func (v *UserViews) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("UserID")
	username := r.FormValue("Username")
	email := r.FormValue("Email")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")

	u, err := v.Manager.GetByID(userID)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/account/edit", http.StatusFound)
		return
	}

	u.Email = email
	u.Username = username
	u.Profile.FirstName = firstName
	u.Profile.LastName = lastName

	// Update avatar
	// Get file data from form
	file, header, err := r.FormFile("AvatarFile")
	if err != http.ErrMissingFile {
		if err != nil && err != http.ErrMissingFile{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Read the file
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//ToDo: Get static path/URL from config
		//ToDo: Crop image to 500x500 pixelsgo
		// Store the file
		filename := path.Join("static/avatars", userID+path.Ext(header.Filename))
		err = ioutil.WriteFile(filename, data, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update Profile.AvatarURL
		u.Profile.AvatarURL = "static/avatars/" + userID+path.Ext(header.Filename)
	}

	err = v.Manager.Update(u)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/account/edit", http.StatusFound)
		return
	}
	v.Alerts.New("Success", "alert-success", "Your profile has been updated successfully")
	http.Redirect(w, r, "/account", http.StatusFound)
	return

}

func (v *UserViews) DisplayViewHandler(w http.ResponseWriter, r *http.Request) {
	a, _ := session.GetUser(r)

	u, err := v.Manager.GetByID(a.ID)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
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
		v.Alerts.New("Error", "alert-danger", "New password and confirmation don't match")
		//ToDo: Set change password URL via config
		http.Redirect(w, r, "/change_password", http.StatusFound)
		return
	}

	user, err := v.Manager.Authenticate(u.Username, password, AuthByUsername)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", "Invalid password")
		http.Redirect(w, r, "/change_password", http.StatusFound)
		return
	}

	user.SetPassword(newpassword)
	session.DestroySession(w)
	v.Manager.Update(user)
	v.Alerts.New("Success", "alert-success", "Your password has been updated. Please login again with your new password.")
	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func (v *UserViews) LoginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("Login")
	password := r.FormValue("Password")
	remember := r.FormValue("Remember")

	_ = remember // Ignore remember for now

	u, err := v.Manager.Authenticate(login, password, AuthByUsernameOrEmail)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", "Invalid login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session.CreateSession(u.ID(), u.Username, w)
	if u.ChangePassword {
		v.Alerts.New("Warning", "alert-warning", "You need to change your password.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	v.Alerts.New("Success", "alert-success", "You have succesfully logged in.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (v *UserViews) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.DestroySession(w)
	v.Alerts.New("Success", "alert-success", "You have succesfully logged out.")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func NewUserViews(manager UserManager, alerts *alerts.Alerts, templates string, new string, edit string, display string, login string, changepw string) *UserViews {
	views := &UserViews{Manager: manager, Alerts: alerts}
	views.NewView = view.NewView("Register", "base", alerts, templates+new)
	views.EditView = view.NewView("Edit", "base", alerts, templates+edit)
	views.LoginView = view.NewView("Login", "base", alerts, templates+login)
	views.ChangePasswordView = view.NewView("Change password", "base", alerts, templates+changepw)
	views.DisplayView = view.NewView("Account", "base", alerts, templates+display)

	return views
}
