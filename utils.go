package main

import (
	"os"
	"io/ioutil"
	"bytes"
	"net/url"
)

func exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

func file_contains(fileName string, what string) (bool, error) {
	data, err := ioutil.ReadFile(fileName)
	return bytes.Contains(data, []byte(what)), err
}

func contains(arr []string, elem string) (bool){
	for _, element := range arr{
		if element == elem{
			return true
		}
	}
	return false
}

func genParameters(query url.Values) (map[string]string){
	out := map[string]string{}
	for parmaeter := range query{
		if !contains(config.KnownKeys, parmaeter){
			out[parmaeter] = query.Get(parmaeter)
		}
	}
	return out
}