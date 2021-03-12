package migration

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func buildV0_0_0(tx *sqlx.Tx) {

	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS metadata (
		metadata_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY, 
		total_results INTEGER, 
		page_number INTEGER, 
		per_page INTEGER)`)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS request (
		request_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY, 
		metadata_id INTEGER,
		CONSTRAINT fk_request_data
			FOREIGN KEY(metadata_id)
			REFERENCES metadata(metadata_id))`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS request_data (
		request_data_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		request_id INTEGER,
		data_id INTEGER,
		school_name VARCHAR(512), 
		school_city VARCHAR(512), 
		student_size_2018 INTEGER, 
		student_size_2017 INTEGER, 
		over_poverty_three_years_after_completetion_2017 INTEGER, 
		three_year_repayment_overall_2016 INTEGER,
		CONSTRAINT fk_request_data
			FOREIGN KEY(request_id)
			REFERENCES request(request_id))`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS state_employment_data (
		state_employment_data_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		state VARCHAR(512),
		occupation_major_title VARCHAR(512),
		total_employment INTEGER, 
		percentile_salary_25th_hourly REAL,
		percentile_salary_25th_annual INTEGER,
		occupation_code VARCHAR(512))`)
	if err != nil {
		log.Fatalln(err)
	}
}

func rollbackV0_0_0(tx *sqlx.Tx) {

	_, err := tx.Exec(`DROP TABLE IF EXISTS state_employment_data`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS request_data CASCADE`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS request CASCADE`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS metadata CASCADE`)
	if err != nil {
		log.Fatalln(err)
	}

}

// MigrateV0_0_0 either builds or rollsback db versions to v0.0.0
func MigrateV0_0_0(build bool, tx *sqlx.Tx) {

	if build {
		buildV0_0_0(tx)
	} else {
		rollbackV0_0_0(tx)
	}
}
