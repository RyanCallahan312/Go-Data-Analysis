package analysis

import (
	"Project1/dto"
	"math"
	"sort"
)

func CollegeGradsToAmountOfJobs(scorecardData []dto.CollegeScoreCardFieldsDTO, jobData []dto.JobDataDTO) []interface{} {

	// interface is [int, int]
	stateToData := make(map[string][2]int)

	for _, scorecardFields := range scorecardData {
		if usc[scorecardFields.SchoolState] == "" {
			continue
		}

		val, exists := stateToData[usc[scorecardFields.SchoolState]]
		if !exists {
			val = [2]int{scorecardFields.StudentSize2018 / 4, 0}
		} else {
			val[0] += scorecardFields.StudentSize2018 / 4
		}

		stateToData[usc[scorecardFields.SchoolState]] = val
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
		ratio := float32(val[0]) / float32(val[1])
		model := CollegeGradsToJobsModel{State: key, CollegeGrads: val[0], NumberOfJobs: val[1], Ratio: float32(math.Floor(float64(ratio)*1000) / 1000)}
		result[i] = model
	}

	return result

}

func DecliningBalanceToSalary(scorecardData []dto.CollegeScoreCardFieldsDTO, jobData []dto.JobDataDTO) []interface{} {
	// interface is [int, int]
	stateToData := make(map[string][2]interface{})

	for _, scorecardFields := range scorecardData {
		if usc[scorecardFields.SchoolState] == "" {
			continue
		}

		val, exists := stateToData[usc[scorecardFields.SchoolState]]
		if !exists {
			val = [2]interface{}{scorecardFields.ThreeYearRepaymentDecliningBalance2016, 0}
		} else {
			val[0] = val[0].(float32) + scorecardFields.ThreeYearRepaymentDecliningBalance2016
		}

		stateToData[usc[scorecardFields.SchoolState]] = val
	}

	for _, jobFields := range jobData {

		val, exists := stateToData[jobFields.State]
		if !exists {
			val = [2]interface{}{float32(0), jobFields.PercentileSalary25thAnnual}
		} else {
			val[1] = val[1].(int) + jobFields.PercentileSalary25thAnnual
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
		ratio := float32(val[1].(int)) / val[0].(float32)
		model := DecliningBalToSalarysModel{State: key, DecliningBalance: val[0].(float32), Salary25Percent: val[1].(int), Ratio: float32(math.Floor(float64(ratio)*1000) / 1000)}
		result[i] = model
	}

	return result
}

// courtacy of tmaiaroto on github https://gist.github.com/tmaiaroto/4ec7668ae986335b0a6d
var usc = map[string]string{
	"AL": "Alabama",
	"AK": "Alaska",
	"AZ": "Arizona",
	"AR": "Arkansas",
	"CA": "California",
	"CO": "Colorado",
	"CT": "Connecticut",
	"DE": "Delaware",
	"FL": "Florida",
	"GA": "Georgia",
	"HI": "Hawaii",
	"ID": "Idaho",
	"IL": "Illinois",
	"IN": "Indiana",
	"IA": "Iowa",
	"KS": "Kansas",
	"KY": "Kentucky",
	"LA": "Louisiana",
	"ME": "Maine",
	"MD": "Maryland",
	"MA": "Massachusetts",
	"MI": "Michigan",
	"MN": "Minnesota",
	"MS": "Mississippi",
	"MO": "Missouri",
	"MT": "Montana",
	"NE": "Nebraska",
	"NV": "Nevada",
	"NH": "New Hampshire",
	"NJ": "New Jersey",
	"NM": "New Mexico",
	"NY": "New York",
	"NC": "North Carolina",
	"ND": "North Dakota",
	"OH": "Ohio",
	"OK": "Oklahoma",
	"OR": "Oregon",
	"PA": "Pennsylvania",
	"RI": "Rhode Island",
	"SC": "South Carolina",
	"SD": "South Dakota",
	"TN": "Tennessee",
	"TX": "Texas",
	"UT": "Utah",
	"VT": "Vermont",
	"VA": "Virginia",
	"WA": "Washington",
	"WV": "West Virginia",
	"WI": "Wisconsin",
	"WY": "Wyoming",
}
