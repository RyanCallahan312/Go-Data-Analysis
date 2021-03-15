// +build wasm

package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type hello struct {
	app.Compo
}

func (h *hello) Render() app.UI {
	return app.Div().Body(app.H1().Text("Hello World!"), app.P().Text("I am :)"))
}

func main() {
	app.Route("/", &hello{})
	app.Run()
}
