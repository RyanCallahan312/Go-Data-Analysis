package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filters := make(map[string]string)
	page := 0
	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	outFile, err := os.Create("./out.txt")
	if err != nil {
		log.Fatalln(err)
	}
	writer := bufio.NewWriter(outFile)

	response := makeRequest(baseURL, filters, httpClient)

	for (page * response.Metadata.ResultsPerPage) < (response.Metadata.TotalResults) {
		page++
		filters["page"] = strconv.Itoa(page)
		_, err = writer.WriteString(makeRequest(baseURL, filters, httpClient).TextOutput())
		if err != nil {
			log.Fatalln(err)
		}

	}

}

func makeRequest(baseURL string, filters map[string]string, httpClient *http.Client) CollegeScoreCardResponseDTO {

	requestURL := getRequestURL(baseURL, filters)

	response := requestData(requestURL, httpClient)

	return response
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
	defer resp.Body.Close()

	var parsedResponse CollegeScoreCardResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return parsedResponse
}
