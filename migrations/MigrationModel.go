package migrations

// MigrationModel models the migration table in the db
type MigrationModel struct {
	ID          int    `db:"migration_id"`
	Version     string `db:"version"`
	IsOnVersion bool   `db:"is_on_version"`
}
