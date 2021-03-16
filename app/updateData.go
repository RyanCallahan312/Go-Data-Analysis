// +build wasm

package main

import (
	"Project1/dto"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type updateData struct { //nolint
	app.Compo

	fileName       string
	sheetName      string
	loadingMessage string
}

func (page *updateData) Render() app.UI { //nolint
	return app.Div().Body(
		app.A().Text("To Homepage").Href("/"),
		app.H1().Text("Update Data"),
		app.H3().Text("***specified file must be in root directory of the go project***"),
		app.Label().
			Text("File Name").
			For("fileInput"),
		app.Input().
			Type("text").
			ID("fileInput").
			Name("fileInput").
			Placeholder("Spreadsheet File Name").
			Value(page.fileName).
			OnKeyup(page.onFileChange).
			OnChange(page.onFileChange),
		app.Br(),
		app.Label().
			Text("Sheet Name").
			For("sheetInput"),
		app.Input().
			Type("text").
			ID("sheetInput").
			Name("sheetInput").
			Placeholder("Spreadsheet Sheet Name").
			Value(page.sheetName).
			OnKeyup(page.onSheetChange).
			OnChange(page.onSheetChange),
		app.Br(),
		app.Button().
			Text("Update Data from api and "+page.fileName+".xlsx").
			OnClick(page.onUpdateData),
		app.P().
			Text(page.loadingMessage),
	)
}

func (page *updateData) onUpdateData(ctx app.Context, e app.Event) { //nolint
	go func() {
		app.Dispatch(func() {
			page.loadingMessage = "Updating data"
			page.Update()
			log.Println("Updating data")
		})
		updateSheetDto := dto.UpdateSheetDTO{FileName: page.fileName, SheetName: page.sheetName}
		sheetBytes, err := json.Marshal(updateSheetDto)
		if err != nil {
			log.Panic(err)
		}
		_, err = http.Post("http://localhost:8000/sheet", "application/json", bytes.NewBuffer(sheetBytes))
		if err != nil {
			log.Println(err)
		}

		_, err = http.Post("http://localhost:8000/api", "application/json", bytes.NewBuffer(make([]byte, 0)))
		if err != nil {
			log.Println(err)
		}
		app.Dispatch(func() {
			page.loadingMessage = "Data updated"
			page.Update()
			log.Println("Data updated")
		})
	}()

}

func (page *updateData) onFileChange(ctx app.Context, e app.Event) { //nolint
	file := ctx.JSSrc.Get("value").String()
	page.fileName = file
	page.Update()

}

func (page *updateData) onSheetChange(ctx app.Context, e app.Event) { //nolint
	file := ctx.JSSrc.Get("value").String()
	page.sheetName = file
	page.Update()

}
