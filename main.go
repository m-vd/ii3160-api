package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
				if CheckAPIKey(keyValue[0]) {
					//Check for NIM query
					nimValue, found := queryValues["nim"]
					if !found {
						nameValue, found := queryValues["nama"]
						if !found {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprintf(w, "400 Bad Request: You have to provide a search query.")
						} else {
							fmt.Fprintf(w, FindByNama(nameValue[0]))
						}
					} else {
						parsedNIM, err := strconv.Atoi(nimValue[0])
						if err != nil {
							fmt.Fprintf(w, "{}")
						} else {
							if FindByNimProdi(parsedNIM) == "{}" {
								fmt.Fprintf(w, FindByNimTPB(parsedNIM))
							} else {
								fmt.Fprintf(w, FindByNimProdi(parsedNIM))
							}
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
			validationResponse, err := http.Get("https://login.itb.ac.id/cas/serviceValidate?service=https://" + url.QueryEscape(r.Host+"/login") + "&ticket=" + url.QueryEscape(ticket[0]))
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
				u := User{
					Mail:   parsed.Mail,
					Cn:     parsed.Cn,
					Sn:     parsed.Sn,
					Nim:    parsed.ItbNIM,
					APIKey: GenerateAPIKey(),
				}

				//Create user and send API key
				if AddNewUser(u) {
					fmt.Printf("%+v", u)
					fmt.Fprintf(w, "Please take note of the API Key, as you will be only given this key once.\n")
					fmt.Fprintf(w, u.APIKey)
				}
			} else {
				fmt.Println("Redirecting to API")
				http.Redirect(w, r, "/api", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, ("https://login.itb.ac.id/cas/login?service=https://" + url.QueryEscape(r.Host+"/login")), http.StatusFound)
		}
	}
}

func main() {
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/login", authHandler)
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")

	if port != "" {
		log.Fatal(http.ListenAndServe(":"+port, nil))
		fmt.Println("Running on port " + port)
	} else {
		log.Fatal(http.ListenAndServe(":3160", nil))
	}
}
