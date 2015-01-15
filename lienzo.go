package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			ErrorHTML(e, w)
		}
	}()
	log.Println(r.URL.Path)
	// parse file
	tinfo, ok := m[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
	}
	tmps := template.New("tmp")
	// Load all templates in every request! . We don't need performace,just to load template
	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, *kind) {
			log.Println("Matched files:", path)
			tmps = template.Must(tmps.ParseFiles(path))
		}
		return nil
	})

	if _, ok := m["global"]; ok {
		tinfo.Sum(m["global"])
	}
	err := tmps.ExecuteTemplate(w, tinfo.Tmpl, tinfo.Data)
	if err != nil {
		log.Println(err)
	}
}

type TemplateMap map[string]TemplateInfo

func (tm TemplateMap) String() string {
	b, _ := json.Marshal(tm)
	return string(b)
}

// Persist template map to file
func (tm TemplateMap) Persist(file string) error {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.WriteString(tm.String())
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Sum 2 maps

func (tm TemplateMap) Routes(r *mux.Router) {
	for route, _ := range tm {
		r.HandleFunc(route, Handler)
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

func NewTInfo() TemplateInfo {
	var tmpl TemplateInfo
	tmpl.Data = make(map[string]interface{})
	return tmpl
}

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
	log.Println("Reloading config..")
	newa := LoadDummy(*apis)
	if newa == nil {
		log.Println("Failed to reload")
	}
	a = newa
}
