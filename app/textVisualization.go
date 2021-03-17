// +build wasm

package main

import (
	"Project1/analysis"
	"log"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type TextVisualization struct { // no-lint
	app.Compo

	AnalysisType int
	TitleBar     []string
	Data         []interface{}
	filter       string
	filteredData []interface{}
}

func (page *TextVisualization) OnMount(ctx app.Context) {
	page.filter = ""
	page.filteredData = page.Data
}

func (page *TextVisualization) Render() app.UI { // no-lint

	return app.Div().
		Class("container").
		Body(
			app.P().Text("Text Render :)"),
			app.Table().
				Body(
					app.Tr().
						Body(app.Th().Body(
							app.P().Text(page.TitleBar[0]),
						),
							app.Th().Body(
								app.P().Text(page.TitleBar[1]),
							),
							app.Th().Body(
								app.P().Text(page.TitleBar[2]),
							),
						),
					app.Range(page.Data).Slice(func(i int) app.UI {
						if page.AnalysisType == 1 {
							dataRow := page.Data[i].(analysis.CollegeGradsToJobsModel)
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
							)

						} else if page.AnalysisType == 2 {
							dataRow := page.Data[i].(analysis.DecliningBalToSalarysModel)
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
							)

						} else {
							log.Panic("Unknown AnalysisType")
						}

						return app.Div()

					}),
				),
		)
}
