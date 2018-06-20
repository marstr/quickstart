package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// DatabaseMetadata stashes a bunch of useful information about reporting the status
// of a your database.
type DatabaseMetadata struct {
	Status         string
	Dialect        string
	Connected      bool
	DatabaseExists bool
	EnvironmentVar string
}

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	c.Data()["db_meta"] = newDatabaseMetadataFromEnvironment(
		c.Param("db_env"),
		"MYSQLCONNSTR_DATABASE_URL",
		"CUSTOMCONNSTR_DATABASE_URL",
		"DATABASE_URL")

	return c.Render(200, r.HTML("index.html"))
}

func newDatabaseMetadataFromEnvironment(connectionStringVariables ...string) (dbMeta DatabaseMetadata) {
	var deets *pop.ConnectionDetails

	dbMeta.Status = "N/A"
	dbMeta.Dialect = "N/A"

	if connectionStringVariables == nil || len(connectionStringVariables) == 0 {
		connectionStringVariables = []string{
			"MYSQLCONNSTR_DATABASE_URL",
			"CUSTOMCONNSTR_DATABASE_URL",
			"DATABASE_URL",
		}
	}

	for _, p := range connectionStringVariables {
		if val := os.Getenv(p); val != "" {
			dbMeta.EnvironmentVar = p
			deets = &pop.ConnectionDetails{
				URL: val,
			}
			break
		}
	}

	if deets == nil {
		dbMeta.Status = "No connection string found."
		return
	}
	deets.Finalize()

	dbMeta.Dialect = deets.Dialect
	connection, err := pop.NewConnection(deets)
	if err != nil {
		dbMeta.Status = "Connection details aren't valid."
		return
	}

	err = connection.Open()
	if err != nil {
		dbMeta.Status = "Unable to connect to database."
		return
	}

	dbMeta.Status = "Connected."
	return
}
