package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filters := make(map[string]string)
	filters["page"] = "0"
	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	requestURL := getRequestURL(baseURL, filters)

	response := requestData(requestURL, httpClient)

	fmt.Println(response.TextOutput())
}

func getRequestURL(baseURL string, filters map[string]string) *url.URL {
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalln(err)
	}

	query := requestURL.Query()
	for key, value := range filters {
		query.Set(key, value)
	}
	requestURL.RawQuery = query.Encode()

	return requestURL
}

func requestData(url *url.URL, httpClient *http.Client) CollegeScoreCardResponseDTO {
	request, _ := http.NewRequest(http.MethodGet, url.String(), nil)

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	var parsedResponse CollegeScoreCardResponseDTO
	json.NewDecoder(resp.Body).Decode(&parsedResponse)

	return parsedResponse
}
