package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type TemplateMap map[string]TemplateInfo

func (tm TemplateMap) String() string {
	b, _ := json.Marshal(tm)
	return string(b)
}

// Sum 2 maps
func (tm TemplateInfo) Sum(t TemplateInfo) {
	for key, value := range t.Data {
		if _, ok := tm.Data[key]; !ok {
			tm.Data[key] = value
		} else {
			log.Println("warn: overwritting data key: "+key+":", value)
		}
	}
	return
}

func (tm TemplateMap) Routes(r *mux.Router, hand func(http.ResponseWriter, *http.Request)) {
	for route, _ := range tm {
		r.HandleFunc(route, hand)
	}
	r.HandleFunc("/reload", ReloadConfig).Methods("PUT")
	return
}

func NewMap() TemplateMap {
	var t TemplateMap
	t = make(map[string]TemplateInfo)
	return t
}

type TemplateInfo struct {
	Tmpl string                 `json:"tmpl"`
	Data map[string]interface{} `json:"data"`
}

// Load from file
func LoadMap(path string) TemplateMap {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to parse file:", err.Error())
	}
	m := NewMap()
	err = json.Unmarshal(f, &m)
	if err != nil {
		log.Fatalf("Failed to unmarshal:", err.Error())
	}
	return m
}

func ReloadConfig(w http.ResponseWriter, r *http.Request) {
	// Reload config
	log.Println("Reloading config..")
	newm := LoadMap(*file)
	if newm == nil {
		log.Println("Failed to reload")
	}
	m = newm
}
