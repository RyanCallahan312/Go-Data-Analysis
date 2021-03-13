package migration

import (
	"errors"
	"fmt"
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

// InitalizeDB inalized the db with given name
func InitalizeDB(name string) {
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
	_ = db.QueryRow(fmt.Sprintf(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = '%s')`, name)).Scan(&dbExists)
	if !dbExists {
		_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, name))
		if err != nil {
			log.Fatalln(err)
		}
	}

}
