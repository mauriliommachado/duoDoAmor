package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
)

var db *sql.DB

func Start() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:admin@localhost/duo?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Conectado\n")
}

func GetDB() (c *sql.DB) {
	return db
}
