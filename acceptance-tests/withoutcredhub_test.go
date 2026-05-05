package acceptance_test

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/apps"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/brokers"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/matchers"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/random"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Without CredHub", Label("withoutcredhub"), func() {
	It("can be accessed by an app", func() {
		broker := brokers.Create(
			brokers.WithPrefix("csb-no-credhub"),
			brokers.WithLatestEnv(),
			brokers.WithEnv(apps.EnvVar{Name: "CH_CRED_HUB_URL", Value: ""}),
		)
		defer broker.Delete()

		By("creating a service instance")
		serviceOffering := "csb-databricks-workspace"
		servicePlan := "default"
		serviceName := random.Name(random.WithPrefix(serviceOffering, servicePlan))
		defer services.Delete(serviceName)
		serviceInstance := services.CreateInstance(
			serviceOffering,
			servicePlan,
			services.WithBroker(broker),
			services.WithName(serviceName),
		)

		By("pushing the unstarted app")
		app := apps.Push(apps.WithApp(apps.Databricks))
		defer apps.Delete(app)

		By("binding the app to the databricks workspace service instance")
		binding := serviceInstance.Bind(app)

		By("starting the app")
		apps.Start(app)

		By("checking that the app environment does not a credhub reference for credentials")
		Expect(binding.Credential()).NotTo(matchers.HaveCredHubRef)
	})
})
