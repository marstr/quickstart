package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// DatabaseMetadata stashes a bunch of useful information about reporting the status
// of a your database.
type DatabaseMetadata struct {
	Status         string
	Name           string
	Host           string
	Dialect        string
	Connected      bool
	DatabaseExists bool
	EnvironmentVar string
}

// ContainerMetadata stashes a bunch of useful information about the configuration
// of your Web App for Containers.
type ContainerMetadata struct {
	SiteName    string
	KuduLink    string
	SiteLink    string
	GoVersion   string
	Environment string
}

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	c.Data()["db_meta"] = newDatabaseMetadataFromEnvironment(
		c.Param("db_env"),
		"MYSQLCONNSTR_DATABASE_URL",
		"CUSTOMCONNSTR_DATABASE_URL",
		"DATABASE_URL")

	c.Data()["container_meta"] = newContainerMetadataFromEnvironment()

	return c.Render(200, r.HTML("index.html"))
}

func newContainerMetadataFromEnvironment() (cMeta ContainerMetadata) {
	cMeta.Environment = os.Getenv("GO_ENV")
	cMeta.SiteName = os.Getenv("APPSETTING_WEBSITE_SITE_NAME")

	if cMeta.SiteName == "" {
		cMeta.SiteName = "Unknown, are you sure you're running in an Azure Web App for Containers?"
		cMeta.SiteLink = "N/A"
		cMeta.KuduLink = "N/A"
		cMeta.GoVersion = "N/A"
	} else {
		cMeta.KuduLink = fmt.Sprintf("https://%s.scm.azurewebsites.net", cMeta.SiteName)
		cMeta.SiteLink = fmt.Sprintf("https://%s.azurewebsites.net", cMeta.SiteName)
	}
	cMeta.GoVersion = os.Getenv("GOLANG_VERSION")
	return
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
	dbMeta.Name = deets.Database
	dbMeta.Dialect = deets.Dialect
	dbMeta.Host = deets.Host
	connection, err := pop.NewConnection(deets)
	if err != nil {
		dbMeta.Status = "Connection details aren't valid."
		return
	}

	err = connection.Open()
	if err != nil {
		dbMeta.Status = "Unable to connect to database host."
		return
	}
	defer connection.Close()

	_, err = connection.NewTransaction()
	if err != nil {
		dbMeta.Status = "The database host was found, but the actual database hasn't been initialized."
		return
	}

	dbMeta.Status = "Connected."
	return
}
