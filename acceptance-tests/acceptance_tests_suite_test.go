package acceptance_test

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/environment"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Tests Suite")
}

var DatabricksMetadata environment.DatabricksMetadata

var _ = BeforeSuite(func() {
	DatabricksMetadata = environment.ReadDatabricksMetadata()
})
