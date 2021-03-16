package analysis

import (
	"Project1/dto"
	"sort"
)

func CollegeGradsToAmountOfJobs(scorecardData []dto.CollegeScoreCardFieldsDTO, jobData []dto.JobDataDTO) []interface{} {

	// interface is [int, int]
	stateToData := make(map[string][2]int)

	for _, scorecardFields := range scorecardData {
		val, exists := stateToData[scorecardFields.SchoolState]
		if !exists {
			val = [2]int{scorecardFields.StudentSize2018, 0}
		} else {
			val[0] += scorecardFields.StudentSize2018
		}

		stateToData[scorecardFields.SchoolState] = val
	}

	for _, jobFields := range jobData {
		if jobFields.OccupationCode[0:1] == "30-39" || jobFields.OccupationCode[0:1] == "40-49" {
			continue
		}

		val, exists := stateToData[jobFields.State]
		if !exists {
			val = [2]int{0, jobFields.TotalEmployment}
		} else {
			val[1] += jobFields.TotalEmployment
		}

		stateToData[jobFields.State] = val

	}

	var keys []string
	for k := range stateToData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]interface{}, len(keys))

	for i, key := range keys {
		val := stateToData[key]
		model := CollegeGradsToJobsModel{State: key, CollegeGrads: val[0], NumberOfJobs: val[1]}
		result[i] = model
	}

	return result

}
