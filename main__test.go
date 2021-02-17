package main

import (
	"Project1/vendor/github.com/joho/godotenv"
	"log"
	"testing"
)

func TestRequestData(t *testing.T) {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatalln(err)
	}

	url := "https://api.data.gov/ed/collegescorecard/v1/schools.json?fields=id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall&school.degrees_awarded.predominant=2,3&page=160&per_page=20&api_key="
	url = url + env["API_KEY"]
}
