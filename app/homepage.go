package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type homepage struct {
	app.Compo
}

func (h *homepage) Render() app.UI {
	return app.Div().Body(app.H1().Text("Homepage"), app.A().Text("Update Data").Href("/updateData"), app.Br(), app.A().Text("Run Visualization").Href("/visualization"))
}
