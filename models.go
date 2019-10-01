package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Site struct {
	Alias    string
	Url      string
	FileName string
}

func (t *ServerMed) findSiteInList(alias string) (error, Site) {
	for _, el := range sites {
		if el.Alias == alias {
			return nil, el
		}
	}
	return errors.New("site not found"), Site{}

}

func (t *ServerMed) StringToJson(st map[string]string) string {
	b, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		return err.Error()
	} else {
		return string(b[:])
	}
}

func (t *ServerMed) Parser(w http.ResponseWriter, r *http.Request, s Site) {
	switch {
	case s.Alias == "galaktika.clinic":
		t.ParserGalactikaClinic(w, r, s)
	case s.Alias == "cidk.ru":
		t.ParserCidkRu(w, r, s)
	case s.Alias == "delight-lancette.ru":
		t.ParserDelightLancetteRu(w, r, s)
	default:
		t.returnError(w, r, errors.New("site not found"))
	}
}

func (t *ServerMed) GetCsv(w http.ResponseWriter, r *http.Request, s Site) {
	switch {
	case s.Alias == "galaktika.clinic":
		t.ReturnFileCsvToClient(w, r, s)
	case s.Alias == "cidk.ru":
		t.ReturnFileCsvToClient(w, r, s)
	case s.Alias == "delight-lancette.ru":
		t.ReturnFileCsvToClient(w, r, s)
	default:
		t.returnError(w, r, errors.New("site not found"))
	}
}

func (t *ServerMed) WriteToCsv(mp map[string]string, s Site) error {
	//currentTime := time.Now()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	dirf := filepath.FromSlash(fmt.Sprintf("%s/%s/%s", dir, DirTemp, s.FileName))
	w, err := os.Create(dirf)
	if err != nil {
		return err
	}
	defer w.Close()
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	defer writer.Flush()
	for k, v := range mp {
		err := writer.Write([]string{k, v})
		if err != nil {
			return err
		}
	}
	return nil
}
