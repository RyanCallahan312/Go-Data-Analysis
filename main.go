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
	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	outFile, _ := os.Create("./out.txt")
	writer := bufio.NewWriter(outFile)

	order := make(chan int, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	response := makeRequest(baseURL, filters, page, httpClient, &wg, writer, order)

	order = make(chan int, response.Metadata.TotalResults/response.Metadata.ResultsPerPage)

	for ((page + 1) * response.Metadata.ResultsPerPage) < (response.Metadata.TotalResults / 4) {
		page++

		wg.Add(1)
		go makeRequest(baseURL, filters, page, httpClient, &wg, writer, order)

	}

	fmt.Println(order)

	wg.Wait()

}

func makeRequest(baseURL string, filters map[string]string, page int, httpClient *http.Client, wg *sync.WaitGroup, writer *bufio.Writer, order chan int) CollegeScoreCardResponseDTO {

	defer func() {
		order <- page
		wg.Done()
	}()

	filters["page"] = strconv.Itoa(page)

	requestURL := getRequestURL(baseURL, filters)

	response := requestData(requestURL, httpClient)

	if page > 1 {
		message := -1
		for message != page-1 {
			select {
			case message = <-order:
				if message == page-1 {

					writer.WriteString(response.TextOutput())
				}
			}
		}
	} else {
		writer.WriteString(response.TextOutput())
	}

	return response
}

func getRequestURL(baseURL string, filters map[string]string) *url.URL {
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalln(err)
	}

	query := requestURL.Query()
	for key, value := range filters {
		if key == "page" {
			fmt.Println(value)
		}
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
