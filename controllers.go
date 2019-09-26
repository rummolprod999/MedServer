package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (t *ServerMed) indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API SERVER"
	tmpl, _ := template.New("data").Parse("<h1>{{.}}</h1>Examples:<p>GET /run/galaktika.clinic - run parser http://galaktika.clinic/prices/<p>GET /get/galaktika.clinic - return CSV file<p>")
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

func (t *ServerMed) ReturnFileCsvToClient(w http.ResponseWriter, r *http.Request, s Site) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.returnError(w, r, err)
		return
	}
	dirf := filepath.FromSlash(fmt.Sprintf("%s/%s/%s", dir, DirTemp, s.FileName))
	if _, err := os.Stat(dirf); os.IsNotExist(err) {
		t.returnError(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+s.FileName)
	w.Header().Set("Content-Type", "application/CSV")
	http.ServeFile(w, r, dirf)
}
