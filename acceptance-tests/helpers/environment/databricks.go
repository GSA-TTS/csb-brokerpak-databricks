// Package environment manages environment variables
package environment

import (
	"os"

	"github.com/onsi/gomega"
)

type DatabricksMetadata struct {
	Host  string
	Token string
}

func ReadDatabricksMetadata() DatabricksMetadata {
	result := DatabricksMetadata{
		Host:  os.Getenv("DATABRICKS_HOST"),
		Token: os.Getenv("DATABRICKS_TOKEN"),
	}

	gomega.Expect(result.Host).NotTo(gomega.BeEmpty(), "must set DATABRICKS_HOST")
	gomega.Expect(result.Token).NotTo(gomega.BeEmpty(), "must set DATABRICKS_TOKEN")

	return result
}
