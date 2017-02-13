package view

import (
	"html/template"
	"log"
	"net/http"

	"github.com/BalkanTech/goilerplate/alerts"
)

type View struct {
	Title    string
	Data     interface{}
	alerts  *alerts.Alerts
	template *template.Template
	layout   string
}

type responseData struct {
	Title  string
	Data   interface{}
	Alerts []alerts.Alert
}

var DefaultFiles = []string{"templates/index.html", "templates/_nav.html"}
//var Alerts = []alerts.Alert{}

func (v *View) ExecuteTemplate(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "html")
	w.WriteHeader(http.StatusOK)

	r := responseData{v.Title, v.Data, v.alerts.Get()}

	err := v.template.ExecuteTemplate(w, v.layout, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *View) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	v.ExecuteTemplate(w)
}

func NewView(title string, layout string, alerts *alerts.Alerts, files ...string) *View {
	fs := append(DefaultFiles, files...)
	t := template.Must(template.ParseFiles(fs...))
	v := &View{template: t, layout: layout, Title: title, alerts: alerts}
	return v
}
