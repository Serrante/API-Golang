package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	USER    = "blockcoin"
	PASS    = "bl0ckc01n"
	DBNAME  = "blockcoin"
	SSLMODE = "disable"
)

func Connect() *sql.DB {
	URL := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", USER, PASS, DBNAME, SSLMODE)
	db, err := sql.Open("postgres", URL)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func TestConnection() {
	con := Connect()
	defer con.Close()
	err := con.Ping()

	if err != nil {
		err := fmt.Errorf("%s", err.Error())
		fmt.Printf("%s", err)
		return
	}

	fmt.Println("Database connected!")
}
