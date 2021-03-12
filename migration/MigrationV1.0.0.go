package migration

import (
	"github.com/jmoiron/sqlx"
)

func buildV1_0_0(tx *sqlx.Tx) {

}

func rollbackV1_0_0(tx *sqlx.Tx) {

}

// MigrateV1_0_0 either builds or rollsback db versions to v1.0.0
func MigrateV1_0_0(build bool, tx *sqlx.Tx) {

	if build {
		buildV0_0_0(tx)
	} else {
		rollbackV0_0_0(tx)
	}
}
