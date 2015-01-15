package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	r *mux.Router
	m TemplateMap
	a DummyMap
)

var (
	file      = flag.String("config", "lienzo.json", "Configuration file")
	apis      = flag.String("api", "dummies.json", "Configuration file for dummy services")
	dir       = flag.String("dir", ".", "Templates directory")
	assets    = flag.String("assets", "", "Assets directory name to serve static files, must be in same the folder")
	port      = flag.String("port", "8989", "default port")
	kind      = flag.String("suffix", ".html", "documents suffix to load as template")
	admin     = flag.Bool("admin", false, "turns web admin on")
	adminport = flag.String("aport", "8990", "admin port")
)

func main() {
	flag.Parse()
	if *assets == "" {
		log.Println("WARNING - Assets directory not set! set using flag:  ./lienzo -assets static")
	}

	r = mux.NewRouter()
	// Load all configutation files as maps in memory
	m = LoadMap(*file)
	log.Println("Loaded File:", m)
	m.Routes(r)
	a = LoadDummy(*apis)
	a.DummyRoutes(r)
	log.Println("dummies:", a)

	if *assets != "" {
		r.PathPrefix("/" + *assets + "/").Handler(http.StripPrefix("/"+*assets+"/",
			http.FileServer(http.Dir(*assets+"/"))))
	}

	if *admin {
		go AdminInit(*adminport)
	}

	err := http.ListenAndServe(":"+*port, r)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
