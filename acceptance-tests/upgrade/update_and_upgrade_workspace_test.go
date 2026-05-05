package upgrade_test

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/brokers"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/random"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpgradeDatabricksWorkspaceTest", Label("workspace"), func() {
	When("upgrading broker version", func() {
		It("should continue to work", func() {
			By("pushing latest released broker version")
			serviceBroker := brokers.Create(
				brokers.WithPrefix("csb-databricks"),
				brokers.WithSourceDir(releasedBuildDir),
				brokers.WithReleasedEnv(releasedBuildDir),
			)
			defer serviceBroker.Delete()

			By("creating a service")
			serviceOffering := "csb-databricks-workspace"
			servicePlan := "default"
			serviceName := random.Name(random.WithPrefix(serviceOffering, servicePlan))
			defer services.Delete(serviceName)
			serviceInstance := services.CreateInstance(
				serviceOffering,
				servicePlan,
				services.WithBroker(serviceBroker),
				services.WithName(serviceName),
			)

			By("pushing the development version of the broker")
			serviceBroker.UpdateBroker(developmentBuildDir)

			By("triggering a no-op update to reapply the terraform for service instance")
			serviceInstance.Update(services.WithParameters(`{}`))

			By("verifying the service instance is still accessible after upgrade")
			Expect(metadata.Host).NotTo(BeEmpty(), "Databricks host should be configured")
			Expect(serviceInstance).NotTo(BeNil(), "Service instance should exist after upgrade")
		})
	})
})
