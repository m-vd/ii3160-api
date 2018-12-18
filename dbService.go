package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//Constants for connecting to the database
const (
	dbUser     = "ii3160"
	dbPassword = "ii3160"
	dbName     = "ii3160_api"
)

//Mahasiswa is a data type for representing a student in the database
type Mahasiswa struct {
	NimTPB   int
	NimProdi int
	Nama     string
	EmailSTD string
	Email    string
}

//User is a data type for representing a user in the database
type User struct {
	Mail   string
	Cn     string
	Sn     string
	Nim    string
	APIKey string
}

//FindByNimProdi is a function to search a student by his/her nimProdi attribute
func FindByNimProdi(nimProdi int) string {
	//Connect to database
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var m Mahasiswa

	//Search the database for a row with the exact nim
	if err := db.QueryRow("SELECT * FROM mahasiswa WHERE nimprodi=$1", nimProdi).Scan(&m.NimTPB, &m.NimProdi, &m.Nama, &m.EmailSTD, &m.Email); err != nil {
		if err == sql.ErrNoRows {
			return "{}"
		}
		log.Fatal(err)
	}
	fmt.Printf(fmt.Sprintf("Query for %d completed\n", nimProdi))
	//convert data to JSON
	value, _ := json.MarshalIndent(m, "", " ")
	//return value
	return string(value)

}

//FindByNimTPB is a function to search a student by his/her nimTPB attribute
func FindByNimTPB(nimTPB int) string {
	//Connect to database
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var m Mahasiswa

	//Search the database for a row with the exact nim
	if err := db.QueryRow("SELECT * FROM mahasiswa WHERE nimtpb=$1", nimTPB).Scan(&m.NimTPB, &m.NimTPB, &m.Nama, &m.EmailSTD, &m.Email); err != nil {
		if err == sql.ErrNoRows {
			return "{}"
		}
		log.Fatal(err)
	}
	fmt.Printf(fmt.Sprintf("Query for %d completed\n", nimTPB))
	//convert data to JSON
	value, _ := json.MarshalIndent(m, "", " ")
	//return value
	return string(value)

}

//AddNewUser is a function to add a new user to the database
func AddNewUser(u User) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlq := "INSERT INTO users (mail, cn, sn, nim, api_key) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(sqlq, u.Mail, u.Cn, u.Sn, u.Nim, u.APIKey)

	if err != nil {
		return false
	}
	return true
}

//CheckUserExistByMail is a function to check whether a user exist or not in the database
//using on email as parameter
func CheckUserExistByMail(mail string) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM user WHERE mail=$1)", mail).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Error checking user by mail")
	}
	return exists

}
