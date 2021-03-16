package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type visualization struct {
	app.Compo
}

func (h *visualization) Render() app.UI {
	return app.Div().Body(app.A().Text("To Homepage").Href("/"),
		app.H1().Text("Visualization"),
		app.Button().Text("Text Viz"),
		app.Button().Text("Map Viz"),
	)
}
