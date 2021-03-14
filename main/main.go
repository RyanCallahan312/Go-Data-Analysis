package main

import (
	"Project1/api"
	"Project1/config"
	"Project1/database"
	"Project1/migration"
	"bufio"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, err := os.Stat(config.ProjectRootPath + "/.env")
	if err == nil {
		err := godotenv.Load(config.ProjectRootPath + "/.env")
		if err != nil {
			log.Panic(err)
		}
	}

	migration.InitalizeDB(os.Getenv("DATABASE_NAME"))
	database.DB, err = sqlx.Open("pgx", os.Getenv("WORKING_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	migration.MigrateToLatest()
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	errFile, err := os.Create("./err.txt")
	if err != nil {
		log.Panic(err)
	}
	errWriter := bufio.NewWriter(errFile)

	api.GetAPIData(errWriter)

	spreadsheet.GetSheetData()

}

func WriteToFile(data string, writer *bufio.Writer, fileLock *sync.Mutex) {

	fileLock.Lock()
	_, err := writer.WriteString(data)
	if err != nil {
		log.Panic(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Panic(err)
	}
	fileLock.Unlock()
}
