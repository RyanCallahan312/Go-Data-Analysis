package main

import (
	"Project1/vendor/github.com/joho/godotenv"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestRequestData(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filterBase := createFilterBase()

	page := 0

	// make first request to get how many pages we need to retrieve
	requestURL := getRequestURL(baseURL, filterBase)
	response, rawResponse := requestData(requestURL, httpClient)
	if rawResponse != "" {
		t.Errorf("Response != 200")
	}

	var responses []CollegeScoreCardResponseDTO
	responses = append(responses, response)

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
				t.Errorf("Response != 200")
			}

			var responses []CollegeScoreCardResponseDTO
			responses = append(responses, response)

		}(page)

	}
	wg.Wait()

	totalData := 0
	for _, response := range responses {
		totalData += len(response.Results)
	}

	if totalData < 1000 {
		t.Errorf("Did not retreive enough data")
	}

}

func TestWriteToDb(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	var dbExists bool
	_ = conn.QueryRow(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = 'comp490project1TEST'
			)`).Scan(&dbExists)

	if !dbExists {
		_, err = conn.Exec(`CREATE DATABASE comp490project1TEST`)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err = sql.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING_TEST"))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	initalizeTables(conn)
	response := CollegeScoreCardResponseDTO{CollegeScoreCardMetadataDTO{10, 11, 12}, []CollegeScoreCardFieldsDTO{{1, "bsu", "bridgew", 1, 2, 3, 4}, {2, "bsu", "bridgew", 1, 2, 3, 4}}}
	writeToDb(response, conn)
}
