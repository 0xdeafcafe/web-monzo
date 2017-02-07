package helpers

import "github.com/maxwellhealth/bongo"

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(connectionString, database string) *bongo.Connection {
	connection, err := bongo.Connect(&bongo.Config{
		ConnectionString: connectionString,
		Database:         database,
	})
	if err != nil {
		panic(err)
	}

	return connection
}
