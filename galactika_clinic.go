package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func (t *ServerMed) ParserGalactikaClinic(w http.ResponseWriter, r *http.Request, s Site) {
	err := t.GalactikaClinic(w, r, s)
	if err != nil {
		t.returnError(w, r, err)
	} else {
		fmt.Fprint(w, t.StringToJson(map[string]string{"Ok": "the parser worked successfully"}))
	}
}

func (t *ServerMed) GalactikaClinic(w http.ResponseWriter, r *http.Request, s Site) error {
	defer SaveStack()
	downString := DownloadPage(s.Url)
	if downString != "" {
		Logging("received empty string")
		return errors.New("received empty string")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(downString))
	if err != nil {
		Logging(err)
		return err
	}
	doc.Find("#serv_spec603_center18 div.item").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())

	})
	return nil
}