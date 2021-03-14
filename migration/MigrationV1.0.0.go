package migration

import (
	"github.com/jmoiron/sqlx"
)

func buildV1_0_0(tx *sqlx.Tx) error {
	_, err := tx.Exec(`ALTER TABLE request_data
		ADD COLUMN school_state VARCHAR(512),
		ADD COLUMN three_year_repayment_declining_balance_2016 FLOAT
		`)
	if err != nil {
		return err
	}

	return nil
}

func rollbackV1_0_0(tx *sqlx.Tx) error {
	_, err := tx.Exec(`ALTER TABLE request_data
		DROP COLUMN school_state,
		DROP COLUMN three_year_repayment_declining_alance_2016
		`)
	if err != nil {
		return err
	}
	return nil
}

//MigrateV1_0_0 either builds or rollsback db versions to v1.0.0
func MigrateV1_0_0(build bool, tx *sqlx.Tx) error {
	if build {
		return buildV1_0_0(tx)
	} else {
		return rollbackV1_0_0(tx)
	}
}
