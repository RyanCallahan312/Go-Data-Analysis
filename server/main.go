// +build !wasm

package main

import (
	"Project1/api"
	"Project1/config"
	"Project1/database"
	"Project1/dto"
	"Project1/migration"
	"Project1/shared"
	"Project1/spreadsheet"
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	appSetup()

	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	http.Handle("/", &app.Handler{
		Name:        "Project1",
		Title:       "Project1",
		Description: "Comp490Project1",
	})

	http.HandleFunc("/api", apiMethods)
	http.HandleFunc("/sheet", sheetMethods)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func appSetup() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.InitEnv()

	migration.InitalizeDB(config.Env["DATABASE_NAME"])

	var err error
	database.DB, err = sqlx.Open("pgx", config.Env["WORKING_CONNECTION_STRING"])
	if err != nil {
		log.Panic(err)
	}
	err = migration.MigrateToLatest()
	if err != nil {
		log.Panic(err)
	}

	errFile, err := os.Create("./err.txt")
	if err != nil {
		log.Panic(err)
	}
	shared.Writer = bufio.NewWriter(errFile)

}

func sheetMethods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		data := spreadsheet.GetSheetData()
		w.WriteHeader(http.StatusOK)
		dataBytes, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(dataBytes)

	case "POST":
		var updateSheetDto dto.UpdateSheetDTO
		err := json.NewDecoder(r.Body).Decode(&updateSheetDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		spreadsheet.UpdateSheetData(updateSheetDto.FileName, updateSheetDto.SheetName)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func apiMethods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		//get data
		data := api.GetApiData()
		w.WriteHeader(http.StatusOK)
		dataBytes, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(dataBytes)
	case "POST":
		//update api data
		api.UpdateApiData()
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
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
