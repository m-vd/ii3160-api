package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Mahasiswa struct {
	NimTPB   int
	NimProdi int
	Nama     string
	EmailSTD string
	Email    string
}

var mahasiswa []Mahasiswa

func ReadJSON(filename string) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Print(err)
	}

	jsonData := string(file)
	json.Unmarshal([]byte(jsonData), &mahasiswa)

}

func Find(nim int) string {
	var returnValue []byte
	for _, element := range mahasiswa {
		if element.NimProdi == nim {
			foundData, err := json.MarshalIndent(element, "", "	")
			if err != nil {
				log.Fatal(err)
			}
			returnValue = foundData
		}
	}
	if returnValue != nil {
		return (string(returnValue))
	} else {
		return "Not Found"
	}

}
