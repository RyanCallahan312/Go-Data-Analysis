package migration

import (
	"Project1/database"
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/jmoiron/sqlx"
)

var versions = []string{
	"V0.0.0",
}

// UpdateDBFromVersion go from a given version of a db to the latest version
func UpdateDBFromVersion(version string) {

	constraint, err := semver.NewConstraint("> " + version)
	if err != nil {
		log.Fatalln(err)
	}

	sematicVersions := getAvailableVersions()

	for _, version := range sematicVersions {
		if constraint.Check(version) {

		}
	}

}

func findFunction(version string, build bool) {
	switch version {
	case "V0.0.0":
		runInTransaction(build, MigrateV0_0_0)
	case "V1.0.0":
		runInTransaction(build, MigrateV1_0_0)
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
	semanticVersions := make([]*semver.Version, len(versions))
	for i, r := range versions {
		semanticVersion, err := semver.NewVersion(r)
		if err != nil {
			log.Fatalln(err)
		}

		semanticVersions[i] = semanticVersion
	}
	return semanticVersions
}
