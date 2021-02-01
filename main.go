package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	httpClient := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "https://api.data.gov/ed/collegescorecard/v1/schools.json?page=0&school.degrees_awarded.predominant=2,3&fields=id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall&api_key=HK4nxy96BDiKgcWABCBVMqdctCnTWWpUcUwSbjjE", nil)

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	var parsedResponse CollegeScoreCardResponseDTO
	json.NewDecoder(resp.Body).Decode(&parsedResponse)

	fmt.Println(parsedResponse.TextOutput())
}
