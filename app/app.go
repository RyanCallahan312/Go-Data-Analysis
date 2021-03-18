// +build wasm

package main

import (
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	app.Route("/", &homepage{})
	app.Route("/visualization", &visualization{})
	app.Route("/updateData", &updateData{})
	app.Run()
}
