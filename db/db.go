package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"os"
)

var db *sql.DB

func Start() {
	var err error
	db, err = sql.Open("postgres", determineDb())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Conectado\n")
}

func GetDB() (c *sql.DB) {
	return db
}


func determineDb() (string) {
	port := os.Getenv("DATABASE_URL")
	if port == "" {
		return "postgres://wrzuvgfdhwhxmd:d073a4d6d274727ada015a2ba9cac3d6f50a08afe9e2271a58409cfafcd16d68@ec2-23-21-158-253.compute-1.amazonaws.com:5432/d1dfnfc8hob5ov"
	}
	return port
}