package users

import (
	"net/http"

	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/view"
	"github.com/BalkanTech/goilerplate/session"
	"github.com/gorilla/mux"
	"io/ioutil"
	"path"
	"reflect"
)

type UserViews struct {
	NewView     *view.View
	EditView    *view.View
	ChangePasswordView    *view.View
	LoginView   *view.View
	DisplayView *view.View

	AdminIndexView *view.View
	AdminNewView *view.View
	AdminEditView *view.View
	AdminDisplayView *view.View

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
	u := &User{}
	a, _ := session.GetUser(r)

	au, err := v.Manager.GetByID(a.ID)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if au.IsAdmin {
		vars := mux.Vars(r)
		id := vars["id"]
		u, err = v.Manager.GetByID(id)
		if err != nil {
			v.Alerts.New("Error", "alert-danger", err.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		u = au
	}

	v.EditView.Data = map [string]interface{}{"User": u}
	v.EditView.ExecuteTemplate(w, r)
}

func (v *UserViews) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	au, _ := session.GetUser(r)

	userID := r.FormValue("UserID")
	username := r.FormValue("Username")
	email := r.FormValue("Email")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")

	u, err := v.Manager.GetByID(userID)
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		if au.IsAdmin {
			http.Redirect(w, r, "/admin/user/edit/" + userID, http.StatusFound)
		} else {
			http.Redirect(w, r, "/account/edit", http.StatusFound)
		}

		return
	}

	u.Email = email
	u.Username = username
	u.Profile.FirstName = firstName
	u.Profile.LastName = lastName
	if au.IsAdmin {
		u.IsAdmin = r.FormValue("IsAdmin") != ""
		u.ChangePassword = r.FormValue("ChangePassword") != ""
	}

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
	if(au.IsAdmin) {
		http.Redirect(w, r, "/admin/user", http.StatusFound)
	} else {
		http.Redirect(w, r, "/account", http.StatusFound)
	}
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

	session.CreateSession(u.ID(), u.Username, u.IsAdmin, w)
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

func (v *UserViews) AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	users, err := v.Manager.Find("", "")
	if err != nil {
		v.Alerts.New("Error", "alert-danger", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	v.AdminIndexView.Data = map [string]interface{}{"Users": users}
	v.AdminIndexView.ExecuteTemplate(w, r)
}

func (v *UserViews) AdminShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	u, err := v.Manager.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v.AdminDisplayView.Data = map [string]interface{}{"User": u}
	v.AdminDisplayView.ExecuteTemplate(w, r)
}

type Templates struct {
	BaseTemplates string
	New string
	Edit string
	Display string
	Login string
	ChangePassword string

	AdminIndex string
	AdminEdit string
}

func (t *Templates) Get(e string) string {
	v := reflect.ValueOf(t).Elem().FieldByName(e)
	if !v.IsValid() {
		return ""
	}
	return t.BaseTemplates + v.String()
}

func Views(manager UserManager, alerts *alerts.Alerts, t *Templates) *UserViews {
	views := &UserViews{Manager: manager, Alerts: alerts}
	views.NewView = view.NewView("Register", "base", alerts, t.Get("New"))
	views.EditView = view.NewView("Edit", "base", alerts, t.Get("Edit"))
	views.LoginView = view.NewView("Login", "base", alerts, t.Get("Login"))
	views.ChangePasswordView = view.NewView("Change password", "base", alerts, t.Get("ChangePassword"))
	views.DisplayView = view.NewView("Account", "base", alerts, t.Get("Display"))

	views.AdminIndexView = view.NewView("Admin - Users", "base", alerts, t.Get("AdminIndex"))

	return views
}

//ToDo: Admin views and handlers
//ToDo: Configuration views
// Todo: Account activation
// Todo: Forgotten password