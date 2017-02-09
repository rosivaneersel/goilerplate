package view

import (
	"github.com/BalkanTech/goilerplate/alerts"
	"html/template"
	"net/http"
	"log"
)


type View struct {
	template *template.Template
	Data interface{}
}

type responseData struct {
	Data interface{}
	Alerts []alerts.Alert
}

var Alerts = []alerts.Alert{}

func (v *View) Execute(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "html")
	w.WriteHeader(http.StatusOK)

	r := responseData{v.Data, Alerts}
	Alerts = []alerts.Alert{}

	err := v.template.Execute(w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *View) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	v.Execute(w)
}

func NewView(filename string) *View{
	t := template.Must(template.ParseFiles(filename))
	v := &View{template: t}
	return v
}