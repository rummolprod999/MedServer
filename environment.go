package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Filelog string

var DirLog = "log_server_med"
var DirTemp = "temp_server_med"
var FileLog Filelog
var Port string
var SecretKey = "11111111" //change it
var mutex sync.Mutex
var sites []Site

func GetPort() {
	flag.Parse()
	ArgS := flag.Arg(0)
	if ArgS == "" {
		fmt.Println("server run on default port 8181")
		Port = "8181"
	} else {
		Port = ArgS
	}
}

func Logging(args ...interface{}) {
	mutex.Lock()
	file, err := os.OpenFile(string(FileLog), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("error writing to log file", err)
		return
	}
	fmt.Fprintf(file, "%v  ", time.Now())
	for _, v := range args {

		fmt.Fprintf(file, " %v", v)
	}
	//fmt.Fprintf(file, " %s", UrlXml)
	fmt.Fprintln(file, "")
	mutex.Unlock()
}

func CreateLogFile() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirlog := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirLog))
	if _, err := os.Stat(dirlog); os.IsNotExist(err) {
		err := os.MkdirAll(dirlog, 0711)

		if err != nil {
			fmt.Println("cannot create logs dir")
			os.Exit(1)
		}
	}
	t := time.Now()
	ft := t.Format("2006-01-02")
	FileLog = Filelog(filepath.FromSlash(fmt.Sprintf("%s/log_med_server_%v.log", dirlog, ft)))
}

func CreateTempDir() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirtemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirTemp))
	if _, err := os.Stat(dirtemp); os.IsNotExist(err) {
		err := os.MkdirAll(dirtemp, 0711)

		if err != nil {
			fmt.Println("cannot create temp dir")
			os.Exit(1)
		}
	}
}
func CreateListSites() {
	sites = make([]Site, 0)
	sites = append(sites, Site{Alias: "galaktika.clinic", Url: "http://galaktika.clinic/prices/", FileName: "galaktika.clinic.csv"})
}

func CreateEnv() {
	CreateLogFile()
	CreateTempDir()
	GetPort()
	CreateListSites()
}
