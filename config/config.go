package config

import (
	"os"
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
		Env = make(map[string]string)
		Env["API_KEY"] = os.Getenv("API_KEY")
		Env["MAINTENANCE_CONNECTION_STRING"] = os.Getenv("MAINTENANCE_CONNECTION_STRING")
		Env["WORKING_CONNECTION_STRING"] = os.Getenv("WORKING_CONNECTION_STRING")
		Env["TEST_CONNECTION_STRING"] = os.Getenv("TEST_CONNECTION_STRING")
		Env["DATABASE_NAME"] = os.Getenv("DATABASE_NAME")
	}
}
