package migration

import (
	"Project1/database"
	"database/sql"
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/jmoiron/sqlx"
)

// GetLastMigration queries the migration table to find the current db version. if no results are found it defaults to V0.0.0
func GetLastMigration() MigrationModel {
	var lastMigration MigrationModel
	err := database.DB.QueryRowx(`SELECT * FROM migration WHERE is_on_version`).StructScan(&lastMigration)
	if err != nil {
		if err != sql.ErrNoRows {
			_, err = database.DB.Exec(`CREATE TABLE IF NOT EXISTS migrations(
				migration_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
				version VARCHAR(16),
				is_on_version BOOLEAN)`)
			if err != nil {
				log.Fatalln(err)
			}
		}
		return MigrationModel{0, "v0.0.0", true}
	}
	return lastMigration
}

// UpdateDB will either roll back or build the db to a specified version
func UpdateDB(lastMigration MigrationModel, isBuild bool, constraint *semver.Constraints) {

	semanticVersions := getAvailableVersions()

	for i, version := range semanticVersions {
		if isBuild {
			if constraint.Check(version) {
				runMigration(version.Original(), isBuild)
			}
		} else {
			reverseOrder := len(semanticVersions) - 1 - i
			if constraint.Check(semanticVersions[reverseOrder]) {
				runMigration(semanticVersions[reverseOrder].Original(), isBuild)
			}
		}
	}
	updateMigrationTable(isBuild, lastMigration)
}

func MakeConstraint(lastMigration MigrationModel, isBuild bool, targetVersion string) *semver.Constraints{
	var constraint *semver.Constraints
	var err error
	if isBuild {
		constraint, err = semver.NewConstraint("> " + lastMigration.Version + ", < " + targetVersion)
	} else {
		constraint, err = semver.NewConstraint("< " + lastMigration.Version + ", > " + targetVersion)
	}
	if err != nil {
		log.Fatalln(err)
	}

	return constraint
}

func updateMigrationTable(isBuild bool, lastMigration MigrationModel) {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for i, version := range Versions {
		var isCurrentVersion bool
		if isBuild {
			isCurrentVersion = i == len(Versions)-1
		} else {
			isCurrentVersion = i == 0
		}

		var migration MigrationModel
		migrationRow := tx.QueryRowx(`SELECT * FROM migrations WHERE version=$1`, version)

		err = migrationRow.StructScan(&migration)
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO migrations VALUES(DEFAULT, $1, $2)`, version, isCurrentVersion)
			if err != nil {
				log.Fatalln(err)
			}
		} else if migration.IsOnVersion != isCurrentVersion {
			_, err = tx.Exec(`UPDATE migrations SET is_on_version=$1 WHERE version='$1'`, isCurrentVersion, version)
			if err != nil {
				log.Fatalln(err)
			}
		}

	}

}

func runMigration(version string, isBuild bool) {
	switch version {
	case "v0.0.1":
		runInTransaction(isBuild, MigrateV0_0_1)
	case "v1.0.0":
		runInTransaction(isBuild, MigrateV1_0_0)
	}
}

func runInTransaction(isBuild bool, scriptToRun func(bool, *sqlx.Tx)) {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	scriptToRun(isBuild, tx)
}

func getAvailableVersions() []*semver.Version {
	semanticVersions := make([]*semver.Version, len(Versions))
	for i, r := range Versions {
		semanticVersion, err := semver.NewVersion(r)
		if err != nil {
			log.Fatalln(err)
		}

		semanticVersions[i] = semanticVersion
	}
	return semanticVersions
}
