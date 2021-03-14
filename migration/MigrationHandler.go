package migration

import (
	"Project1/database"
	"database/sql"
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/jmoiron/sqlx"
)

// GetLastMigration queries the migration table to find the current db version. if no results are found it defaults to V0.0.0
func GetLastMigration() Model {
	var lastMigration Model
	err := database.DB.QueryRowx(`SELECT * FROM migrations WHERE is_on_version=TRUE`).StructScan(&lastMigration)
	if err != nil {
		if err != sql.ErrNoRows {
			_, err = database.DB.Exec(`CREATE TABLE IF NOT EXISTS migrations(
				migration_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
				version VARCHAR(16),
				is_on_version BOOLEAN)`)
			if err != nil {
				log.Panic(err)
			}
		}
		return Model{0, "v0.0.0", true}
	}
	return lastMigration
}

// UpdateDB will either roll back or build the db to a specified version
func UpdateDB(lastMigration Model, isBuild bool, constraint *semver.Constraints) error {

	semanticVersions := getAvailableVersions()

	for i, version := range semanticVersions {
		if isBuild {
			if constraint.Check(version) {
				err := runMigration(version.Original(), isBuild)
				if err != nil {
					return err
				}
			}
		} else {
			reverseOrder := len(semanticVersions) - 1 - i
			if constraint.Check(semanticVersions[reverseOrder]) {
				err := runMigration(semanticVersions[reverseOrder].Original(), isBuild)
				if err != nil {
					return err
				}
			}
		}
	}
	return updateMigrationTable(isBuild, lastMigration)
}

// MakeConstraint makes a constraint for building or rolling back to a taret db version
func MakeConstraint(lastMigration Model, isBuild bool, targetVersion string) *semver.Constraints {
	var constraint *semver.Constraints
	var err error
	if isBuild {
		constraint, err = semver.NewConstraint("> " + lastMigration.Version + ", <= " + targetVersion)
	} else {
		constraint, err = semver.NewConstraint("< " + lastMigration.Version + ", >= " + targetVersion)
	}
	if err != nil {
		log.Panic(err)
	}

	return constraint
}

func updateMigrationTable(isBuild bool, lastMigration Model) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			if err != nil {
				log.Println(err)
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	for i, version := range Versions {
		var isCurrentVersion bool
		if isBuild {
			isCurrentVersion = i == len(Versions)-1
		} else {
			isCurrentVersion = i == 0
		}

		var migration Model
		migrationRow := tx.QueryRowx(`SELECT * FROM migrations WHERE version=$1`, version)

		err = migrationRow.StructScan(&migration)
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO migrations VALUES(DEFAULT, $1, $2)`, version, isCurrentVersion)
			if err != nil {
				return err
			}
		} else if migration.IsOnVersion != isCurrentVersion {
			_, err = tx.Exec(`UPDATE migrations SET is_on_version=$1 WHERE version='$1'`, isCurrentVersion, version)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func runMigration(version string, isBuild bool) error {
	switch version {
	case "v0.0.1":
		err := runInTransaction(isBuild, MigrateV0_0_1)
		if err != nil {
			return err
		}
	case "v1.0.0":
		err := runInTransaction(isBuild, MigrateV1_0_0)
		if err != nil {
			return err
		}
	}
	return nil
}

func runInTransaction(isBuild bool, scriptToRun func(bool, *sqlx.Tx) error) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Panic(err)
	}
	err = scriptToRun(isBuild, tx)

	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func getAvailableVersions() []*semver.Version {
	semanticVersions := make([]*semver.Version, len(Versions))
	for i, r := range Versions {
		semanticVersion, err := semver.NewVersion(r)
		if err != nil {
			log.Panic(err)
		}

		semanticVersions[i] = semanticVersion
	}
	return semanticVersions
}
