package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
)

func (t *ServerMed) indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API SERVER"
	tmpl, _ := template.New("data").Parse("<h1>{{.}}</h1>Examples:<p>GET /run/{galaktika.clinic} - run parser http://galaktika.clinic/prices/<p>GET /get/{galaktika.clinic} - return CSV file<p>")
	tmpl.Execute(w, data)
}

func (t *ServerMed) parserSite(w http.ResponseWriter, r *http.Request) {
	defer SaveStack()
	vars := mux.Vars(r)
	siteParam := strings.TrimSpace(vars["site"])
	err, st := t.findSiteInList(siteParam)
	if err != nil {
		t.returnError(w, r, err)
	}
	t.Parser(w, r, st)
}

func (t *ServerMed) returnCsv(w http.ResponseWriter, r *http.Request) {
	defer SaveStack()
	vars := mux.Vars(r)
	siteParam := strings.TrimSpace(vars["site"])
	err, st := t.findSiteInList(siteParam)
	if err != nil {
		t.returnError(w, r, err)
	}
	t.GetCsv(w, r, st)
}

func (t *ServerMed) returnError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	Logging(err)
	fmt.Fprint(w, t.StringToJson(map[string]string{"Error": err.Error()}))
}
