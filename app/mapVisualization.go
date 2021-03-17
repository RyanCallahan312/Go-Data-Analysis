// +build wasm

package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type MapVisualization struct { // no-lint
	app.Compo

	AnalysisType int
	TitleBar     []string
	Data         []interface{}
}

func (page *MapVisualization) Render() app.UI { // no-lint
	return app.Div().Body(
		app.P().Text("Map Render :)"),
	)
}
