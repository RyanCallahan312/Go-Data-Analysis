package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filters := make(map[string]string)
	page := 0
	filters["page"] = strconv.Itoa(page)
	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	outFile, _ := os.Create("./out.txt")
	writer := bufio.NewWriter(outFile)

	wg := sync.WaitGroup{}
	wg.Add(1)
	response := makeRequest(baseURL, filters, httpClient, &wg, writer)

	for ((page + 1) * response.Metadata.ResultsPerPage) < response.Metadata.TotalResults {
		page++
		wg.Add(1)
		filters["page"] = strconv.Itoa(page)

		go makeRequest(baseURL, filters, httpClient, &wg, writer)

	}

	wg.Wait()

}

func makeRequest(baseURL string, filters map[string]string, httpClient *http.Client, wg *sync.WaitGroup, writer *bufio.Writer) CollegeScoreCardResponseDTO {

	defer wg.Done()

	requestURL := getRequestURL(baseURL, filters)

	response := requestData(requestURL, httpClient)

	writer.WriteString(response.TextOutput())

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
	fmt.Println(url.String())
	request, _ := http.NewRequest(http.MethodGet, url.String(), nil)

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	var parsedResponse CollegeScoreCardResponseDTO
	json.NewDecoder(resp.Body).Decode(&parsedResponse)

	return parsedResponse
}
