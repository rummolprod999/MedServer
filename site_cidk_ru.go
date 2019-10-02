package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func (t *ServerMed) ParserCidkRu(w http.ResponseWriter, r *http.Request, s Site) {
	err := t.CidkRu(w, r, s)
	if err != nil {
		t.returnError(w, r, err)
	} else {
		fmt.Fprint(w, t.StringToJson(map[string]string{"Ok": "the parser worked successfully"}))
	}
}

func (t *ServerMed) CidkRu(w http.ResponseWriter, r *http.Request, s Site) error {
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
	doc.Find("item").Each(func(i int, ss *goquery.Selection) {
		gname := ss.Parent().AttrOr("title", "")
		ggname := ss.Parent().Parent().AttrOr("title", "")
		name := strings.TrimSpace(ss.Find("desc").First().AttrOr("value", ""))
		if gname != "" {
			name = fmt.Sprintf("%s | %s", gname, name)
		}
		if ggname != "" {
			name = fmt.Sprintf("%s | %s", ggname, name)
		}
		price := strings.TrimSpace(ss.Find("price").First().AttrOr("value", ""))
		if price == "Цена (руб.)" {
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
