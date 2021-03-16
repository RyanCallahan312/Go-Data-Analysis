package dto

import (
	"fmt"
)

// JobDataDTO holds the job data from a 2019 massachusetts job/wage data spreadsheet
type JobDataDTO struct {
	State                      string  `json:"state" db:"state"`
	OccupationMajorTitle       string  `json:"occupationMajorTitle" db:"occupation_major_title"`
	TotalEmployment            int     `json:"totalEmployment" db:"total_employment"`
	PercentileSalary25thHourly float32 `json:"percentileSalary25thHourly" db:"percentile_salary_25th_hourly"`
	PercentileSalary25thAnnual int     `json:"percentileSalary25thAnnual" db:"percentile_salary_25th_annual"`
	OccupationCode             string  `json:"occupationCode" db:"occupation_code"`
}

// TextOutput job data to string
func (jobDataDto JobDataDTO) TextOutput() string {
	return fmt.Sprintf("State: %s,\n\tOccupationMajorTitle: %s,\n\tTotalEmployment: %d,\n\tPercentileSalary25thHourly: %f,\n\tPercentileSalary25thAnual: %d,\n\tOccupationCode: %s\n\n",
		jobDataDto.State,
		jobDataDto.OccupationMajorTitle,
		jobDataDto.TotalEmployment,
		jobDataDto.PercentileSalary25thHourly,
		jobDataDto.PercentileSalary25thAnnual,
		jobDataDto.OccupationCode)
}
