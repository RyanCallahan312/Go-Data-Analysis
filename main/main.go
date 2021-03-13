package main

import (
	"Project1/config"
	"Project1/database"
	"Project1/dto"
	"Project1/migration"
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

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, err := os.Stat(config.ProjectRootPath + "/.env")
	if err == nil {
		err := godotenv.Load(config.ProjectRootPath + "/.env")
		if err != nil {
			log.Fatalln(err)
		}
	}

	migration.InitalizeDB(os.Getenv("DATABASE_NAME"))
	database.DB, err = sqlx.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}
	migration.MigrateToLatest()
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	errFile, err := os.Create("./err.txt")
	if err != nil {
		log.Fatalln(err)
	}
	errWriter := bufio.NewWriter(errFile)

	getAPIData(errWriter)

	getSheetData()

}

func getSheetData() {
	rows := getSheetRows("state_M2019_dl.xlsx", "State_M2019_dl")
	jobDataDTOs := getJobDataDTOs(rows)
	for _, jobDataDTO := range jobDataDTOs {
		writeJobDataToDb(jobDataDTO)
	}
}

func getJobDataDTOs(rows [][]string) []dto.JobDataDTO {

	jobDataDTOs := make([]dto.JobDataDTO, 0)
	for _, row := range rows {
		if row[9] == "major" && row[2] == "2" && row[1] != "District of Columbia" {
			jobDataDTOs = append(jobDataDTOs, getJobDataDTO(row))
		}
	}

	return jobDataDTOs
}

func getSheetRows(fileName string, sheetName string) [][]string {
	sheet, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := sheet.GetRows(sheetName)
	if err != nil {
		log.Fatalln(err)
	}

	return rows
}

func getJobDataDTO(row []string) dto.JobDataDTO {
	totalEmployemnt, err := strconv.ParseInt(row[10], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	percentileSalary25thHourly, err := strconv.ParseFloat(row[19], 32)
	if err != nil {
		log.Fatalln(err)
	}

	percentileSalary25thAnnual, err := strconv.ParseInt(row[24], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	jobDataDTO := dto.JobDataDTO{
		State:                      row[1],
		OccupationMajorTitle:       row[8],
		TotalEmployment:            int(totalEmployemnt),
		PercentileSalary25thHourly: float32(percentileSalary25thHourly),
		PercentileSalary25thAnnual: int(percentileSalary25thAnnual),
		OccupationCode:             row[7],
	}

	return jobDataDTO

}

func writeJobDataToDb(data dto.JobDataDTO) {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	_, err = tx.Exec(`INSERT INTO state_employment_data VALUES(DEFAULT, $1, $2, $3, $4, $5, $6)`,
		data.State,
		data.OccupationMajorTitle,
		data.TotalEmployment,
		data.PercentileSalary25thHourly,
		data.PercentileSalary25thAnnual,
		data.OccupationCode)
	if err != nil {
		log.Fatalln(err)
	}

}

func getAPIData(writer *bufio.Writer) {

	httpClient := &http.Client{}
	baseURL := "https://api.data.gov/ed/collegescorecard/v1/schools.json"

	filterBase := createFilterBase()

	page := 0

	var fileLock sync.Mutex

	// make first request to get how many pages we need to retrieve
	requestURL := getRequestURL(baseURL, filterBase)
	response, rawResponse := requestData(requestURL, httpClient)

	if rawResponse != "" {
		writeToFile(rawResponse, writer, &fileLock)
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

func createFilterBase() map[string]string {

	filters := make(map[string]string)

	//using this to get around rate limiting
	filters["per_page"] = "100"

	filters["school.degrees_awarded.predominant"] = "2,3"
	filters["fields"] = "id,school.name,school.city,2018.student.size,2017.student.size,2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line,2016.repayment.3_yr_repayment.overall"
	filters["api_key"] = os.Getenv("API_KEY")

	return filters
}

func initalizeTables() {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}
	}()

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

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS state_employment_data (
		state_employment_data_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		state VARCHAR(512),
		occupation_major_title VARCHAR(512),
		total_employment INTEGER, 
		percentile_salary_25th_hourly REAL,
		percentile_salary_25th_annual INTEGER,
		occupation_code VARCHAR(512))`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS migration (
		migration_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		version VARCHAR(16),
		is_on_version BOOLEAN)`)
	if err != nil {
		log.Fatalln(err)
	}

	var lastMigration migration.MigrationModel
	err = tx.QueryRowx(`SELECT * FROM migration WHERE is_on_version`).StructScan(&lastMigration)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Fatalln(err)
		}
	}
}

func writeCollegeScoreCardDataToDb(data dto.CollegeScoreCardResponseDTO) {
	tx, err := database.DB.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	metadata := data.Metadata
	lastInsertID := 0
	_ = tx.QueryRow("INSERT INTO metadata VALUES (DEFAULT, $1, $2, $3) RETURNING metadata_id", metadata.TotalResults, metadata.PageNumber, metadata.ResultsPerPage).Scan(&lastInsertID)

	metadataID := lastInsertID
	if err != nil {
		log.Fatalln(err)
	}

	_ = tx.QueryRow(`INSERT INTO request VALUES (DEFAULT, $1) RETURNING request_id`, metadataID).Scan(&lastInsertID)

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
