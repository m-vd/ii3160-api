package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//Deny access to other paths
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found")
		return
	}
	//Render index
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	} else {
		t.Execute(w, "")
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("key")
	keyAPI := string(file)
	//Only handle GET requests
	switch r.Method {
	case http.MethodGet:
		queryValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		} else {
			//Check for API Key
			keyValue, found := queryValues["key"]
			if !found {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "403 Forbidden: You are forbidden to access without an API key.")
			} else {
				if keyValue[0] == keyAPI {
					//Check for NIM query
					nimValue, found := queryValues["nim"]
					if !found {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(w, "400 Bad Request: You have to provide a NIM value.")
					} else {
						parsedNIM, err := strconv.Atoi(nimValue[0])
						if err != nil {
							fmt.Fprintf(w, "{}")
						} else {
							fmt.Fprintf(w, FindByNimProdi(parsedNIM))
						}
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "403 Forbidden: You are forbidden to access without an API key.")
				}
			}
		}
	//Deny any requests other than GET.
	default:
		fmt.Fprintf(w, "404 Not Found")
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	} else {
		ticket, ticketPresent := query["ticket"]
		if ticketPresent {
			validationResponse, err := http.Get("https://login.itb.ac.id/cas/serviceValidate?service=" + url.QueryEscape(r.Host+"/login") + "&ticket=" + url.QueryEscape(ticket[0]))
			if err != nil {
				log.Fatal(err)
			}
			validationBody, err := ioutil.ReadAll(validationResponse.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer validationResponse.Body.Close()

			parsed := ParseResponseXML(validationBody)

			if !CheckUserExistByMail(parsed.Mail) {
				//create api key
				u := User{
					Mail:   parsed.Mail,
					Cn:     parsed.Cn,
					Sn:     parsed.Sn,
					Nim:    parsed.ItbNIM,
					APIKey: parsed.Mail,
				}

				//Create user and send API key
				if AddNewUser(u) {
					fmt.Fprintf(w, u.APIKey)
				}
			}

		} else {
			http.Redirect(w, r, ("https://login.itb.ac.id/cas/login?service=" + url.QueryEscape(r.Host+"/login")), http.StatusFound)
		}
	}
}

func main() {
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/login", authHandler)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":3160", nil))
}
