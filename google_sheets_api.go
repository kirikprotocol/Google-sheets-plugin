package main

import (
	"io/ioutil"
	"gopkg.in/Iwark/spreadsheet.v2"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"log"
	"strconv"
	"regexp"
)

var service = new(spreadsheet.Service)
var sprsheet = *new(spreadsheet.Spreadsheet)

func initialize_sheet() {
	data, err := ioutil.ReadFile(config.PathToGoogleKeyJson)
	if err != nil{
		log.Fatal("Cannot read from json key file: ",config.PathToGoogleKeyJson,";")
	}
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())
	service = spreadsheet.NewServiceWithClient(client)
}

func updSheet(id string)(error){
	err := *new(error)
	sprsheet, err = service.FetchSpreadsheet(id)
	return err
}

func addEntry(timestamp string, user_id string, protocol string, wnumber string, tgId string, markable bool, mark int, params map[string]string) {
	log.Print("Adding entry: ", timestamp, " ", user_id, " ", protocol, " ", wnumber, " ", params)
	sheet, err := sprsheet.SheetByIndex(0)
	checkError(err)
	err = sheet.Synchronize()
	checkError(err)
	log.Print(".")
	//log.Println("sheet rows [0]=",sheet.Rows[0])
	if sheet.Rows[0][0].Value != "timestamp"{
		sheet.Update(0,0, "timestamp")
		sheet.Update(0,1, "user_id")
		sheet.Update(0,2,"telegram_id")
		sheet.Update(0,3, "protocol")
		sheet.Update(0,4, "wnumber")
		sheet.Update(0,5, "mark")
		//err = sheet.Synchronize()
		//checkError(err)
		//go addEntry(timestamp,user_id,protocol,wnumber,markable,mark,params)
		//return
	}
	//log.Println("rows: ",sheet.Rows)
	err = sheet.Synchronize()
	checkError(err)
	pgNamesCells := sheet.Rows[0][6:]
	pgNames := []string{}
	for _, cell := range pgNamesCells {
		pgNames = append(pgNames, cell.Value)
	}

	emptyRowIdx := getEmpty(sheet.Columns[0])
	emptyColumnIdx := getEmpty(sheet.Rows[0])

	log.Print(".")
	r := regexp.MustCompile(`^[A-Za-z]`)
	match :=!r.MatchString(user_id)
	if match {
		sheet.Update(emptyRowIdx, 0, timestamp)
		sheet.Update(emptyRowIdx, 1, user_id)
		sheet.Update(emptyRowIdx, 2, protocol)
		sheet.Update(emptyRowIdx, 3, tgId)
		sheet.Update(emptyRowIdx, 4, wnumber)
	}else {
		sheet.Update(emptyRowIdx, 0, timestamp)
		sheet.Update(emptyRowIdx, 1, "0")
		sheet.Update(emptyRowIdx, 2, protocol)
		sheet.Update(emptyRowIdx, 3, tgId)
		sheet.Update(emptyRowIdx, 4, user_id)
	}
	//123
	if markable {
		sheet.Update(emptyRowIdx, 5, strconv.Itoa(mark))
	}else {
		sheet.Update(emptyRowIdx, 5, "not markable")
	}
	for key, value := range params {
		//log.Println("Cols: ",sheet.Columns[0],"; rows: ",sheet.Rows[0])
		emptyRowIdx = getEmpty(sheet.Columns[0]) - 1
		emptyColumnIdx = getEmpty(sheet.Rows[0])
		//log.Println("iterating params: key:",key,"=",value,";")
		if !contains(pgNames, key) {
			//log.Println("!Contains; len(sheet.Rows[0])=",emptyColumnIdx)
			fillColumn(sheet, emptyColumnIdx, emptyRowIdx)
			fillUnfilledCols(sheet, emptyRowIdx, emptyColumnIdx)
			sheet.Update(0, emptyColumnIdx, key)
			//log.Println("Setting value")
			sheet.Update(emptyRowIdx, emptyColumnIdx, value)
			//log.Println("ended updating...")
		} else {
			//log.Println("Filling: ",emptyRowIdx, findColumn(sheet, key), value)
			fillUnfilledCols(sheet, emptyRowIdx, emptyColumnIdx)
			sheet.Update(emptyRowIdx, findColumn(sheet, key), value)
		}
	}
	log.Print(".")
	err = sheet.Synchronize()
	checkError(err)
	log.Println("Done!")
}

func getEmpty(cells []spreadsheet.Cell) (int) {
	out := len(cells)
	for i, cell := range cells {
		//log.Println("cell[",i,"]=",cell.Value)
		if cell.Value == "" {
			//log.Println("Found emty cell!(",i,")")
			out = i
			break
		}
	}
	return out
}

func findColumn(sheet *spreadsheet.Sheet, key string) (int) {
	out := 0
	for i, column := range sheet.Rows[0][6:] {
		if column.Value == key {
			out = i
		}
	}
	return 5 + out
}

func fillUnfilledCols(sheet *spreadsheet.Sheet, row int, lastColumnIdx int){
	for i:=0; i<=lastColumnIdx-1; i++{
		if sheet.Rows[row][i].Value == "" {
			sheet.Update(row, i, "0")
		}
	}
}

func fillColumn(sheet *spreadsheet.Sheet, columnIdx int, emptyRowIdx int) {
	//log.Println("Filling column(length = ",emptyRowIdx,") ",columnIdx)
	for i := 0; i <= emptyRowIdx; i++ {
		sheet.Update(i, columnIdx, "0")
	}
	//log.Println("Done!")
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
