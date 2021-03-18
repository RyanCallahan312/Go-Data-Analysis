package database

import "github.com/jmoiron/sqlx"

var (
	// DB global db connection for project
	DB *sqlx.DB
)
