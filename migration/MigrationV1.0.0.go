package migration

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func buildV1_0_0(tx *sqlx.Tx) {
	_, err := tx.Exec(`ALTER TABLE request_data
		ADD COLUMN school_state VARCHAR(512),
		ADD COLUMN three_year_repayment_declining_balance_2016 VARCHAR(512)
		`)
	if err != nil {
		log.Fatalln(err)
	}
}

func rollbackV1_0_0(tx *sqlx.Tx) {
	_, err := tx.Exec(`ALTER TABLE request_data
		DROP COLUMN school_state,
		DROP COLUMN three_year_repayment_declining_alance_2016
		`)
	if err != nil {
		log.Fatalln(err)
	}

}

//MigrateV1_0_0 either builds or rollsback db versions to v1.0.0
func MigrateV1_0_0(build bool, tx *sqlx.Tx) {
	if build {
		buildV1_0_0(tx)
	} else {
		rollbackV1_0_0(tx)
	}
}
