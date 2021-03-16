package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// ProjectRootPath is used to tell gotdotenv where the .env file is
	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
	// Env enviorment variables from .env file
	Env map[string]string
)

// InitEnv initailizes the env variables
func InitEnv() {
	var err error
	Env, err = godotenv.Read(ProjectRootPath + "/.env")
	if err != nil {
		log.Panic(err)
	}
}
