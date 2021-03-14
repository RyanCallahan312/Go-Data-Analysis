package config

import (
	"path/filepath"
	"runtime"
)

// ProjectRootPath is used to tell gotdotenv where the .env file is
var (
	_, b, _, _      = runtime.Caller(0)
	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
)
