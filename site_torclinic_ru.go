package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func (t *ServerMed) ParserTorClinic(w http.ResponseWriter, r *http.Request, s Site) {
	err := t.TorClinic(w, r, s)
	if err != nil {
		t.returnError(w, r, err)
	} else {
		fmt.Fprint(w, t.StringToJson(map[string]string{"Ok": "the parser worked successfully"}))
	}
}

func (t *ServerMed) TorClinic(w http.ResponseWriter, r *http.Request, s Site) error {
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
	sliceGal := make([]ItemArr, 0)
	title := ""
	doc.Find("li.price_list--item").Each(func(i int, ss *goquery.Selection) {
		name1 := strings.TrimSpace(ss.Find("div:nth-child(2)").First().Text())
		name := strings.TrimSpace(ss.Find("div:nth-child(3)").First().Text())
		price := strings.TrimSpace(ss.Find("div:nth-child(1)").First().Text())
		if name1 == "" && price == "" {
			title = name
			return
		}
		if name1 != "" {
			name = fmt.Sprintf("%s | %s", name1, name)
		}

		if title != "" {
			name = fmt.Sprintf("%s | %s", title, name)
		} else {
			return
		}
		if name != "" && price != "" {
			sliceGal = append(sliceGal, ItemArr{Name: name, Price: price})
		}
	})
	if len(sliceGal) > 0 {
		err := t.WriteSliceToCsvNew(sliceGal, s)
		if err != nil {
			Logging(err)

		}
	}
	return nil
}
