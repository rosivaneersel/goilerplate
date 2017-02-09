package view

import (
	"html/template"
	"log"
	"net/http"

	"github.com/BalkanTech/goilerplate/alerts"
)

type View struct {
	Title    string
	template *template.Template
	layout   string
	Data     interface{}
}

type responseData struct {
	Title  string
	Data   interface{}
	Alerts []alerts.Alert
}

var Alerts = []alerts.Alert{}

func (v *View) Execute(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "html")
	w.WriteHeader(http.StatusOK)

	r := responseData{v.Title, v.Data, Alerts}
	Alerts = []alerts.Alert{}

	err := v.template.ExecuteTemplate(w, v.layout, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *View) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	v.Execute(w)
}

func NewView(title string, layout string, files ...string) *View {
	t := template.Must(template.ParseFiles(files...))
	v := &View{template: t}
	return v
}
