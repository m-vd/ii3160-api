package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//Constants for connecting to the database

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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var m Mahasiswa

	//Search the database for a row with the exact nim
	if err := db.QueryRow("SELECT * FROM mahasiswa WHERE nimtpb=$1", nimTPB).Scan(&m.NimTPB, &m.NimProdi, &m.Nama, &m.EmailSTD, &m.Email); err != nil {
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

//FindByNama is a function to find students by their names
func FindByNama(nama string) string {
	//Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var m Mahasiswa
	var ma []Mahasiswa

	//Search the database for a row with the exact nim
	rows, err := db.Query("SELECT * FROM mahasiswa WHERE nama LIKE '%' || $1 || '%'", nama)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m.NimTPB, &m.NimProdi, &m.Nama, &m.EmailSTD, &m.Email); err != nil {
			log.Fatal(err)
		}
		ma = append(ma, m)
	}
	value, _ := json.MarshalIndent(ma, "", " ")
	//return value
	return string(value)

}

//AddNewUser is a function to add a new user to the database
func AddNewUser(u User) bool {
	//Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlq := "INSERT INTO users (mail, cn, sn, nim, api_key) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(sqlq, u.Mail, u.Cn, u.Sn, u.Nim, u.APIKey)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//CheckUserExistByMail is a function to check whether a user exist or not in the database
//using on email as parameter
func CheckUserExistByMail(mail string) bool {
	//Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE mail=$1)", mail).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Error checking user by mail")

	}
	return exists

}

//CheckAPIKey is a function to compare whether the given parameter of key exists as an api_key in the database
func CheckAPIKey(key string) bool {
	//Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var k string
	err = db.QueryRow("SELECT api_key FROM users WHERE api_key=$1", key).Scan(&k)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
