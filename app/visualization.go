// +build wasm

package main

import (
	"Project1/analysis"
	"Project1/dto"
	"encoding/json"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type visualization struct { //nolint
	app.Compo

	analysisType int
	vizType      int
	titleBar     []string
	data         []interface{}
}

func (page *visualization) Render() app.UI { //nolint
	return app.Div().Body(app.A().Text("To Homepage").Href("/"),
		app.H1().Text("Visualization"),
		app.Label().
			Text("Compare amount of college grads to amount of jobs").
			For("analysisType1"),
		app.Input().
			Type("radio").
			ID("analysisType1").
			Name("analysisType").
			OnChange(func(ctx app.Context, e app.Event) {
				page.selectAnalysisType(1, ctx, e)
			}),
		app.Br(),
		app.Label().
			Text("3 year graduate cohort declining balance percentage to the 25 percent salary in the state").
			For("analysisType2"),
		app.Input().
			Type("radio").
			ID("analysisType2").
			Name("analysisType").
			OnChange(func(ctx app.Context, e app.Event) {
				page.selectAnalysisType(2, ctx, e)
			}),
		app.Br(),
		app.Br(),
		app.Label().
			Text("Text Visualization").
			For("vizType1"),
		app.Input().
			Type("radio").
			ID("vizType1").
			Name("vizType").
			OnChange(func(ctx app.Context, e app.Event) {
				page.selectVizType(1, ctx, e)
				page.getData(ctx, e)
			}),
		app.Br(),
		app.Label().
			Text("Map Visualization").
			For("vizType2"),
		app.Input().
			Type("radio").
			ID("vizType2").
			Name("vizType").
			OnChange(func(ctx app.Context, e app.Event) {
				page.selectVizType(2, ctx, e)
				page.getData(ctx, e)
			}),
		app.If(page.vizType == 1 && (page.data != nil && len(page.data) > 0), &TextVisualization{Data: page.data, TitleBar: page.titleBar}),
		app.If(page.vizType == 2 && (page.data != nil && len(page.data) > 0), &MapVisualization{Data: page.data, TitleBar: page.titleBar}),
	)
}

func (page *visualization) selectAnalysisType(analysisType int, ctx app.Context, e app.Event) {
	page.analysisType = analysisType
}

func (page *visualization) selectVizType(vizType int, ctx app.Context, e app.Event) {
	page.vizType = vizType
}

func (page *visualization) getData(ctx app.Context, e app.Event) {
	go func() {
		var scorecardData []dto.CollegeScoreCardFieldsDTO
		resp, err := http.Get("http://localhost:8000/api")
		if err != nil {
			log.Panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Panic(resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&scorecardData)
		if err != nil {
			log.Panic(err)
		}

		var jobData []dto.JobDataDTO
		resp, err = http.Get("http://localhost:8000/sheet")
		if err != nil {
			log.Panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Panic(resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&jobData)
		if err != nil {
			log.Panic(err)
		}

		if page.analysisType == 1 {
			page.data = analysis.CollegeGradsToAmountOfJobs(scorecardData, jobData)
		} else {
			page.data = analysis.CollegeGradsToAmountOfJobs(scorecardData, jobData)
		}
		page.titleBar = []string{"State", "College Grads 2018", "Total Employment"}

		log.Println("gotData")
		app.Dispatch(func() { page.Update() })
	}()

}
