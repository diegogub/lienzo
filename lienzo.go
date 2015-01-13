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

func (tm TemplateMap) Routes(r *mux.Router, hand func(http.ResponseWriter, *http.Request)) {
	for route, _ := range tm {
		r.HandleFunc(route, hand)
	}
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
	json.Unmarshal(f, &m)
	return m
}