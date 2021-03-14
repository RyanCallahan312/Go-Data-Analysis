package spreadsheet

import (
	"Project1/config"
	"Project1/database"
	"Project1/dto"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func GetSheetData() {
	rows := getSheetRows(config.ProjectRootPath+"/state_M2019_dl.xlsx", "State_M2019_dl")
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

func writeJobDataToDb(data dto.JobDataDTO) {
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

	_, err = tx.Exec(`INSERT INTO state_employment_data VALUES(DEFAULT, $1, $2, $3, $4, $5, $6)`,
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
