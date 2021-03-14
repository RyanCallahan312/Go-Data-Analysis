package main

import (
	"Project1/config"
	"Project1/database"
	"Project1/dto"
	"Project1/migration"
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

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, err := os.Stat(config.ProjectRootPath + "/.env")
	if err == nil {
		err := godotenv.Load(config.ProjectRootPath + "/.env")
		if err != nil {
			log.Panic(err)
		}
	}

	err = setUp()
	if err != nil {
		tearDown()
		log.Panic(err)
	}

	var retCode int
	defer func() {
		tearDown()
		os.Exit(retCode)
	}()

	retCode = m.Run()

}

func setUp() error {
	return buildTestDB()
}

func tearDown() {
	tearTestDownDB()
}

func TestGetSheetData(t *testing.T) {

	getSheetData()

	var stateCount int
	err := database.DB.QueryRow(`SELECT COUNT(DISTINCT state) FROM state_employment_data;`).Scan(&stateCount)
	if err != nil {
		t.Error(err)
	}

	if stateCount != 50 {
		t.Errorf("Expected 50 states; Got %d", stateCount)
	}

}

func TestGetJobDataDTOs(t *testing.T) {
	rows := getSheetRows("MajorAndOccCodeTest.xlsx", "Sheet1")
	jobDataDTOs := getJobDataDTOs(rows)

	if len(jobDataDTOs) != 4 {
		t.Errorf("Expected 4 rows; Got %d", len(jobDataDTOs))
	}

}

func TestGetJobDataDTO(t *testing.T) {
	row := []string{"01", "Alabama", "2", "000000", "Cross-industry", "cross-industry", "1235", "11-0000", "Management Occupations", "major", "83760", "1.2", "42.428", "0.77", "", "51.86", "107860", "0.6", "22.72", "31.80", "45.03", "63.07", "90.16", "47250", "66140", "93660", "131180", "187530"}

	jobDataDTO := getJobDataDTO(row)

	if !(jobDataDTO.State == "Alabama" &&
		jobDataDTO.OccupationMajorTitle == "Management Occupations" &&
		jobDataDTO.TotalEmployment == 83760 &&
		jobDataDTO.PercentileSalary25thHourly == 31.799999 &&
		jobDataDTO.PercentileSalary25thAnnual == 66140 &&
		jobDataDTO.OccupationCode == "11-0000") {
		t.Errorf("Mismatch Data: \n%s", jobDataDTO.TextOutput())
	}

}

func TestRequestData(t *testing.T) {

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

	responses := make([]dto.CollegeScoreCardResponseDTO, 0)
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

	testResponse := dto.CollegeScoreCardResponseDTO{
		Metadata: dto.CollegeScoreCardMetadataDTO{
			TotalResults:   10,
			PageNumber:     11,
			ResultsPerPage: 12,
		},
		Results: []dto.CollegeScoreCardFieldsDTO{
			{
				ID:              1,
				SchoolName:      "bsu",
				SchoolCity:      "bridgew",
				SchoolState:     "MA",
				StudentSize2018: 1,
				StudentSize2017: 2,
				StudentsOverPovertyLineThreeYearsAfterCompletion2017: 3,
				ThreeYearRepaymentOverall2016:                        4,
				ThreeYearRepaymentDecliningBalance2016:               4,
			},
			{
				ID:              2,
				SchoolName:      "bsu",
				SchoolCity:      "bridgew",
				SchoolState:     "MA",
				StudentSize2018: 1,
				StudentSize2017: 2,
				StudentsOverPovertyLineThreeYearsAfterCompletion2017: 3,
				ThreeYearRepaymentOverall2016:                        4,
				ThreeYearRepaymentDecliningBalance2016:               4,
			},
		},
	}

	writeCollegeScoreCardDataToDb(testResponse)

	idRows, err := database.DB.Query(`SELECT DISTINCT request_id FROM request`)
	if err != nil {
		t.Error(err)
	}

	scoreCards := make([]dto.CollegeScoreCardResponseDTO, 0)
	for idRows.Next() {
		var requestDataID int
		err := idRows.Scan(&requestDataID)
		if err != nil {
			t.Error(err)
		}

		var metadata dto.CollegeScoreCardMetadataDTO
		metadataRow := database.DB.QueryRowx(`SELECT total_results, page_number, per_page FROM metadata WHERE metadata_id = $1`, requestDataID)

		err = metadataRow.StructScan(&metadata)
		if err != nil {
			t.Error(err)
		}

		results := make([]dto.CollegeScoreCardFieldsDTO, 0)
		var result dto.CollegeScoreCardFieldsDTO
		dataRows, err := database.DB.Queryx(`SELECT data_id,
			school_name,
			school_city,
			school_state,
			student_size_2018,
			student_size_2017,
			over_poverty_three_years_after_completetion_2017,
			three_year_repayment_overall_2016,
			three_year_repayment_declining_balance_2016
			FROM request_data WHERE request_id = $1`,
			requestDataID)
		if err != nil {
			t.Error(err)
		}

		for dataRows.Next() {
			err = dataRows.StructScan(&result)
			if err != nil {
				t.Error(err)
			}

			results = append(results, result)
		}

		scoreCards = append(scoreCards, dto.CollegeScoreCardResponseDTO{Metadata: metadata, Results: results})

	}

	if !reflect.DeepEqual(scoreCards[0], testResponse) {
		log.Println(scoreCards[0].TextOutput())
		log.Println(testResponse.TextOutput())
		t.Errorf("Inserted data does not equal queried data")
	}
}

func buildTestDB() error {
	migration.InitalizeDB(os.Getenv("DATABASE_NAME") + "test")
	var err error
	database.DB, err = sqlx.Open("pgx", os.Getenv("TEST_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	return migration.MigrateToLatest()
}

func tearTestDownDB() {

	err := database.DB.Close()
	if err != nil {
		log.Panic(err)
	}

	database.DB, err = sqlx.Open("pgx", os.Getenv("MAINTENANCE_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	statement := fmt.Sprintf(`DROP DATABASE %stest`, os.Getenv("DATABASE_NAME"))

	_, err = database.DB.Exec(statement)
	if err != nil {
		log.Panic(err)
	}
}
