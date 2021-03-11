package main

import (
	"database/sql"
	"log"
	"os"
)

func InitalizeDB() {
	conn, err := sql.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {

		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var dbExists bool
	_ = conn.QueryRow(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = 'comp490project1'
			)`).Scan(&dbExists)

	if !dbExists {
		_, err = conn.Exec(`CREATE DATABASE comp490project1`)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
