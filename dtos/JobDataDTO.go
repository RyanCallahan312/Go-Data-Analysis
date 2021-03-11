package dtos

import (
	"fmt"
)

// JobDataDTO holds the job data from a 2019 massachusetts job/wage data spreadsheet
type JobDataDTO struct {
	State                      string
	OccupationMajorTitle       string
	TotalEmployment            int
	PercentileSalary25thHourly float32
	PercentileSalary25thAnnual int
	OccupationCode             string
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
