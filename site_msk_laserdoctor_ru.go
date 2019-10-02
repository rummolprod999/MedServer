package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func (t *ServerMed) ParserMskLaserDoctor(w http.ResponseWriter, r *http.Request, s Site) {
	err := t.MskLaserDoctor(w, r, s)
	if err != nil {
		t.returnError(w, r, err)
	} else {
		fmt.Fprint(w, t.StringToJson(map[string]string{"Ok": "the parser worked successfully"}))
	}
}

func (t *ServerMed) MskLaserDoctor(w http.ResponseWriter, r *http.Request, s Site) error {
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
	doc.Find("div.price_list_item_header").EachWithBreak(func(i int, ss *goquery.Selection) bool {
		url := ss.AttrOr("data-href", "")
		url = fmt.Sprintf("https://msk.laserdoctor.ru%s", url)
		name := strings.TrimSpace(ss.Text())
		ret := t.MskLaserDoctorInner(url, name, &sliceGal, s)
		if !ret {
			return false
		}
		return true
	})
	if len(sliceGal) > 0 {
		err := t.WriteSliceToCsvNew(sliceGal, s)
		if err != nil {
			Logging(err)

		}
	}
	return nil
}

func (t *ServerMed) MskLaserDoctorInner(url, name string, sliceGal *[]ItemArr, s Site) bool {
	downString := DownloadPage(url)
	if downString == "" {
		Logging("received empty string")
		return false
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(downString))
	if err != nil {
		Logging(err)
		return false
	}
	doc.Find("div.price_list_item_body_list_item_header").EachWithBreak(func(i int, ss *goquery.Selection) bool {
		url := ss.AttrOr("data-href", "")
		url = fmt.Sprintf("https://msk.laserdoctor.ru%s", url)
		name1 := strings.TrimSpace(ss.Text())
		t.MskLaserDoctorInner1(url, name, name1, sliceGal, s)
		return true
	})
	return true
}

func (t *ServerMed) MskLaserDoctorInner1(url, name, name1 string, sliceGal *[]ItemArr, s Site) bool {
	downString := DownloadPage(url)
	if downString == "" {
		Logging("received empty string")
		return false
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(downString))
	if err != nil {
		Logging(err)
		return false
	}
	doc.Find("div.price_table_row").Each(func(i int, ss *goquery.Selection) {
		name2 := strings.TrimSpace(ss.Find("div.name").First().Text())
		if name2 == "" {
			return
		}
		fullName := cleanString(fmt.Sprintf("%s | %s | %s", name, name1, name2))
		price := strings.TrimSpace(strings.Replace(ss.Find("div:contains('Стоимость:')").First().Text(), "Стоимость:", "", -1))
		if fullName != "" && price != "" {
			*sliceGal = append(*sliceGal, ItemArr{Name: fullName, Price: price})
		}
	})
	return true
}
