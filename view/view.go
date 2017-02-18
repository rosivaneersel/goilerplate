package view

import (
	"html/template"
	"log"
	"net/http"
	"github.com/BalkanTech/goilerplate/alerts"
	"github.com/BalkanTech/goilerplate/session"
	"github.com/gorilla/csrf"
)


type View struct {
	Title    string
	Data     map[string]interface{}
	alerts  *alerts.Alerts
	template *template.Template
	layout   string
}

type responseData struct {
	Title   string
	Data    map[string]interface{}
	Alerts  []alerts.Alert
	Session session.ActiveUser
	csrfField  template.HTML
}

var DefaultFiles = []string{"templates/index.html", "templates/_nav.html"}

func (v *View) ExecuteTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "html")
	w.WriteHeader(http.StatusOK)

	u, err := session.GetUser(r)
	res := responseData{v.Title, v.Data, v.alerts.Get(), u, csrf.TemplateField(r)}

	err = v.template.ExecuteTemplate(w, v.layout, res)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *View) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	v.ExecuteTemplate(w, r)
}

func NewView(title string, layout string, alerts *alerts.Alerts, files ...string) *View {
	fs := append(DefaultFiles, files...)
	t := template.Must(template.ParseFiles(fs...))
	v := &View{template: t, layout: layout, Title: title, alerts: alerts}
	return v
}
