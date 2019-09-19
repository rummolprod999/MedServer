package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"path/filepath"
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

func (t *ServerMed) GetGalactikaClinic(w http.ResponseWriter, r *http.Request, s Site) {
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
	http.ServeFile(w, r, dirf)
}

func (t *ServerMed) GalactikaClinic(w http.ResponseWriter, r *http.Request, s Site) error {
	defer SaveStack()
	downString := DownloadPage(s.Url)
	if downString == "" {
		Logging("received empty string")
		return errors.New("received empty string")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(downString))
	if err != nil {
		Logging(err)
		return err
	}
	mapGal := make(map[string]string)
	doc.Find("#serv_spec603_center18 div.item").Each(func(i int, ss *goquery.Selection) {
		name := strings.TrimSpace(ss.Find("div.name").First().Text())
		//name = fmt.Sprintf("'%s'", name)
		price := strings.TrimSpace(ss.Find("div.price").First().Text())
		//price = fmt.Sprintf("'%s'", price)
		if name != "" && price != "" {
			mapGal[name] = price
		}
		if len(mapGal) > 0 {
			err := t.WriteToCsv(mapGal, s)
			if err != nil {
				Logging(err)

			}
		}

	})
	return nil
}
