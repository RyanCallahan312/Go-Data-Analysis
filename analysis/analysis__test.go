package analysis

import (
	"Project1/config"
	"Project1/dto"
	"log"
	"math"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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
func TestCollegeGradsToAmountOfJobs(t *testing.T) {

	inputCollegeData := make([]dto.CollegeScoreCardFieldsDTO, 0)
	inputCollegeData = append(inputCollegeData, dto.CollegeScoreCardFieldsDTO{
		ID:              1,
		SchoolName:      "bsu",
		SchoolCity:      "bridgew",
		SchoolState:     "MA",
		StudentSize2018: 4,
		StudentSize2017: 2,
		StudentsOverPovertyLineThreeYearsAfterCompletion2017: 3,
		ThreeYearRepaymentOverall2016:                        4,
		ThreeYearRepaymentDecliningBalance2016:               4.3,
	},
	)
	inputJobData := make([]dto.JobDataDTO, 0)
	inputJobData = append(inputJobData, dto.JobDataDTO{
		State:                      "Massachusetts",
		OccupationMajorTitle:       "Software Developer",
		TotalEmployment:            1000,
		PercentileSalary25thHourly: 35.75,
		PercentileSalary25thAnnual: 90000,
		OccupationCode:             "99-9999",
	})

	result := CollegeGradsToAmountOfJobs(inputCollegeData, inputJobData)

	expectedResult := float32(inputCollegeData[0].StudentSize2018/4) / float32(inputJobData[0].TotalEmployment)
	expectedResult = float32(math.Floor(float64(expectedResult)*1000) / 1000)

	if result[0].(CollegeGradsToJobsModel).Ratio != expectedResult {
		t.Errorf("Expected Ratio: %f; Got %f", expectedResult, result[0].(CollegeGradsToJobsModel).Ratio)
	}

}

func TestDecliningBalanceToSalary(t *testing.T) {

	inputCollegeData := make([]dto.CollegeScoreCardFieldsDTO, 0)
	inputCollegeData = append(inputCollegeData, dto.CollegeScoreCardFieldsDTO{
		ID:              1,
		SchoolName:      "bsu",
		SchoolCity:      "bridgew",
		SchoolState:     "MA",
		StudentSize2018: 4,
		StudentSize2017: 2,
		StudentsOverPovertyLineThreeYearsAfterCompletion2017: 3,
		ThreeYearRepaymentOverall2016:                        4,
		ThreeYearRepaymentDecliningBalance2016:               4.3,
	},
	)
	inputJobData := make([]dto.JobDataDTO, 0)
	inputJobData = append(inputJobData, dto.JobDataDTO{
		State:                      "Massachusetts",
		OccupationMajorTitle:       "Software Developer",
		TotalEmployment:            1000,
		PercentileSalary25thHourly: 35.75,
		PercentileSalary25thAnnual: 90000,
		OccupationCode:             "99-9999",
	})

	result := DecliningBalanceToSalary(inputCollegeData, inputJobData)

	expectedResult := float32(inputJobData[0].PercentileSalary25thAnnual) / float32(inputCollegeData[0].ThreeYearRepaymentDecliningBalance2016)
	expectedResult = float32(math.Floor(float64(expectedResult)*1000) / 1000)

	if result[0].(DecliningBalToSalarysModel).Ratio != expectedResult {
		t.Errorf("Expected Ratio: %f; Got %f", expectedResult, result[0].(DecliningBalToSalarysModel).Ratio)
	}

}

func setUp() error {
	config.InitEnv()
	return nil
}

func tearDown() {
}
