package terraformtests

import (
	"path"

	tfjson "github.com/hashicorp/terraform-json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "csbbrokerpakdatabricks/terraform-tests/helpers"
)

var _ = Describe("databricks workspace", Label("workspace-terraform"), Ordered, func() {
	const databricksClusterResource = "databricks_cluster"

	var (
		plan                  tfjson.Plan
		terraformProvisionDir string
	)

	defaultVars := map[string]any{
		"databricks_host":         databricksHost,
		"databricks_token":        databricksToken,
		"cluster_name":            "test-cluster",
		"spark_version":           "13.3.x-scala2.12",
		"node_type_id":            "i3.xlarge",
		"num_workers":             1,
		"autotermination_minutes": 30,
		"labels":                  map[string]string{"label1": "value1"},
	}

	BeforeAll(func() {
		terraformProvisionDir = path.Join(workingDir, "databricks/provision")
		Init(terraformProvisionDir)
	})

	Context("default values", func() {
		BeforeAll(func() {
			plan = ShowPlan(terraformProvisionDir, buildVars(defaultVars, map[string]any{}))
		})

		It("maps parameters to corresponding values", func() {
			Expect(AfterValuesForType(plan, databricksClusterResource)).To(
				MatchKeys(IgnoreExtras, Keys{
					"cluster_name":            Equal("test-cluster"),
					"spark_version":           Equal("13.3.x-scala2.12"),
					"node_type_id":            Equal("i3.xlarge"),
					"num_workers":             BeNumerically("==", 1),
					"autotermination_minutes": BeNumerically("==", 30),
					"custom_tags":             MatchAllKeys(Keys{"label1": Equal("value1")}),
				}),
			)
		})
	})

	Context("custom values", func() {
		BeforeAll(func() {
			plan = ShowPlan(terraformProvisionDir, buildVars(defaultVars, map[string]any{
				"cluster_name":            "my-custom-cluster",
				"spark_version":           "14.3.x-scala2.12",
				"node_type_id":            "m5d.large",
				"num_workers":             3,
				"autotermination_minutes": 60,
			}))
		})

		It("maps custom parameters to corresponding values", func() {
			Expect(AfterValuesForType(plan, databricksClusterResource)).To(
				MatchKeys(IgnoreExtras, Keys{
					"cluster_name":            Equal("my-custom-cluster"),
					"spark_version":           Equal("14.3.x-scala2.12"),
					"node_type_id":            Equal("m5d.large"),
					"num_workers":             BeNumerically("==", 3),
					"autotermination_minutes": BeNumerically("==", 60),
				}),
			)
		})
	})
})
