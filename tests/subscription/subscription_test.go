package subscription

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/subscription"
)

// TestSubscriptionAliasCreateValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionAliasCreateValid(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("alias").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("billing_scope_id").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("subscription_name").HasValue(v["subscription_display_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("tags").HasValue(v["subscription_tags"]).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
// This test uses the azapi provider.
func TestSubscriptionAliasCreateValidAzApi(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_use_azapi"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("name").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.billingScope").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.displayName").HasValue(v["subscription_display_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.tags").HasValue(v["subscription_tags"]).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateValidWithManagementGroup tests the
// validation functions with valid data, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionAliasCreateValidWithManagementGroup(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "testdeploy"
	v["subscription_management_group_association_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(2).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("alias").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("billing_scope_id").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("subscription_name").HasValue(v["subscription_display_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_subscription.this[0]").Key("workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)

	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	check.InPlan(test.PlanStruct).That("azurerm_management_group_subscription_association.this[0]").Key("management_group_id").HasValue(mgResId).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateValidWithManagementGroupAzApi tests the
// validation functions with valid data, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
// This test uses the azapi provider.
func TestSubscriptionAliasCreateValidWithManagementGroupAzApi(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "testdeploy"
	v["subscription_management_group_association_enabled"] = true
	v["subscription_use_azapi"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureCliAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(2).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("name").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.billingScope").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.tags").HasValue(v["subscription_tags"]).ErrorIsNil(t)

	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.managementGroupId").HasValue(mgResId).ErrorIsNil(t)
}

// TestSubscriptionExistingWithManagementGroup tests the
// validation functions with an existing subscription id, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionExistingWithManagementGroup(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "testdeploy"
	v["subscription_management_group_association_enabled"] = true
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)

	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	check.InPlan(test.PlanStruct).That("azurerm_management_group_subscription_association.this[0]").Key("management_group_id").HasValue(mgResId).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateInvalidBillingScope tests the validation function of the subscription_billing_scope variable.
func TestSubscriptionAliasCreateInvalidBillingScope(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidWorkload tests the validation function of the subscription_workload variable.
func TestSubscriptionAliasCreateInvalidWorkload(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_workload"] = "PRoduction"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.ErrorContains(t, err, "The workload type can be either Production or DevTest and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "invalid/chars"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdLength tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdLength(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "tooooooooooooooooooooooooooloooooooooooooooooooooonnnnnnnnnnnnnnnnnnngggggggggggggggggggggg"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]any {
	return map[string]any{
		"subscription_alias_enabled": true,
		"subscription_alias_name":    "test-subscription-alias",
		"subscription_billing_scope": "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_display_name":  "test-subscription-alias",
		"subscription_use_azapi":     false,
		"subscription_workload":      "Production",
		"subscription_tags": map[string]any{
			"test-tag":   "test-value",
			"test-tag2:": "test-value2",
		},
	}
}
