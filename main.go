package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

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
	var fileLock sync.Mutex

	orderChannel := make(chan int, 1)
	orderChannel <- 0
	response := writeRequestToFile(baseURL, filters, page, httpClient, writer, &fileLock, orderChannel)

	wg := sync.WaitGroup{}
	for (page)*response.Metadata.ResultsPerPage < response.Metadata.TotalResults {
		page++

		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			flt := make(map[string]string)
			for k, v := range filters {
				flt[k] = v
			}
			flt["page"] = strconv.Itoa(page)

			_ = writeRequestToFile(baseURL, flt, page, httpClient, writer, &fileLock, orderChannel)

		}(page)

	}

	wg.Wait()

}

func shouldWrite(requestNumber int, ch chan int) bool {
	message := <-ch

	if requestNumber != message {
		ch <- message
	}

	return requestNumber == message
}

func writeRequestToFile(baseURL string, filters map[string]string, page int, httpClient *http.Client, writer *bufio.Writer, fileLock *sync.Mutex, orderChannel chan int) CollegeScoreCardResponseDTO {
	requestURL := getRequestURL(baseURL, filters)

	retry := true
	retryCount := 0

	var response CollegeScoreCardResponseDTO
	var rawResponse string

	for retry {
		response, rawResponse = requestData(requestURL, httpClient)
		if rawResponse == "" || retryCount > 3 {
			retry = false
		}
		retryCount++
	}

	for !shouldWrite(page, orderChannel) {

	}

	fileLock.Lock()
	if retry {
		_, err := writer.WriteString(rawResponse)
		if err != nil {
			log.Fatalln(err)
		}

	} else {
		_, err := writer.WriteString(response.TextOutput())
		if err != nil {
			log.Fatalln(err)
		}
	}
	fileLock.Unlock()
	orderChannel <- page + 1

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

func requestData(url *url.URL, httpClient *http.Client) (CollegeScoreCardResponseDTO, string) {
	request, _ := http.NewRequest(http.MethodGet, url.String(), nil)

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var parsedResponse CollegeScoreCardResponseDTO
	rawResponse := ""

	if resp.StatusCode == http.StatusOK {

		err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
		if err != nil {
			log.Fatalln(err)
		}

	} else {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		rawResponse = string(bodyBytes)
	}

	return parsedResponse, rawResponse
}
