package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type visualization struct { //nolint
	app.Compo
}

func (h *visualization) Render() app.UI { //nolint
	return app.Div().Body(app.A().Text("To Homepage").Href("/"),
		app.H1().Text("Visualization"),
		app.Button().Text("Text Viz"),
		app.Button().Text("Map Viz"),
	)
}
