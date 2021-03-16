// +build wasm

package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type TextVisualization struct {
	app.Compo

	TitleBar []string
	Data     []interface{}
}

func (page *TextVisualization) Render() app.UI {
	return app.P().Text("Text Render :)")
}
