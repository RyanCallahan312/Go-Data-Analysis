// +build wasm

package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type MapVisualization struct {
	app.Compo

	TitleBar []string
	Data     []interface{}
}

func (page *MapVisualization) Render() app.UI {
	return app.P().Text("Map Render :)")
}
