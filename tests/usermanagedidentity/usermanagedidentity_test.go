package usermanagedidentity

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/usermanagedidentity"
)

func TestUserManagedIdentity(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(3).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg[0]").Key("name").HasValue(v["resource_group_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg_lock[0]").Key("name").HasValue("lock-" + v["resource_group_name"].(string)).ErrorIsNil(t)
}

func TestUserManagedIdentityWithGitHub(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["federated_credentials_github"] = map[string]any{
		"gh1": map[string]any{
			"organization": "my-organization",
			"repository":   "my-repository",
			"entity":       "branch",
			"value":        "my-branch",
		},
		"gh2": map[string]any{
			"organization": "my-organization2",
			"repository":   "my-repository2",
			"entity":       "pull_request",
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(5).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg[0]").Key("name").HasValue(v["resource_group_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg_lock[0]").Key("name").HasValue("lock-" + v["resource_group_name"].(string)).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["ghbch-gh1"]`).Key("body").Query("properties.subject").HasValue("repo:my-organization/my-repository:ref:refs/heads/my-branch").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["ghpr-gh2"]`).Key("body").Query("properties.subject").HasValue("repo:my-organization2/my-repository2:pull_request").ErrorIsNil(t)
}

func TestUserManagedIdentityWithTFCloud(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["federated_credentials_terraform_cloud"] = map[string]any{
		"tfc1": map[string]any{
			"organization": "my-organization",
			"project":      "my-repository",
			"workspace":    "my-workspace",
			"run_phase":    "plan",
		},
		"tfc2": map[string]any{
			"organization": "my-organization",
			"project":      "my-repository",
			"workspace":    "my-workspace",
			"run_phase":    "apply",
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(5).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg[0]").Key("name").HasValue(v["resource_group_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg_lock[0]").Key("name").HasValue("lock-" + v["resource_group_name"].(string)).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["tfc-tfc1"]`).Key("body").Query("properties.subject").HasValue("organization:my-organization:project:my-repository:workspace:my-workspace:run_phase:plan").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["tfc-tfc2"]`).Key("body").Query("properties.subject").HasValue("organization:my-organization:project:my-repository:workspace:my-workspace:run_phase:apply").ErrorIsNil(t)
}

func TestUserManagedIdentityWithAdvancedFederatedCredentials(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["federated_credentials_advanced"] = map[string]any{
		"adv1": map[string]any{
			"name":               "myadvancedcred1",
			"subject_identifier": "field:value",
			"issuer_url":         "https://test",
		},
		"adv2": map[string]any{
			"name":               "myadvancedcred2",
			"subject_identifier": "field:value",
			"issuer_url":         "https://test",
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(5).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.umi").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg[0]").Key("name").HasValue(v["resource_group_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg_lock[0]").Key("name").HasValue("lock-" + v["resource_group_name"].(string)).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["adv-adv1"]`).Key("body").Query("properties.subject").HasValue("organization:my-organization:project:my-repository:workspace:my-workspace:run_phase:plan").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(`azapi_resource.umi_federated_credentials["adv-adv2"]`).Key("body").Query("properties.subject").HasValue("organization:my-organization:project:my-repository:workspace:my-workspace:run_phase:apply").ErrorIsNil(t)
}

func getMockInputVariables() map[string]any {
	return map[string]any{
		"name":                "test",
		"location":            "westeurope",
		"resource_group_name": "rg-test",
		"subscription_id":     "00000000-0000-0000-0000-000000000000",
	}
}
