package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
	"strconv"
)

type Config struct {
	Markable bool
	ServerRoot string
	Port string
	UnmarkableXML string
	MarkableXML string
	PathToGoogleKeyJson string
	SpreadsheetId string
}

// added for test commit

var knownKeys = []string{"ref_sid", "event.id", "event.order", "subscriber", "abonent", "protocol", "user_id", "service", "event.text", "event.referer", "event", "lang", "serviceId", "wnumber"}

func init_system() (*Config, []byte, []byte, error) {
	config := new(Config)
	cfg_bytes, err := ioutil.ReadFile("cfg.json")
	json.Unmarshal(cfg_bytes, config)
	if !exists("out.csv") {
		ioutil.WriteFile("out.csv", []byte("page,button,user_id,wnumber,protocol\n"), 0644)
	}
	//f, err := os.OpenFile("out.csv", os.O_APPEND|os.O_WRONLY, 0600)
	resp_xml, err := ioutil.ReadFile(config.UnmarkableXML)
	mark_resp_xml, err := ioutil.ReadFile(config.MarkableXML)
	initialize_sheet()
	return config, resp_xml, mark_resp_xml, err
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
	//for parameter := range r.URL.Query() {
		//if contains(config.PageNames, parameter) {
			params:=genParameters(r.URL.Query())
			mark = calcMark(params)
			go addEntry(time.Now().String(),
				r.URL.Query().Get("subscriber"),
				r.URL.Query().Get("protocol"),
				r.URL.Query().Get("wnumber"),
					config.Markable,
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
	if !config.Markable {
		fmt.Fprintf(w, string(responseXml))
	}else{
		fmt.Fprintf(w, string(markableResponseXml), strconv.Itoa(mark))
	}
}

var config = new(Config)
//var outputFile = new(os.File)
var responseXml = []byte{}
var markableResponseXml = []byte{}

func main() {
	log.Println("Starting...")
	cfg, respXml, markRespXml, err := init_system()
	config = cfg
	//outputFile = f
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
