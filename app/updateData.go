package main

import (
	"Project1/dto"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type updateData struct {
	app.Compo

	fileName  string
	sheetName string
}

func (page *updateData) Render() app.UI {
	return app.Div().Body(
		app.A().Text("To Homepage").Href("/"),
		app.H1().Text("Update Data"),
		app.Input().
			Type("text").
			Placeholder("Spreadsheet File Name").
			Value(page.fileName).
			OnKeyup(page.onFileChange).
			OnChange(page.onFileChange),
		app.Input().
			Type("text").
			Placeholder("Spreadsheet Sheet Name").
			Value(page.sheetName).
			OnKeyup(page.onSheetChange).
			OnChange(page.onSheetChange),
		app.Button().
			Text("Update Data from api and "+page.fileName+".xlsx").
			OnClick(page.onUpdateData),
		app.P().
			Text(page.fileName),
	)
}

func (page *updateData) onUpdateData(ctx app.Context, e app.Event) {
	updateSheetDto := dto.UpdateSheetDTO{FileName: page.fileName, SheetName: page.sheetName}
	sheetBytes, err := json.Marshal(updateSheetDto)
	if err != nil {
		log.Panic(err)
	}
	http.Post("http://localhost:8000/sheet", "application/json", bytes.NewBuffer(sheetBytes))

	http.Post("http://localhost:8000/api", "application/json", bytes.NewBuffer(make([]byte, 0)))
	page.Update()

}

func (page *updateData) onFileChange(ctx app.Context, e app.Event) {
	file := ctx.JSSrc.Get("value").String()
	page.fileName = file
	page.Update()

}

func (page *updateData) onSheetChange(ctx app.Context, e app.Event) {
	file := ctx.JSSrc.Get("value").String()
	page.sheetName = file
	page.Update()

}
