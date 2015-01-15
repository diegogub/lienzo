package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func AdminInit(port string) {
	ar := mux.NewRouter()

	ar.HandleFunc("/", Redirect).Methods("GET")

	ar.HandleFunc("/panel", AdminPanel)

	log.Println("Starting admin..")
	err := http.ListenAndServe(":"+port, ar)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func AdminPanel(w http.ResponseWriter, r *http.Request) {
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/panel", 302)
}
