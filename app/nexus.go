package main

import (
	"html/template"
	"net/http"
)

func main() {
	main := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/index.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", main)

	http.ListenAndServe("localhost:8080", nil)
}
