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
	http.HandleFunc("/:any", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404 Not Found")
	})
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":3160", nil))
}
