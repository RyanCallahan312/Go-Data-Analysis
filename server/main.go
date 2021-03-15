// +build !wasm

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	http.Handle("/", &app.Handler{
		Name:        "Project1",
		Title:       "Project1",
		Description: "Comp490Project1",
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// package main

// import (
// 	"Project1/api"
// 	"Project1/config"
// 	"Project1/database"
// 	"Project1/migration"
// 	"Project1/spreadsheet"
// 	"bufio"
// 	"log"
// 	"os"

// 	_ "github.com/jackc/pgx/v4"
// 	_ "github.com/jackc/pgx/v4/stdlib"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	_, err := os.Stat(config.ProjectRootPath + "/.env")
// 	if err == nil {
// 		err := godotenv.Load(config.ProjectRootPath + "/.env")
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 	}

// 	migration.InitalizeDB(os.Getenv("DATABASE_NAME"))
// 	database.DB, err = sqlx.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING"))
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	err = migration.MigrateToLatest()
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	defer func() {
// 		err := database.DB.Close()
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 	}()

// 	errFile, err := os.Create("./err.txt")
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	errWriter := bufio.NewWriter(errFile)

// 	api.GetAPIData(errWriter)

// 	spreadsheet.GetSheetData()

// }
