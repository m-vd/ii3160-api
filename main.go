package main

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"net/http"
	"net/http/httputil"
	"strconv"
	"io/ioutil"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found")
		return
	} else {
		t, _ := template.ParseFiles("static/hello.html")
		t.Execute(w, "Golang for serving static content")
	}
	//this should be changed since hello.html is a static asset
}

func main() {
	ReadJSON("fillerData.json")
	file, _ := ioutil.ReadFile("key")

	keyAPI := string(file)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		} else {
			keyValue, found := queryValues["key"]
			if !found {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "403 Forbidden: You are forbidden to access without an API key.")
			} else {
				if keyValue[0] == keyAPI {
					parsedNIM, err := strconv.Atoi(queryValues["nim"][0])
					if err != nil {
						fmt.Fprintf(w, "{ }")
					} else {
						fmt.Fprintf(w, Find(parsedNIM))
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "403 Forbidden: You are forbidden to access without an API key.")
				}
			}
		}
	})

	http.HandleFunc("/", helloHandler)

	http.HandleFunc("/readClientRequest", func(w http.ResponseWriter, r *http.Request){
		a, _ := httputil.DumpRequest(r, true)
		fmt.Fprintf(w, string(a))
		fmt.Fprintf(w, r.URL.RawQuery)
	})
	log.Fatal(http.ListenAndServe(":3160", nil))
}
