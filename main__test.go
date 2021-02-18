package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
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
	}

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

	testResponse := CollegeScoreCardResponseDTO{
		CollegeScoreCardMetadataDTO{10, 11, 12},
		[]CollegeScoreCardFieldsDTO{
			{1, "bsu", "bridgew", 1, 2, 3, 4},
			{2, "bsu", "bridgew", 1, 2, 3, 4}}}

	writeToDb(testResponse, conn)

	idRows, err := conn.Query(`SELECT DISTINCT request_id FROM request`)
	if err != nil {
		log.Fatalln(err)
	}

	sqlxConn, err := sqlx.Open("pgx", os.Getenv("TEST_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}

	scoreCards := make([]CollegeScoreCardResponseDTO, 0)
	for idRows.Next() {
		var requestDataID int
		err := idRows.Scan(&requestDataID)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(requestDataID)

		var metadata CollegeScoreCardMetadataDTO
		metadataRow := sqlxConn.QueryRowx(`SELECT total_results, page_number, per_page FROM metadata WHERE metadata_id = $1`, requestDataID)

		err = metadataRow.StructScan(&metadata)
		if err != nil {
			log.Fatalln(err)
		}

		results := make([]CollegeScoreCardFieldsDTO, 0)
		var result CollegeScoreCardFieldsDTO
		dataRows, err := sqlxConn.Queryx(`SELECT data_id, school_name, school_name, school_city, student_size_2018, student_size_2017, over_poverty_three_years_after_completetion_2017, three_year_repayment_overall_2016 FROM request_data WHERE request_id = $1`, requestDataID)
		if err != nil {
			log.Fatalln(err)
		}

		for dataRows.Next() {
			dataRows.StructScan(&result)
			results = append(results, result)
		}

		scoreCards = append(scoreCards, CollegeScoreCardResponseDTO{metadata, results})

	}
	err = sqlxConn.Close()
	if err != nil {
		log.Fatalln(err)
	}

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

	if !reflect.DeepEqual(scoreCards[0], testResponse) {
		t.Errorf("Inserted data does not equal queried data")
	}
}
