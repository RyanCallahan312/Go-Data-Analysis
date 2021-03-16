package spreadsheet

import (
	"Project1/config"
	"Project1/database"
	"Project1/dto"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jmoiron/sqlx"
)

// UpdateSheetData gets sheet data and writes it to db
func UpdateSheetData(fileName string, sheetName string) {
	rows := getSheetRows(config.ProjectRootPath+"/"+fileName+".xlsx", sheetName)
	jobDataDTOs := getJobDataDTOs(rows)

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

	_, err = tx.Exec(`TRUNCATE state_employment_data CASCADE`)
	if err != nil {
		log.Panic(err)
	}

	_, err = tx.Exec(`SELECT setval(pg_get_serial_sequence('state_employment_data', 'state_employment_data_id'), coalesce(max(state_employment_data_id),0) + 1, false) FROM state_employment_data`)
	if err != nil {
		log.Panic(err)
	}

	for _, jobDataDTO := range jobDataDTOs {
		writeJobDataToDb(jobDataDTO, tx)
	}
}

// GetSheetData retrieves all spreadsheet data in db
func GetSheetData() []dto.JobDataDTO {
	rows, err := database.DB.Queryx(`SELECT state, occupation_major_title, total_employment, percentile_salary_25th_hourly, percentile_salary_25th_annual, occupation_code FROM state_employment_data`)
	if err != nil {
		log.Panic(err)
	}

	allJobData := make([]dto.JobDataDTO, 0)
	var jobData dto.JobDataDTO
	for rows.Next() {
		err = rows.StructScan(&jobData)
		if err != nil {
			log.Panic(err)
		}

		allJobData = append(allJobData, jobData)
	}

	return allJobData
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
		log.Panic(err)
	}

	rows, err := sheet.GetRows(sheetName)
	if err != nil {
		log.Panic(err)
	}

	return rows
}

func getJobDataDTO(row []string) dto.JobDataDTO {
	totalEmployemnt, err := strconv.ParseInt(row[10], 10, 32)
	if err != nil {
		log.Panic(err)
	}

	percentileSalary25thHourly, err := strconv.ParseFloat(row[19], 32)
	if err != nil {
		log.Panic(err)
	}

	percentileSalary25thAnnual, err := strconv.ParseInt(row[24], 10, 32)
	if err != nil {
		log.Panic(err)
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

func writeJobDataToDb(data dto.JobDataDTO, tx *sqlx.Tx) {

	_, err := tx.Exec(`INSERT INTO state_employment_data VALUES(DEFAULT, $1, $2, $3, $4, $5, $6)`,
		data.State,
		data.OccupationMajorTitle,
		data.TotalEmployment,
		data.PercentileSalary25thHourly,
		data.PercentileSalary25thAnnual,
		data.OccupationCode)
	if err != nil {
		log.Panic(err)
	}

}
