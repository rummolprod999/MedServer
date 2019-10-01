package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func (t *ServerMed) ParserDelightLancetteRu(w http.ResponseWriter, r *http.Request, s Site) {
	err := t.DelightLancetteRu(w, r, s)
	if err != nil {
		t.returnError(w, r, err)
	} else {
		fmt.Fprint(w, t.StringToJson(map[string]string{"Ok": "the parser worked successfully"}))
	}
}

func (t *ServerMed) DelightLancetteRu(w http.ResponseWriter, r *http.Request, s Site) error {
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
	doc.Find("dd").Each(func(i int, ss *goquery.Selection) {
		nameParent := ss.Prev().Prev().Text()
		ss.Find("ul li").Each(func(i int, tt *goquery.Selection) {
			name := strings.TrimSpace(tt.Find("span").First().Text())
			price := strings.TrimSpace(tt.Find("span").Last().Text())
			if name == "Процедура" {
				return
			}
			name = fmt.Sprintf("%s | %s", nameParent, name)
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
	})
	return nil
}
