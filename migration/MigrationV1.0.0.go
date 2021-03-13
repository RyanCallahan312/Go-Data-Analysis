package migration

import (
	"githb.com/jmoiron/sqlx"


fnc buildV1_0_0(tx *sqlx.Tx) {
tx.Exec(`ALTER TABLE request_data 
		ADD COLUMN school_state VARCHR(512),
	ADD COLUMN three_year_repayment_declining_balance_2016 VARCHAR(512)
	`)
}

fnc rollbackV1_0_0(tx *sqlx.Tx) {
tx.Exec(`ALTER TABLE request_data 
		DROP COLUMN school_state,
		DROP COLUMN three_year_repayment_declining_alance_2016
		`)

}

//MigrateV1_0_0 either builds or rollsback db versions to v1.0.0
fnc MigrateV1_0_0(build bool, tx *sqlx.Tx) {
	if build {
		buildV1_0_0(tx)
	} else {
		rollbackV1_0_0(tx)
	}
}
