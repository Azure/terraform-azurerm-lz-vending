package subscription

import (
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
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
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
	terraform.RequirePlannedValuesMapKeyExists(t, plan, "azurerm_subscription.this[0]")

	subalias := plan.ResourcePlannedValuesMap["azurerm_subscription.this[0]"]
	require.Contains(t, subalias.AttributeValues, "alias")
	require.Contains(t, subalias.AttributeValues, "billing_scope_id")
	require.Contains(t, subalias.AttributeValues, "subscription_name")
	require.Contains(t, subalias.AttributeValues, "workload")
	require.Contains(t, subalias.AttributeValues, "tags")

	assert.Equal(t, v["subscription_alias_name"], subalias.AttributeValues["alias"])
	assert.Equal(t, v["subscription_billing_scope"], subalias.AttributeValues["billing_scope_id"])
	assert.Equal(t, v["subscription_display_name"], subalias.AttributeValues["subscription_name"])
	assert.Equal(t, v["subscription_workload"], subalias.AttributeValues["workload"])
	assert.Equal(t, v["subscription_tags"], subalias.AttributeValues["tags"])
}

// TestSubscriptionAliasCreateValidWithManagementGroup tests the
// validation functions with valid data, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionAliasCreateValidWithManagementGroup(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_management_group_id"] = "testdeploy"
	v["subscription_management_group_association_enabled"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 2, len(plan.ResourcePlannedValuesMap))
	terraform.RequirePlannedValuesMapKeyExists(t, plan, "azurerm_subscription.this[0]")
	terraform.RequirePlannedValuesMapKeyExists(t, plan, "azurerm_management_group_subscription_association.this[0]")

	subalias := plan.ResourcePlannedValuesMap["azurerm_subscription.this[0]"]
	require.Contains(t, subalias.AttributeValues, "alias")
	require.Contains(t, subalias.AttributeValues, "billing_scope_id")
	require.Contains(t, subalias.AttributeValues, "subscription_name")
	require.Contains(t, subalias.AttributeValues, "workload")

	assert.Equal(t, v["subscription_alias_name"], subalias.AttributeValues["alias"])
	assert.Equal(t, v["subscription_billing_scope"], subalias.AttributeValues["billing_scope_id"])
	assert.Equal(t, v["subscription_display_name"], subalias.AttributeValues["subscription_name"])
	assert.Equal(t, v["subscription_workload"], subalias.AttributeValues["workload"])

	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	mg := plan.ResourcePlannedValuesMap["azurerm_management_group_subscription_association.this[0]"]
	require.Contains(t, mg.AttributeValues, "management_group_id")
	assert.Equal(t, mgResId, mg.AttributeValues["management_group_id"])
}

// TestSubscriptionExistingWithManagementGroup tests the
// validation functions with an existing subscription id, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionExistingWithManagementGroup(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_management_group_id"] = "testdeploy"
	v["subscription_management_group_association_enabled"] = true
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
	terraform.RequirePlannedValuesMapKeyExists(t, plan, "azurerm_management_group_subscription_association.this[0]")

	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_management_group_id"].(string)
	mg := plan.ResourcePlannedValuesMap["azurerm_management_group_subscription_association.this[0]"]
	require.Contains(t, mg.AttributeValues, "management_group_id")
	assert.Equal(t, mgResId, mg.AttributeValues["management_group_id"])
}

// TestSubscriptionAliasCreateInvalidBillingScope tests the validation function of the subscription_billing_scope variable.
func TestSubscriptionAliasCreateInvalidBillingScope(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	terraformOptions.Vars = v
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidWorkload tests the validation function of the subscription_workload variable.
func TestSubscriptionAliasCreateInvalidWorkload(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_workload"] = "PRoduction"
	terraformOptions.Vars = v
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The workload type can be either Production or DevTest and is case sensitive.")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdInvalidChars(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_management_group_id"] = "invalid/chars"
	terraformOptions.Vars = v
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// TestSubscriptionAliasCreateInvalidManagementGroupIdLength tests the validation function of the
// subscription_alias_management_group_id variable.
func TestSubscriptionAliasCreateInvalidManagementGroupIdLength(t *testing.T) {
	t.Parallel()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_management_group_id"] = "tooooooooooooooooooooooooooloooooooooooooooooooooonnnnnnnnnnnnnnnnnnngggggggggggggggggggggg"
	terraformOptions.Vars = v
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_alias_enabled": true,
		"subscription_alias_name":    "test-subscription-alias",
		"subscription_display_name":  "test-subscription-alias",
		"subscription_billing_scope": "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_workload":      "Production",
		"subscription_tags": map[string]interface{}{
			"test-tag":   "test-value",
			"test-tag2:": "test-value2",
		},
	}
}
