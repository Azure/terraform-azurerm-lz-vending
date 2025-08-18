package subscription

import (
	"os"
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
// This test uses the azapi provider.
func TestSubscriptionAliasCreateValid(t *testing.T) {

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		"azapi_resource.subscription[0]",
		"azapi_resource_action.subscription_rename[0]",
		"azapi_update_resource.subscription_tags[0]",
		"azapi_resource_action.subscription_cancel[0]",
		"time_sleep.wait_for_subscription_before_subscription_operations[0]",
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)

	for _, res := range resources {
		check.InPlan(test.PlanStruct).That(res).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("name").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.billingScope").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.displayName").HasValue(v["subscription_display_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.tags").HasValue(v["subscription_tags"]).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateValidWithManagementGroup tests the
// validation functions with valid data, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
// This test uses the azapi provider.
func TestSubscriptionAliasCreateValidWithManagementGroup(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_management_group_id"] = os.Getenv("ARM_TENANT_ID")
	v["subscription_management_group_association_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		"terraform_data.replacement[0]",
		"azapi_resource.subscription[0]",
		"azapi_resource_action.subscription_rename[0]",
		"azapi_update_resource.subscription_tags[0]",
		"azapi_resource_action.subscription_cancel[0]",
		"azapi_resource_action.subscription_association[0]",
		"time_sleep.wait_for_subscription_before_subscription_operations[0]",
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)

	for _, res := range resources {
		check.InPlan(test.PlanStruct).That(res).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("name").HasValue(v["subscription_alias_name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.billingScope").HasValue(v["subscription_billing_scope"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.workload").HasValue(v["subscription_workload"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.tags").HasValue(v["subscription_tags"]).ErrorIsNil(t)

	mgResID := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	check.InPlan(test.PlanStruct).That("azapi_resource.subscription[0]").Key("body").Query("properties.additionalProperties.managementGroupId").HasValue(mgResID).ErrorIsNil(t)
}

// TestSubscriptionAliasCreateInvalidBillingScope tests the validation function of the subscription_billing_scope variable.
func TestSubscriptionAliasCreateInvalidBillingScope(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidWorkload tests the validation function of the subscription_workload variable.
func TestSubscriptionAliasCreateInvalidWorkload(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_workload"] = "PRoduction"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.ErrorContains(t, err, "The workload type can be either Production or DevTest and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "invalid/chars"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdLength tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdLength(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_management_group_id"] = "tooooooooooooooooooooooooooloooooooooooooooooooooonnnnnnnnnnnnnnnnnnngggggggggggggggggggggg"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()

	assert.Contains(t, utils.SanitiseErrorMessage(err), "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

func TestSubscriptionInvalidTagValue(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_tags"] = map[string]any{
		"illegal-value": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum mattis velit quis nisl dictum, nec aliquet velit bibendum. Sed et ante nec arcu convallis rutrum. Nulla sed velit ac quam finibus volutpat! Duis malesuada leo nec eros laoreet, vel consectetur enim eleifend. Sed at fermentum libero. Proin sodales lectus quis est volutpat, id suscipit purus eleifend. Vivamus dignissim nulla nec dui sollicitudin, quis pharetra ipsum posuere. Pellentesque eget magna sit amet metus fermentum hendrerit ut non velit. Donec accumsan eros nec nibh porttitor, non interdum elit laoreet. Nam gravida elit ac turpis tristique, a facilisis orci suscipit. Sed eget luctus velit. Integer quis nulla nec ante tempus congue vitae id sem. Nam eget felis non risus fringilla tempor. Integer aliquam facilisis aliquam&.",
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.Contains(t, utils.SanitiseErrorMessage(err), "Tag values must be between 0-256 characters.")
}

func TestSubscriptionInvalidTagName(t *testing.T) {
	var tagname string
	for i := 0; i < 513; i++ {
		tagname += "a"
	}
	v := getMockInputVariables()
	v["subscription_tags"] = map[string]any{
		tagname: "illegal-name",
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.Contains(t, utils.SanitiseErrorMessage(err), "Tag name must contain neither `<>%&\\?/` nor control characters, and must be between 0-512 characters.")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]any {
	return map[string]any{
		"subscription_alias_enabled": true,
		"subscription_alias_name":    "test-subscription-alias",
		"subscription_billing_scope": "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_display_name":  "test-subscription-alias",
		"subscription_workload":      "Production",
		"subscription_tags": map[string]any{
			"test-tag":   "test-value",
			"test-tag2:": "test-value2",
		},
	}
}
