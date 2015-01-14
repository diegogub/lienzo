package main

import (
	"html/template"
	"net/http"
)

var errtmpl string = `
<html>
<body>
<p>Error while loading resource:</p>
<b> {{ .Error }} </b>
</body>
</html>
`

func ErrorHTML(e interface{}, w http.ResponseWriter) {
	tmp, err := template.New("error").Parse(errtmpl)
	if err != nil {
		panic(err)
	}

	data := struct {
		Error interface{}
	}{
		e,
	}

	err = tmp.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
