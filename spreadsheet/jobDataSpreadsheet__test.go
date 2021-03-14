package spreadsheet

import (
	"Project1/config"
	"Project1/database"
	"Project1/migration"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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

func TestGetJobDataDTOs(t *testing.T) {
	rows := getSheetRows(config.ProjectRootPath+"/MajorAndOccCodeTest.xlsx", "Sheet1")
	jobDataDTOs := getJobDataDTOs(rows)

	if len(jobDataDTOs) != 4 {
		t.Errorf("Expected 4 rows; Got %d", len(jobDataDTOs))
	}

}

func TestGetSheetData(t *testing.T) {

	GetSheetData()

	var stateCount int
	err := database.DB.QueryRow(`SELECT COUNT(DISTINCT state) FROM state_employment_data;`).Scan(&stateCount)
	if err != nil {
		t.Error(err)
	}

	if stateCount != 50 {
		t.Errorf("Expected 50 states; Got %d", stateCount)
	}

}

func setUp() error {
	return buildTestDB()
}

func tearDown() {
	tearTestDownDB()
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
