package migration

import (
	"Project1/database"
	"errors"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/jmoiron/sqlx"
)

// Version all DB versions
var Versions = []string{
	"v0.0.0",
	"v0.0.1",
	"v1.0.0",
}

// MigrateToLatest Migrates the database to the latest version
func MigrateToLatest() {
	initalizeDB()
	lastMigration := GetLastMigration()
	constraint := MakeConstraint(lastMigration, true, Versions[len(Versions)-1])
	UpdateDB(lastMigration, true, constraint)
}

// MigrateToVersion Migrates the database to a specified version
func MigrateToVersion(targetVersion string) error {

	isValidVersion := false
	for _, v := range Versions {
		if targetVersion == v {
			isValidVersion = true
		}
	}

	if isValidVersion == false {
		return errors.New("Invalid version number")
	}

	initalizeDB()
	lastMigration := GetLastMigration()

	isBuildConstriant, err := semver.NewConstraint("< " + targetVersion)
	if err != nil {
		log.Fatalln(err)
	}

	lastMigrationVersion, err := semver.NewVersion(lastMigration.Version)
	if err != nil {
		log.Fatalln(err)
	}

	isBuild := isBuildConstriant.Check(lastMigrationVersion)

	constraint := MakeConstraint(lastMigration, isBuild, targetVersion)
	UpdateDB(lastMigration, isBuild, constraint)

	return nil
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
			WHERE datname = $1
			)`, os.Getenv("DATABASE_NAME")).Scan(&dbExists)

	if !dbExists {
		_, err = db.Exec(`CREATE DATABASE $1`, os.Getenv("DATABASE_NAME"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	database.DB, err = sqlx.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING"))

}
