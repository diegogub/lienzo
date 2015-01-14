package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Handle all dummy api request, load file or do an http request
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			ErrorHTML(e, w)
		}
	}()
	// get dummy info
	if _, ok := a[r.URL.Path]; !ok {
		http.NotFound(w, r)
		return
	}
	dinfo, _ := a[r.URL.Path]
	switch dinfo.Type {
	case "url", "file":
	default:
		panic(errors.New("Invalid dummy type, must be 'url' or 'file'"))
	}

	if dinfo.Type == "file" {
		f, err := os.Open(dinfo.File)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		w.Write(b)
	} else {
		b, err := dinfo.DoRequest()
		if err != nil {
			panic(err)
		}
		w.Write(b)
	}
}

// map endpoints to data sources
type DummyMap map[string]DummyInfo

func (dm DummyMap) DummyRoutes(r *mux.Router) {
	for endpoint, _ := range dm {
		log.Println(endpoint)
		r.HandleFunc(endpoint, DummyHandler)
	}
}

func NewDMap() DummyMap {
	dm := make(map[string]DummyInfo)
	return dm
}

// Load Dummy map, :O note: code a generic util for json loading
func LoadDummy(path string) DummyMap {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to parse dummy file:", err.Error())
	}
	dm := NewDMap()
	err = json.Unmarshal(f, &dm)
	if err != nil {
		log.Fatalf("Failed to unmarshal:", err.Error())
	}
	log.Println("->>", dm)
	return dm
}

type DummyInfo struct {
	// url or file
	Type     string `json:"type"`
	Url      string `json:"url"`
	Method   string `json:"method"`
	Data     []byte `json:"data"`
	BodyType string `json:"bodyType"`
	File     string `json:"file"`
}

func (di DummyInfo) DoRequest() ([]byte, error) {
	var res *http.Response
	var err error
	switch di.Method {
	case "POST":
		res, err = http.Post(di.Url, di.BodyType, bytes.NewBuffer(di.Data))
	default:
		res, err = http.Get(di.Url)
	}
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	return b, err
}
