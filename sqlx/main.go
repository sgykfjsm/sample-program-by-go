package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

const (
	schemaPlace = `CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer);`
	schemaPhoneBook = `CREATE TABLE phone_book (
    telcode integer,
	name    text);`
	insertCityTelcode        = `INSERT INTO place (country, telcode) VALUES (?, ?)`
	insertCountryCityTelcode = `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
	insertTelcodeName        = `INSERT INTO phone_book (telcode, name) VALUES (?, ?)`
	selectAll                = `SELECT country, city, telcode FROM place`
	selectQueryx             = `SELECT * FROM place`
	selectQueryRow           = `SELECT * FROM place WHERE telcode = ?`
	selectJoinQueryx         = `SELECT b.name, a.country, a.city, a.telcode FROM place a JOIN phone_book b ON a.telcode = b.telcode`
)

type Place struct {
	Country       string
	City          sql.NullString
	TelephoneCode int `db:"telcode"`
}

type PlaceWithName struct {
	Name          string
	Country       string
	City          sql.NullString
	TelephoneCode int `db:"telcode"`
}

func FatalIfNotNil(message string, e error) {
	if e != nil {
		log.Fatalf(message, e.Error())
	}
}

func main() {
	db, err := sqlx.Open("sqlite3", ":memory:")
	FatalIfNotNil("Failed to open database: %s", err)
	fmt.Println("Opened the database connection")

	err = db.Ping()
	FatalIfNotNil("Failed to ping database: %s", err)
	fmt.Println("Pinging the database succesfully")

	_, err = db.Exec(schemaPlace)
	FatalIfNotNil("Failed to create table: %s", err)
	fmt.Println("Created Schema Place")

	_, err = db.Exec(schemaPhoneBook)
	FatalIfNotNil("Failed to create table: %s", err)
	fmt.Println("Created Schema PhoneBook")

	db.MustExec(insertCityTelcode, "Hong Kong", 852)
	db.MustExec(insertCityTelcode, "Singapore", 65)
	db.MustExec(insertCountryCityTelcode, "South Africa", "Johannesburg", 27)
	db.MustExec(insertTelcodeName, 852, "Alice")
	db.MustExec(insertTelcodeName, 65, "Bob")
	db.MustExec(insertTelcodeName, 27, "Charlie")

	// fetch all places from the db
	rows, err := db.Query(selectAll)
	fmt.Println("Selected the rows from database succesfully")

	// iterate over each row
	// treat the Rows like a database cursor rather than a materialized list of results
	for rows.Next() {
		var country string
		var city sql.NullString
		var telcode int
		err = rows.Scan(&country, &city, &telcode)
		FatalIfNotNil("Failed to scan the rows: %s", err)
		fmt.Printf("Row(Query) >>> country: %s, city: %s, telcode: %d\n", country, city.String, telcode)
	}

	rowx, err := db.Queryx(selectQueryx)
	for rowx.Next() {
		var p Place
		err = rowx.StructScan(&p)
		FatalIfNotNil("Failed to execute Queryx: %s", err)
		fmt.Printf("Row(Queryx) >>> country: %s, city: %s, telcode: %d\n", p.Country, p.City.String, p.TelephoneCode)
	}

	row := db.QueryRow(selectQueryRow, 852)
	var country string
	var city sql.NullString
	var telcode int
	err = row.Scan(&country, &city, &telcode)
	FatalIfNotNil("Failed to execute db.QueryRow(): %s", err)
	fmt.Printf("Row(QueryRow) >>> country: %s, city: %s, telcode: %d\n", country, city.String, telcode)

	var pn []PlaceWithName
	err = db.Select(&pn, selectJoinQueryx)
	FatalIfNotNil("Failed to execute JOIN query: %s", err)
	for i, p := range pn {
		fmt.Printf("Row(JOIN) >>> %d: name: %s,  country: %s, city: %s, telcode: %d\n", i, p.Name, p.Country, p.City.String, p.TelephoneCode)
	}
}
