package tests

import (
	"encoding/json"
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/models"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateNewAliasValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionCreateNewAliasValid(t *testing.T) {
	v := getMockInputVariables()
	terraformOptions := utils.GetDefaultTerraformOptions(v)

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))

	// Extract values from the plan and compare to the input variables.
	name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["name"]
	bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["body"]

	var body models.SubscriptionAliasBody
	err = json.Unmarshal([]byte(bodyText.(string)), &body)
	require.NoErrorf(t, err, "Failed to unmarshal body JSON: %s", bodyText)

	assert.Equal(t, v["subscription_alias_name"], name)
	assert.Equal(t, v["subscription_alias_billing_scope"], *body.Properties.BillingScope)
	assert.Equal(t, v["subscription_alias_display_name"], *body.Properties.DisplayName)
	assert.Equal(t, v["subscription_alias_workload"], *body.Properties.Workload)
	assert.Nil(t, body.Properties.SubscriptionId)
}

// TestSubscriptionCreateNewAliasValidWithManagementGroup tests the
// validation functions with valid data, including a destination management group,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionCreateNewAliasValidWithManagementGroup(t *testing.T) {
	v := getMockInputVariables()
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	v["subscription_alias_management_group_id"] = "testdeploy"

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))

	// Extract values from the plan and compare to the input variables.
	name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["name"]
	bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["body"]

	var body models.SubscriptionAliasBody
	err = json.Unmarshal([]byte(bodyText.(string)), &body)
	require.NoErrorf(t, err, "Failed to unmarshal body JSON: %s", bodyText)

	assert.Equal(t, v["subscription_alias_name"], name)
	assert.Equal(t, v["subscription_alias_billing_scope"], *body.Properties.BillingScope)
	assert.Equal(t, v["subscription_alias_display_name"], *body.Properties.DisplayName)
	assert.Equal(t, v["subscription_alias_workload"], *body.Properties.Workload)
	mgResId := "/providers/Microsoft.Management/managementGroups/" + v["subscription_alias_management_group_id"].(string)
	assert.Equal(t, mgResId, *body.Properties.AdditionalProperties.ManagementGroupId)
	assert.Nil(t, body.Properties.SubscriptionId)
}

// TestSubscriptionCreateNewAliasExistingSubscriptionId tests the validation functions with valid data,
// that will create an alias for an existing subscription id.
// Then it creates a plan and compares the input variables to the planned values.
func TestSubscriptionCreateNewAliasExistingSubscriptionId(t *testing.T) {
	v := map[string]interface{}{
		"subscription_id":                 "00000000-0000-0000-0000-000000000000",
		"subscription_alias_name":         "test-subscription-alias",
		"subscription_alias_display_name": "test-subscription-alias",
	}
	terraformOptions := utils.GetDefaultTerraformOptions(v)

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	require.Equal(t, 0, len(plan.ResourcePlannedValuesMap))
	_, err = terraform.ApplyE(t, terraformOptions)
	require.NoError(t, err)
	sid := terraform.Output(t, terraformOptions, "subscription_id")
	assert.Equal(t, v["subscription_id"], sid)
	// This is commented out as we don't support creation of alias for existing subscription
	// due to complexities with testing
	//
	// name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias_existing[0]"].AttributeValues["name"]
	// bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias_existing[0]"].AttributeValues["body"]

	// var body models.SubscriptionAliasBody
	// err = json.Unmarshal([]byte(bodyText.(string)), &body)
	// require.NoErrorf(t, err, "Failed to unmarshal body JSON: %s", bodyText)

	// assert.Equal(t, v["subscription_alias_name"], name)
	// assert.Nil(t, body.Properties.DisplayName)
	// assert.Equal(t, v["subscription_id"], *body.Properties.SubscriptionId)
	// assert.Nil(t, body.Properties.BillingScope)
	// assert.Nil(t, body.Properties.Workload)
}

// TestSubscriptionCreateDisabledAlias tests the validation function with subscription_alias_enabled
// set to false.
// This should result in no resources being deployed
func TestSubscriptionCreateDisabledAlias(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_enabled"] = false
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	assert.Equal(t, 0, len(plan.ResourcePlannedValuesMap))
	_, err = terraform.ApplyE(t, terraformOptions)
	require.NoError(t, err)
	sid := terraform.Output(t, terraformOptions, "subscription_id")
	assert.Equal(t, "", sid)
}

// TestCreateNewAliasInvalidBillingScope tests the validation function of the subscription_alias_billing_scope variable.
func TestSubscriptionCreateNewAliasInvalidBillingScope(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestCreateNewAliasInvalidWorkload tests the validation function of the subscription_alias_workload variable.
func TestSubscriptionCreateNewAliasInvalidWorkload(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_workload"] = "PRoduction"
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The workload type can be either Production or DevTest and is case sensitive.")
}

// TestCreateNewAliasInvalidManagementGroupId tests the validation function of the
// subscription_alias_management_group_id variable.
func TestCreateNewAliasInvalidManagementGroupId(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_management_group_id"] = "invalid/chars"
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// TestCreateNewAliasInvalidManagementGroupId2 tests the validation function of the
// subscription_alias_management_group_id variable.
func TestCreateNewAliasInvalidManagementGroupId2(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_management_group_id"] = "tooooooooooooooooooooooooooloooooooooooooooooooooonnnnnnnnnnnnnnnnnnngggggggggggggggggggggg"
	terraformOptions := utils.GetDefaultTerraformOptions(v)
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_alias_name":          "test-subscription-alias",
		"subscription_alias_display_name":  "test-subscription-alias",
		"subscription_alias_billing_scope": "/providers/Microsoft.Billing/billingAccounts/test-billing-account",
		"subscription_alias_workload":      "Production",
	}
}
