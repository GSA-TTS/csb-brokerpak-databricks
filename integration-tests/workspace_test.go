package integration_test

import (
	testframework "github.com/cloudfoundry/cloud-service-broker/v2/brokerpaktestframework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

const (
	workspaceServiceName        = "csb-databricks-workspace"
	workspaceServiceID          = "8a8da4d7-1f3e-4d3e-9c1a-2b4f8c3e1d5a"
	workspaceServiceDisplayName = "Databricks Workspace"
	workspaceServiceDescription = "Databricks workspace cluster provisioned via OpenTofu."
	workspaceServiceSupportURL  = "https://databricks.com/support"
	workspaceDefaultPlanName    = "default"
	workspaceDefaultPlanID      = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
)

var customDatabricksWorkspacePlans = []map[string]any{
	customDatabricksWorkspacePlan,
}

var customDatabricksWorkspacePlan = map[string]any{
	"name": workspaceDefaultPlanName,
	"id":   workspaceDefaultPlanID,
	"metadata": map[string]any{
		"displayName": workspaceServiceDisplayName,
	},
}

var _ = Describe("Databricks Workspace", Label("workspace"), func() {
	BeforeEach(func() {
		Expect(mockTerraform.SetTFState([]testframework.TFStateValue{})).To(Succeed())
	})

	AfterEach(func() {
		Expect(mockTerraform.Reset()).To(Succeed())
	})

	It("publishes in the catalog", func() {
		catalog, err := broker.Catalog()
		Expect(err).NotTo(HaveOccurred())

		service := testframework.FindService(catalog, workspaceServiceName)
		Expect(service.ID).To(Equal(workspaceServiceID))
		Expect(service.Description).To(Equal(workspaceServiceDescription))
		Expect(service.Tags).To(ConsistOf("databricks"))
		Expect(service.Metadata.ImageUrl).To(ContainSubstring("data:image/png;base64,"))
		Expect(service.Metadata.DisplayName).To(Equal(workspaceServiceDisplayName))
		Expect(service.Metadata.DocumentationUrl).To(Equal(cloudServiceBrokerDocumentationURL))
		Expect(service.Metadata.ProviderDisplayName).To(Equal(providerDisplayName))
		Expect(service.Metadata.SupportUrl).To(Equal(workspaceServiceSupportURL))
		Expect(service.Plans).To(
			ConsistOf(
				MatchFields(IgnoreExtras, Fields{
					ID:   Equal(workspaceDefaultPlanID),
					Name: Equal(workspaceDefaultPlanName),
				}),
			),
		)
	})

	Describe("provisioning", func() {
		It("should provision a plan with defaults", func() {
			instanceID, err := broker.Provision(workspaceServiceName, workspaceDefaultPlanName, map[string]any{})
			Expect(err).NotTo(HaveOccurred())

			Expect(mockTerraform.FirstTerraformInvocationVars()).To(
				SatisfyAll(
					HaveKeyWithValue("cluster_name", ContainSubstring(instanceID)),
					HaveKeyWithValue("spark_version", "13.3.x-scala2.12"),
					HaveKeyWithValue("node_type_id", "i3.xlarge"),
					HaveKeyWithValue("num_workers", BeNumerically("==", 1)),
					HaveKeyWithValue("autotermination_minutes", BeNumerically("==", 30)),
					HaveKeyWithValue("databricks_host", brokerDatabricksHost),
					HaveKeyWithValue("databricks_token", brokerDatabricksToken),
				),
			)
		})

		It("should allow properties to be set on provision", func() {
			_, err := broker.Provision(workspaceServiceName, workspaceDefaultPlanName, map[string]any{
				"cluster_name":            "my-cluster",
				"spark_version":           "14.3.x-scala2.12",
				"node_type_id":            "m5d.large",
				"num_workers":             3,
				"autotermination_minutes": 60,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(mockTerraform.FirstTerraformInvocationVars()).To(
				SatisfyAll(
					HaveKeyWithValue("cluster_name", "my-cluster"),
					HaveKeyWithValue("spark_version", "14.3.x-scala2.12"),
					HaveKeyWithValue("node_type_id", "m5d.large"),
					HaveKeyWithValue("num_workers", BeNumerically("==", 3)),
					HaveKeyWithValue("autotermination_minutes", BeNumerically("==", 60)),
				),
			)
		})

		Describe("updating instance", func() {
			var instanceID string

			BeforeEach(func() {
				var err error
				instanceID, err = broker.Provision(workspaceServiceName, workspaceDefaultPlanName, nil)
				Expect(err).NotTo(HaveOccurred())
			})

			DescribeTable(
				"preventing updates with `prohibit_update` as it can force resource replacement or re-creation",
				func(prop string, value any) {
					err := broker.Update(instanceID, workspaceServiceName, workspaceDefaultPlanName, map[string]any{prop: value})

					Expect(err).To(MatchError(
						ContainSubstring(
							"attempt to update parameter that may result in service instance re-creation and data loss",
						),
					))

					const initialProvisionInvocation = 1
					Expect(mockTerraform.ApplyInvocations()).To(HaveLen(initialProvisionInvocation))
				},
				Entry("cluster_name", "cluster_name", "some-other-name"),
			)
		})
	})
})
