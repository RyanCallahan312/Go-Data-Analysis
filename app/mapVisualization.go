// +build wasm

package main

import (
	"Project1/analysis"
	"log"

	grob "github.com/RyanCallahan312/go-plotly-clone/graph_objects"
	"github.com/RyanCallahan312/go-plotly-clone/offline"
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type MapVisualization struct { // no-lint
	app.Compo

	AnalysisType     int
	lastAnalysisType int
	TitleBar         []string
	Data             []interface{}
	HtmlMap          string
}

func (page *MapVisualization) OnMount(ctx app.Context) { // no-lint
	page.lastAnalysisType = page.AnalysisType
	page.getMap()
}

func (page *MapVisualization) Render() app.UI { // no-lint

	if page.AnalysisType != page.lastAnalysisType || page.HtmlMap == "" {
		page.lastAnalysisType = page.AnalysisType
		page.getMap()

	}
	log.Println(page.HtmlMap)
	return app.Div().Body(
		app.P().Text("Map Render :)"),
		app.If(page.HtmlMap != "", app.IFrame().SrcDoc(page.HtmlMap)),
	)
}

func (page *MapVisualization) getMap() {

	locations := make([]string, 0)
	locationsAbrv := make([]string, 0)
	metrics := make([]float32, 0)
	for _, val := range page.Data {
		if page.AnalysisType == 1 {
			parsedVal := val.(analysis.CollegeGradsToJobsModel)

			locations = append(locations, parsedVal.State)
			locationsAbrv = append(locationsAbrv, nameToAbbrv[parsedVal.State])
			metrics = append(metrics, parsedVal.Ratio)

		} else if page.AnalysisType == 2 {
			parsedVal := val.(analysis.DecliningBalToSalarysModel)

			locations = append(locations, parsedVal.State)
			locationsAbrv = append(locationsAbrv, nameToAbbrv[parsedVal.State])
			metrics = append(metrics, parsedVal.Ratio)
		} else {
			log.Panic("Unknown AnalysisType")
		}

	}

	var title string
	if page.AnalysisType == 1 {
		title = "Ratio of College Grads to Jobs"
	} else if page.AnalysisType == 2 {
		title = "Ratio of Three Year College Graduate Declining Balance to 25th Percentile Salary"
	} else {
		log.Panic("Unknown AnalysisType")
	}

	fig := &grob.Fig{
		Data: grob.Traces{
			&grob.Choropleth{
				Type:           grob.TraceTypeChoropleth,
				Autocolorscale: grob.True,
				Locationmode:   grob.ChoroplethLocationmode_USA_states,
				Locations:      locationsAbrv,
				Z:              metrics,
				Text:           locations,
			},
		},
		Layout: &grob.Layout{
			Title: &grob.LayoutTitle{
				Text: title,
			},
			Geo: &grob.LayoutGeo{
				Scope: grob.LayoutGeoScope_usa,
			},
		},
	}

	bytes := offline.FigToBuffer(fig)
	page.HtmlMap = bytes.String()
	log.Println(page.HtmlMap)

	page.Update()
}

var nameToAbbrv = map[string]string{
	"Alabama":        "AL",
	"Alaska":         "AK",
	"Arizona":        "AZ",
	"Arkansas":       "AR",
	"California":     "CA",
	"Colorado":       "CO",
	"Connecticut":    "CT",
	"Delaware":       "DE",
	"Florida":        "FL",
	"Georgia":        "GA",
	"Hawaii":         "HI",
	"Idaho":          "ID",
	"Illinois":       "IL",
	"Indiana":        "IN",
	"Iowa":           "IA",
	"Kansas":         "KS",
	"Kentucky":       "KY",
	"Louisiana":      "LA",
	"Maine":          "ME",
	"Maryland":       "MD",
	"Massachusetts":  "MA",
	"Michigan":       "MI",
	"Minnesota":      "MN",
	"Mississippi":    "MS",
	"Missouri":       "MO",
	"Montana":        "MT",
	"Nebraska":       "NE",
	"Nevada":         "NV",
	"New Hampshire":  "NH",
	"New Jersey":     "NJ",
	"New Mexico":     "NM",
	"New York":       "NY",
	"North Carolina": "NC",
	"North Dakota":   "ND",
	"Ohio":           "OH",
	"Oklahoma":       "OK",
	"Oregon":         "OR",
	"Pennsylvania":   "PA",
	"Rhode Island":   "RI",
	"South Carolina": "SC",
	"South Dakota":   "SD",
	"Tennessee":      "TN",
	"Texas":          "TX",
	"Utah":           "UT",
	"Vermont":        "VT",
	"Virginia":       "VA",
	"Washington":     "WA",
	"West Virginia":  "WV",
	"Wisconsin":      "WI",
	"Wyoming":        "WY",
}
