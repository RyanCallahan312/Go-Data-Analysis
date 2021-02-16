package main

import (
	"bufio"
	"database/sql"
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

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open("pgx", os.Getenv("CONNECTION_STRING"))
	defer conn.Close()

	initalizeTables(conn)

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filters := make(map[string]string)
	page := 0

	//using this to get around rate limiting
	filters["per_page"] = "100"

	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	outFile, err := os.Create("./err.txt")
	if err != nil {
		log.Fatalln(err)
	}

	writer := bufio.NewWriter(outFile)
	var fileLock sync.Mutex

	orderChannel := make(chan int, 1)

	// make first request to get how many pages we need to retrieve
	requestURL := getRequestURL(baseURL, filters)
	response, rawResponse := requestData(requestURL, httpClient)
	if rawResponse != "" {
		writeToFile(rawResponse, writer, &fileLock)
	} else {
		writeToDb(response, conn)
	}

	orderChannel <- page + 1

	wg := sync.WaitGroup{}
	for (page+1)*response.Metadata.ResultsPerPage <= response.Metadata.TotalResults {
		page++

		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			flt := make(map[string]string)
			for k, v := range filters {
				flt[k] = v
			}
			flt["page"] = strconv.Itoa(page)

			requestURL := getRequestURL(baseURL, flt)

			response, rawResponse := requestData(requestURL, httpClient)

			for !shouldWrite(page, orderChannel) {

			}

			if rawResponse != "" {
				writeToFile(rawResponse, writer, &fileLock)
			} else {
				writeToDb(response, conn)
			}

			orderChannel <- page + 1

		}(page)

	}
	wg.Wait()

}

func initalizeTables(conn *sql.DB) {
	tx, err := conn.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS metadata (
		metadata_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY, 
		total_results INTEGER, 
		page_number INTEGER, 
		per_page INTEGER)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS request (
		request_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY, 
		metadata_id INTEGER,
		CONSTRAINT fk_request_data
			FOREIGN KEY(metadata_id)
			REFERENCES metadata(metadata_id))`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS request_data (
		request_data_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		request_id INTEGER,
		data_id INTEGER,
		school_name VARCHAR(512), 
		school_city VARCHAR(512), 
		student_size_2018 INTEGER, 
		student_size_2017 INTEGER, 
		over_poverty_three_years_after_completetion_2017 INTEGER, 
		three_year_repayment_overall_2016 INTEGER,
		CONSTRAINT fk_request_data
			FOREIGN KEY(request_id)
			REFERENCES request(request_id))`)
	if err != nil {
		log.Fatalln(err)
	}

	tx.Commit()
}

func writeToDb(data CollegeScoreCardResponseDTO, conn *sql.DB) {
	tx, err := conn.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	metadata := data.Metadata
	lastInsertID := 0
	err = tx.QueryRow("INSERT INTO metadata VALUES (DEFAULT, $1, $2, $3) RETURNING metadata_id", metadata.TotalResults, metadata.PageNumber, metadata.ResultsPerPage).Scan(&lastInsertID)
	if err != nil {
		log.Fatalln(err)
	}
	metadataID := lastInsertID
	if err != nil {
		log.Fatalln(err)
	}

	err = tx.QueryRow(`INSERT INTO request VALUES (DEFAULT, $1) RETURNING request_id`, metadataID).Scan(&lastInsertID)
	if err != nil {
		log.Fatalln(err)
	}
	requestID := lastInsertID
	if err != nil {
		log.Fatalln(err)
	}

	results := data.Results
	for _, requestData := range results {
		_, err = tx.Exec("INSERT INTO request_data VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8)", requestID, requestData.ID, requestData.SchoolName, requestData.SchoolCity, requestData.StudentSize2018, requestData.StudentSize2017, requestData.StudentsOverPovertyLineThreeYearsAfterCompletion2017, requestData.ThreeYearRepaymentOverall2016)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func shouldWrite(requestNumber int, ch chan int) bool {
	message := <-ch

	if requestNumber != message {
		ch <- message
	}
	return requestNumber == message
}

func writeToFile(data string, writer *bufio.Writer, fileLock *sync.Mutex) {

	fileLock.Lock()
	_, err := writer.WriteString(data)
	if err != nil {
		log.Fatalln(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	fileLock.Unlock()
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

	retry := true
	retryCount := 0

	var parsedResponse CollegeScoreCardResponseDTO
	rawResponse := ""

	var resp *http.Response
	var err error

	//retry count at 10 seems to work to get around rate limiting
	for retry && retryCount < 10 {

		resp, err = httpClient.Do(request)
		if err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
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
		log.Fatalln(err)
	}

	return parsedResponse, rawResponse
}
