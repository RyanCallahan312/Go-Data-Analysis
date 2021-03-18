// +build wasm

package main

import (
	"Project1/analysis"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type TextVisualization struct { // no-lint
	app.Compo

	AnalysisType     int
	TitleBar         []string
	Data             []interface{}
	lastAnalysisType int
	filterField      string
	filter           string
	filterIndex      int
	filteredData     []interface{}
	sortDirection    string
	sortField        string
	sortIndex        int
}

func (page *TextVisualization) OnMount(ctx app.Context) {
	page.lastAnalysisType = page.AnalysisType
	page.filter = ""
	page.filterField = page.TitleBar[0]
	page.filteredData = make([]interface{}, len(page.Data))
	copy(page.filteredData, page.Data)
	page.sortDirection = "desc"
	page.sortField = page.TitleBar[0]
	page.Update()
}

func (page *TextVisualization) rerunFilterAndSort() {
	page.lastAnalysisType = page.AnalysisType
	page.sortField = page.TitleBar[page.sortIndex]
	page.filterField = page.TitleBar[page.filterIndex]
	page.updateFilteredData()
	page.sortData()
	page.Update()
}

func (page *TextVisualization) Render() app.UI { // no-lint

	if page.AnalysisType != page.lastAnalysisType {
		page.rerunFilterAndSort()
	}
	return app.Div().
		Class("container").
		Body(
			app.Label().
				Text("Column to filter: ").
				For("fieldsMenu"),
			app.Select().
				Name("fieldsMenu").
				ID("fieldsMenu").
				Body(
					app.Option().
						Value(page.TitleBar[0]).
						Text(page.TitleBar[0]),
					app.Option().
						Value(page.TitleBar[1]).
						Text(page.TitleBar[1]),
					app.Option().
						Value(page.TitleBar[2]).
						Text(page.TitleBar[2]),
					app.Option().
						Value(page.TitleBar[3]).
						Text(page.TitleBar[3]),
				).
				OnChange(page.selectFilterField),
			app.Input().
				Type("text").
				ID("filterText").
				Name("filterText").
				Placeholder("Filter Value").
				OnKeyup(page.selectFilterString),
			app.Table().
				Body(
					app.Tr().
						Body(app.Th().
							OnClick(func(ctx app.Context, e app.Event) {
								page.changeSort(page.TitleBar[0], 0, ctx, e)
							}).
							Body(
								app.P().Text(page.TitleBar[0]),
								app.If(page.sortField == page.TitleBar[0],
									app.If(page.sortDirection == "asc",
										app.P().Text("/\\")).
										Else(
											app.P().Text("\\/"),
										),
								),
							),
							app.Th().
								OnClick(func(ctx app.Context, e app.Event) {
									page.changeSort(page.TitleBar[1], 1, ctx, e)
								}).Body(
								app.P().Text(page.TitleBar[1]),
								app.If(page.sortField == page.TitleBar[1],
									app.If(page.sortDirection == "asc",
										app.P().Text("/\\")).
										Else(
											app.P().Text("\\/"),
										),
								),
							),
							app.Th().
								OnClick(func(ctx app.Context, e app.Event) {
									page.changeSort(page.TitleBar[2], 2, ctx, e)
								}).Body(
								app.P().Text(page.TitleBar[2]),
								app.If(page.sortField == page.TitleBar[2],
									app.If(page.sortDirection == "asc",
										app.P().Text("/\\")).
										Else(
											app.P().Text("\\/"),
										),
								),
							),
							app.Th().
								OnClick(func(ctx app.Context, e app.Event) {
									page.changeSort(page.TitleBar[3], 3, ctx, e)
								}).Body(
								app.P().Text(page.TitleBar[3]),
								app.If(page.sortField == page.TitleBar[3],
									app.If(page.sortDirection == "asc",
										app.P().Text("/\\")).
										Else(
											app.P().Text("\\/"),
										),
								),
							),
						),
					app.Range(page.filteredData).Slice(func(i int) app.UI {
						if page.AnalysisType == 1 {
							dataRow := page.filteredData[i].(analysis.CollegeGradsToJobsModel)
							return app.Tr().Body(
								app.Td().Body(
									app.P().Text(dataRow.State),
								),
								app.Td().Body(
									app.P().Text(dataRow.CollegeGrads),
								),
								app.Td().Body(
									app.P().Text(dataRow.NumberOfJobs),
								),
								app.Td().Body(
									app.P().Text(dataRow.Ratio),
								),
							)

						} else if page.AnalysisType == 2 {
							dataRow := page.filteredData[i].(analysis.DecliningBalToSalarysModel)
							return app.Tr().Body(
								app.Td().Body(
									app.P().Text(dataRow.State),
								),
								app.Td().Body(
									app.P().Text(dataRow.DecliningBalance),
								),
								app.Td().Body(
									app.P().Text(dataRow.Salary25Percent),
								),
								app.Td().Body(
									app.P().Text(dataRow.Ratio),
								),
							)

						} else {
							log.Panic("Unknown AnalysisType")
						}

						return app.Div()

					}),
				),
		)
}

func (page *TextVisualization) selectFilterField(ctx app.Context, e app.Event) {
	filterField := ctx.JSSrc.Get("value").String()
	for i, v := range page.TitleBar {
		if filterField == v {
			page.filterIndex = i
			break
		}
	}
	page.filterField = filterField

	page.updateFilteredData()

	page.sortData()
	page.Update()
}

func (page *TextVisualization) selectFilterString(ctx app.Context, e app.Event) {
	filter := ctx.JSSrc.Get("value").String()
	page.filter = filter

	page.updateFilteredData()

	page.sortData()
	page.Update()
}

func (page *TextVisualization) updateFilteredData() {
	newFilteredData := make([]interface{}, 0)
	for _, row := range page.Data {

		compareVal := page.getCompareVal(row, page.filterField)

		if strings.HasPrefix(strings.ToUpper(compareVal), strings.ToUpper(page.filter)) {
			newFilteredData = append(newFilteredData, row)
		}
	}

	page.filteredData = newFilteredData
}

func (page *TextVisualization) changeSort(field string, fieldIndex int, ctx app.Context, e app.Event) {
	if page.sortField == field {
		if page.sortDirection == "asc" {
			page.sortDirection = "desc"
		} else {
			page.sortDirection = "asc"
		}
	} else {
		page.sortField = field
		page.sortIndex = fieldIndex
		page.sortDirection = "desc"
	}

	page.sortData()

	page.Update()

}

func (page *TextVisualization) getCompareVal(row interface{}, selectedField string) string {
	var compareVal string
	if page.AnalysisType == 1 {
		parsedVal := row.(analysis.CollegeGradsToJobsModel)

		switch selectedField {
		case page.TitleBar[0]:
			compareVal = parsedVal.State

		case page.TitleBar[1]:
			compareVal = strconv.Itoa(parsedVal.CollegeGrads)

		case page.TitleBar[2]:
			compareVal = strconv.Itoa(parsedVal.NumberOfJobs)

		case page.TitleBar[3]:
			compareVal = strconv.FormatFloat(float64(parsedVal.Ratio), 'f', 6, 64)

		}

	} else if page.AnalysisType == 2 {
		parsedVal := row.(analysis.DecliningBalToSalarysModel)

		switch selectedField {
		case page.TitleBar[0]:
			compareVal = parsedVal.State

		case page.TitleBar[1]:
			compareVal = strconv.FormatFloat(float64(parsedVal.DecliningBalance), 'f', 6, 64)

		case page.TitleBar[2]:
			compareVal = strconv.Itoa(parsedVal.Salary25Percent)

		case page.TitleBar[3]:
			compareVal = strconv.FormatFloat(float64(parsedVal.Ratio), 'f', 6, 64)

		}

	} else {
		log.Panic("Unknown AnalysisType")
	}

	return compareVal
}

func (page *TextVisualization) sortData() {
	sort.Slice(page.filteredData, func(i, j int) bool {
		val1 := page.getCompareVal(page.filteredData[i], page.sortField)
		val2 := page.getCompareVal(page.filteredData[j], page.sortField)

		if int1, err := strconv.Atoi(val1); err == nil {
			int2, err := strconv.Atoi(val2)
			if err != nil {
				log.Panic("Mismatch Datatypes")
			}
			if page.sortDirection == "asc" {
				return int1 > int2
			} else {
				return int1 < int2
			}
		}

		if float1, err := strconv.ParseFloat(val1, 32); err == nil {
			float2, err := strconv.ParseFloat(val2, 32)
			if err != nil {
				log.Panic("Mismatch Datatypes")
			}
			if page.sortDirection == "asc" {
				return float1 > float2
			} else {
				return float1 < float2
			}
		}

		if page.sortDirection == "asc" {
			return val1 > val2
		}
		return val1 < val2

	})
}
