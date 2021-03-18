package migration

import (
	"Project1/config"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/semver/v3"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// Versions all DB versions
var Versions = []string{
	"v0.0.0",
	"v0.0.1",
	"v1.0.0",
}

// MigrateToLatest Migrates the database to the latest version
func MigrateToLatest() error {
	lastMigration := GetLastMigration()
	if lastMigration.Version == Versions[len(Versions)-1] {
		return nil
	}
	constraint := MakeConstraint(lastMigration, true, Versions[len(Versions)-1])
	return UpdateDB(lastMigration, true, constraint)
}

// MigrateToVersion Migrates the database to a specified version
func MigrateToVersion(targetVersion string) error {

	isValidVersion := false
	for _, v := range Versions {
		if targetVersion == v {
			isValidVersion = true
		}
	}

	if !isValidVersion {
		return errors.New("Invalid version number")
	}

	lastMigration := GetLastMigration()

	if targetVersion == lastMigration.Version {
		return nil
	}

	isBuildConstriant, err := semver.NewConstraint("< " + targetVersion)
	if err != nil {
		return err
	}

	lastMigrationVersion, err := semver.NewVersion(lastMigration.Version)
	if err != nil {
		return err
	}

	isBuild := isBuildConstriant.Check(lastMigrationVersion)

	constraint := MakeConstraint(lastMigration, isBuild, targetVersion)
	err = UpdateDB(lastMigration, isBuild, constraint)
	if err != nil {
		return err
	}

	return nil
}

// InitalizeDB inalized the db with given name
func InitalizeDB(name string) {
	db, err := sqlx.Open("pgx", config.Env["MAINTENANCE_CONNECTION_STRING"])
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	var dbExists bool
	_ = db.QueryRow(fmt.Sprintf(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = '%s')`, name)).Scan(&dbExists)
	if !dbExists {
		_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, name))
		if err != nil {
			log.Panic(err)
		}
	}

}
