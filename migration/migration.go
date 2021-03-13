package migration

import (
	"Project1/database"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func Migrate() {
	initalizeDB()
	lastMigration := GetLastMigration()
	UpdateDBFromVersion(lastMigration, true)
}

func initalizeDB() {
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

	database.DB, err = sqlx.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING"))

}
