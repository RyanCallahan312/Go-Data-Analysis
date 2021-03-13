package dto

import (
	"fmt"
	"strings"
)

// CollegeScoreCardResponseDTO is the college score card data and the metadata
type CollegeScoreCardResponseDTO struct {
	Metadata CollegeScoreCardMetadataDTO `json:"metadata"`
	Results  []CollegeScoreCardFieldsDTO `json:"results"`
}

//TextOutput is exported,it formats the data to plain text.
func (resDTO CollegeScoreCardResponseDTO) TextOutput() string {
	var sb strings.Builder

	sb.WriteString("MetaData: \n")
	resDTO.Metadata.TextOutput(&sb)

	for i, result := range resDTO.Results {
		sb.WriteString(fmt.Sprintf("Result %d: \n", i+1))
		result.TextOutput(&sb)
	}

	return sb.String()
}

// CollegeScoreCardMetadataDTO holds the response metadata
type CollegeScoreCardMetadataDTO struct {
	TotalResults   int `json:"total" db:"total_results"`
	PageNumber     int `json:"page" db:"page_number"`
	ResultsPerPage int `json:"per_page" db:"per_page"`
}

//TextOutput is exported,it formats the data to plain text.
func (metadataDTO CollegeScoreCardMetadataDTO) TextOutput(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("\t%-30s %-10d\n", "Total Results:", metadataDTO.TotalResults))
	sb.WriteString(fmt.Sprintf("\t%-30s %-10d\n", "Page Number:", metadataDTO.PageNumber))
	sb.WriteString(fmt.Sprintf("\t%-30s %-10d\n", "Results Per Page:", metadataDTO.ResultsPerPage))
}

// CollegeScoreCardFieldsDTO holds the required fields from college score card response
type CollegeScoreCardFieldsDTO struct {
	ID                                                   int    `json:"id" db:"data_id"`
	SchoolName                                           string `json:"school.name" db:"school_name"`
	SchoolCity                                           string `json:"school.city" db:"school_city"`
	StudentSize2018                                      int    `json:"2018.student.size" db:"student_size_2018"`
	StudentSize2017                                      int    `json:"2017.student.size" db:"student_size_2017"`
	StudentsOverPovertyLineThreeYearsAfterCompletion2017 int    `json:"2017.earnings.3_yrs_after_completion.overall_count_over_poverty_line" db:"over_poverty_three_years_after_completetion_2017"`
	ThreeYearRepaymentOverall2016                        int    `json:"2016.repayment.3_yr_repayment.overall" db:"three_year_repayment_overall_2016"`
}

//TextOutput is exported,it formats the data to plain text.
func (fieldsDTO CollegeScoreCardFieldsDTO) TextOutput(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("\t%-70s %-30d\n", "ID:", fieldsDTO.ID))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30s\n", "School Name:", fieldsDTO.SchoolName))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30s\n", "School City:", fieldsDTO.SchoolCity))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30d\n", "2018 Student Size:", fieldsDTO.StudentSize2018))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30d\n", "2017 Student Size:", fieldsDTO.StudentSize2017))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30d\n", "2017 Students Over Poverty Line Three Years After Completetion:", fieldsDTO.StudentsOverPovertyLineThreeYearsAfterCompletion2017))
	sb.WriteString(fmt.Sprintf("\t%-70s %-30d\n", "2016 Three Year Repayment Overall:", fieldsDTO.ThreeYearRepaymentOverall2016))
}
