package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
	"strconv"
	"os"
	"io"
)

type Config struct {
	ServerRoot string
	Port string
	UnmarkableXML string
	MarkableXML string
	ErrorXML string
	LogPath string
	PathToGoogleKeyJson string
	KnownKeys []string
}

var config = new(Config)
//var outputFile = new(os.File)
var responseXml = []byte{}
var markableResponseXml = []byte{}
var errorXml = []byte{}

// added for test commit

//var knownKeys = []string{"ref_sid", "event.id", "event.order", "subscriber", "abonent", "protocol", "user_id", "service", "event.text", "event.referer", "event", "lang", "serviceId", "wnumber"}

func init_system() (*Config, []byte, []byte, []byte, error) {
	cfg_bytes, err := ioutil.ReadFile(os.Args[1])
	json.Unmarshal(cfg_bytes, config)
	//log.Println("config: ",config)
	/*
	if !exists("out.csv") {
		ioutil.WriteFile("out.csv", []byte("page,button,user_id,wnumber,protocol\n"), 0644)
	}
	*/
	//f, err := os.OpenFile("out.csv", os.O_APPEND|os.O_WRONLY, 0600)
	resp_xml, err := ioutil.ReadFile(config.UnmarkableXML)
	mark_resp_xml, err := ioutil.ReadFile(config.MarkableXML)
	err_xml, err := ioutil.ReadFile(config.ErrorXML)

	logFile, err := os.OpenFile(config.LogPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Println("Logging to file and console!")

	initialize_sheet()
	return config, resp_xml, err_xml, mark_resp_xml, err
}

func calcMark(params map[string]string) (int){
	out := 0
	for _, value := range params{
		num, err := strconv.Atoi(value)
		if err == nil{
			out+=num
		}else if value == "0"{
			out+=0
		}
	}
	return out
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	var mark = 0
	evaluableInt, err := strconv.Atoi(r.URL.Query().Get("evaluable"))
	evaluable := true
	evaluable = evaluableInt == 1
	if err != nil || evaluableInt > 1{
		evaluable = false
	}
	//for parameter := range r.URL.Query() {
		//if contains(config.PageNames, parameter) {
			params:=genParameters(r.URL.Query())
			mark = calcMark(params)
			updErr := updSheet(r.URL.Query().Get("spreadsheetId"))
			if updErr != nil{
				fmt.Fprintf(w, string(errorXml), string(updErr.Error()))
				return
			}
			go addEntry(time.Now().String(),
				r.URL.Query().Get("subscriber"),
				r.URL.Query().Get("protocol"),
				r.URL.Query().Get("wnumber"),
					evaluable,
					mark,
					params)
			//go outputFile.Write([]byte(time.Now().String() + "," +
			//	parameter + "," +
			//	r.URL.Query().Get(parameter) + "," +
			//	r.URL.Query().Get("subscriber") + "," +
			//	r.URL.Query().Get("wnumber") + "," +
			//	r.URL.Query().Get("protocol") + "\n"))
		//}
	//}
	if !evaluable {
		fmt.Fprintf(w, string(responseXml))
	}else{
		fmt.Fprintf(w, string(markableResponseXml), strconv.Itoa(mark))
	}
}

func main() {
	log.Println("Starting...")
	if len(os.Args) < 2{
		log.Fatal("You should pass me a config name like: ",os.Args[0]," <json config name>")
	}
	cfg, respXml, errXml, markRespXml, err := init_system()
	config = cfg
	//outputFile = f
	errorXml = errXml
	responseXml = respXml
	markableResponseXml = markRespXml
	//log.Println(string(response_xml))
	if err != nil {
		//outputFile.Close()
		panic(err)
	}
	log.Println("Done! Listening...")
	http.HandleFunc(config.ServerRoot, handler)
	http.ListenAndServe(":"+config.Port, nil)
}
