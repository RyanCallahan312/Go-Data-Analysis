package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

func TestRequestData(t *testing.T) {
	if _, err := os.Stat("./.env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
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

	responses := make([]CollegeScoreCardResponseDTO, 0)
	responses = append(responses, response)
	sliceLock := &sync.Mutex{}

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
			sliceLock.Lock()
			responses = append(responses, response)
			sliceLock.Unlock()

		}(page)

	}
	wg.Wait()

	totalData := 0
	for _, val := range responses {
		totalData += len(val.Results)
		fmt.Println(len(val.Results))
	}
	fmt.Println(len(responses))

	if totalData < 1000 {
		t.Errorf("Did not retreive enough data; got %d", totalData)
	}

}

func TestWriteToDb(t *testing.T) {
	if _, err := os.Stat("./.env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	conn, err := sql.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	var dbExists bool
	_ = conn.QueryRow(`SELECT EXISTS (
			SELECT FROM pg_database 
			WHERE datname = 'comp490project1test'
			)`).Scan(&dbExists)

	if !dbExists {
		_, err = conn.Exec(`CREATE DATABASE comp490project1test`)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err = sql.Open("pgx", os.Getenv("TEST_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	initalizeTables(conn)

	response := CollegeScoreCardResponseDTO{
		CollegeScoreCardMetadataDTO{10, 11, 12},
		[]CollegeScoreCardFieldsDTO{
			{1, "bsu", "bridgew", 1, 2, 3, 4},
			{2, "bsu", "bridgew", 1, 2, 3, 4}}}

	writeToDb(response, conn)

	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err = sql.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Exec(`DROP DATABASE comp490project1test`)
	if err != nil {
		log.Fatalln(err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
