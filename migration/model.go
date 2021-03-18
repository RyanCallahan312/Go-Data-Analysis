package migration

// Model Models the migration table in the db
type Model struct {
	ID          int    `db:"migration_id"`
	Version     string `db:"version"`
	IsOnVersion bool   `db:"is_on_version"`
}
