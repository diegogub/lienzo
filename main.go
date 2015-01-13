package main

import (
	"flag"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	r *mux.Router
	m TemplateMap
)

var (
	file = flag.String("config", "lienzo.json", "Configuration file")
	dir  = flag.String("dir", ".", "Templates directory")
	port = flag.String("port", "8989", "default port")
	kind = flag.String("suffix", ".html", "documents suffix to load as template")
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
			tmps.ParseFiles(path)
		}
		return nil
	})

	err := tmps.ExecuteTemplate(w, tinfo.Tmpl, tinfo.Data)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	r = mux.NewRouter()
	m = LoadMap(*file)

	log.Println("Loaded File:", m)
	m.Routes(r, Handler)

	err := http.ListenAndServe(":"+*port, r)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
