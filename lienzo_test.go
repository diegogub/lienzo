package main

import (
	"testing"
)

func TestPersist(t *testing.T) {
	m := NewMap()

	tinfo := NewTInfo()
	tinfo.Tmpl = "2adadasd1main"
	tinfo.Data["miau"] = 1312312
	tinfo.Data["t"] = []string{"est", "asdasdasd"}
	m["/test"] = tinfo

	err := m.Persist("test2.json")
	if err != nil {
		panic(err)
	}

}
