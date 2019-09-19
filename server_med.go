package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ServerMed struct {
	Port string
}

func (t *ServerMed) run() {
	router := mux.NewRouter()
	router.HandleFunc(`/run/{site:}`, t.parserSite)
	router.HandleFunc(`/get/{site:}`, t.returnCsv)
	router.HandleFunc("/", t.indexHandler)
	http.Handle("/", router)
	if err := http.ListenAndServe(t.Port, nil); err != nil {
		Logging(err)
		return
	}

}
