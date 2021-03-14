package api

import (
	"Project1/database"
	"Project1/dto"
	"Project1/main"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// GetAPIData gets the api data and writes it to the DB
func GetAPIData(writer *bufio.Writer) {

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filterBase := createFilterBase()

	page := 0

	var fileLock sync.Mutex

	// make first request to get how many pages we need to retrieve
	requestURL := getRequestURL(baseURL, filterBase)
	response, rawResponse := requestData(requestURL, httpClient)

	if rawResponse != "" {
		main.WriteToFile(rawResponse, writer, &fileLock)
	} else {
		writeCollegeScoreCardDataToDb(response)
	}

	wg := sync.WaitGroup{}
	for (page+1)*response.Metadata.ResultsPerPage <= response.Metadata.TotalResults {
		page++

		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			filters := make(map[string]string)
			for k, v := range filterBase {
				filters[k] = v
			}
			filters["page"] = strconv.Itoa(page)

			requestURL := getRequestURL(baseURL, filters)
			response, rawResponse := requestData(requestURL, httpClient)

			if rawResponse != "" {
				writeToFile(rawResponse, writer, &fileLock)
			} else {
				writeCollegeScoreCardDataToDb(response)
			}

		}(page)

	}
	wg.Wait()
}

func getRequestURL(baseURL string, filters map[string]string) *url.URL {
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		log.Panic(err)
	}

	query := requestURL.Query()
	for key, value := range filters {
		query.Set(key, value)
	}
	requestURL.RawQuery = query.Encode()

	return requestURL
}

func requestData(url *url.URL, httpClient *http.Client) (dto.CollegeScoreCardResponseDTO, string) {
	request, _ := http.NewRequest(http.MethodGet, url.String(), nil)

	retry := true
	retryCount := 0

	var parsedResponse dto.CollegeScoreCardResponseDTO
	rawResponse := ""

	var resp *http.Response
	var err error

	//retry count at 10 seems to work to get around rate limiting
	for retry && retryCount < 10 {

		resp, err = httpClient.Do(request)
		if err != nil {
			log.Panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			retryCount++

			//spacing out the retry requests randomly because if two goroutines retry at the same time it might cause it to rate limit again
			time.Sleep(time.Duration(200+rand.Int31n(300)) * time.Millisecond)
		} else {
			retry = false
		}
	}

	if resp.StatusCode == http.StatusOK {

		err := json.NewDecoder(resp.Body).Decode(&parsedResponse)
		if err != nil {
			log.Panic(err)
		}

	} else {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		rawResponse = string(bodyBytes)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	return parsedResponse, rawResponse
}

func writeCollegeScoreCardDataToDb(data dto.CollegeScoreCardResponseDTO) {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Panic(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Panic(err)
		}
	}()

	metadata := data.Metadata
	lastInsertID := 0
	_ = tx.QueryRow("INSERT INTO metadata VALUES (DEFAULT, $1, $2, $3) RETURNING metadata_id", metadata.TotalResults, metadata.PageNumber, metadata.ResultsPerPage).Scan(&lastInsertID)

	metadataID := lastInsertID
	if err != nil {
		log.Panic(err)
	}

	_ = tx.QueryRow(`INSERT INTO request VALUES (DEFAULT, $1) RETURNING request_id`, metadataID).Scan(&lastInsertID)

	requestID := lastInsertID
	if err != nil {
		log.Panic(err)
	}

	results := data.Results
	for _, requestData := range results {
		_, err = tx.Exec(`INSERT INTO request_data 
			(request_data_id, request_id, data_id, school_name, school_city, school_state, student_size_2018, student_size_2017, over_poverty_three_years_after_completetion_2017, three_year_repayment_overall_2016, three_year_repayment_declining_balance_2016) 
			VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			requestID, requestData.ID, requestData.SchoolName, requestData.SchoolCity, requestData.SchoolState, requestData.StudentSize2018, requestData.StudentSize2017, requestData.StudentsOverPovertyLineThreeYearsAfterCompletion2017, requestData.ThreeYearRepaymentOverall2016, requestData.ThreeYearRepaymentDecliningBalance2016)
		if err != nil {
			log.Panic(err)
		}
	}
}

func createFilterBase() map[string]string {

	filters := make(map[string]string)

	//using this to get around rate limiting
	filters["per_page"] = "100"

	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,school.state,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall,2016.repayment.repayment_cohort.3_year_declining_balance"
	filters["api_key"] = os.Getenv("API_KEY")

	return filters
}
