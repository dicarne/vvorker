package funcs

// MigratePostgreSQLDatabaseFunc is a function type for migrating PostgreSQL database
// This allows the actual implementation to be set during package initialization
type MigratePostgreSQLDatabaseFunc func(userID uint64, dbid string) error

// migratePostgreSQLDatabase is the actual function that will be called
var migratePostgreSQLDatabase MigratePostgreSQLDatabaseFunc

// SetMigratePostgreSQLDatabase sets the function that will be used for database migration
// This should be called during package initialization
func SetMigratePostgreSQLDatabase(fn MigratePostgreSQLDatabaseFunc) {
	migratePostgreSQLDatabase = fn
}

// MigratePostgreSQLDatabase migrates a PostgreSQL database
// This is a wrapper around the actual implementation that can be set during initialization
func MigratePostgreSQLDatabase(userID uint64, dbid string) error {
	if migratePostgreSQLDatabase == nil {
		panic("MigratePostgreSQLDatabase function not initialized. Call SetMigratePostgreSQLDatabase during package initialization.")
	}
	return migratePostgreSQLDatabase(userID, dbid)
}
