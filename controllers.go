package main

import (
	"html/template"
	"net/http"
)

func (t *ServerMed) indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API SERVER"
	tmpl, _ := template.New("data").Parse("<h1>{{.}}</h1>Examples:<p>GET /run/{galaktika.clinic} - run parser http://galaktika.clinic/prices/<p>GET /get/{galaktika.clinic} - return CSV file<p>")
	tmpl.Execute(w, data)
}

func (t *ServerMed) parserSite(w http.ResponseWriter, r *http.Request) {

}

func (t *ServerMed) returnCsv(w http.ResponseWriter, r *http.Request) {

}