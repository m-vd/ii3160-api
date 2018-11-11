package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/hello.html")
	t.Execute(w, "Golang for serving static content")
	//this should be changed since hello.html is a static asset
}

func main() {
	ReadJSON("fillerData.json")
	http.HandleFunc("/apitester", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, Find(18216017))
	})
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":3160", nil))
}
