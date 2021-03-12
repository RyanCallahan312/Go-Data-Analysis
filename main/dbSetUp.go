package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// InitalizeDB creates the database using the maintence connection string
func InitalizeDB() {
	db, err := sqlx.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {

		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var dbExists bool
	_ = db.QueryRow(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = 'comp490project1'
			)`).Scan(&dbExists)

	if !dbExists {
		_, err = db.Exec(`CREATE DATABASE comp490project1`)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
