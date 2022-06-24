package subscription

import (
	"encoding/json"
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/models"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/subscription"
)

// TestCreateNewAliasValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionCreateNewAliasValid(t *testing.T) {
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	terraformOptions.Vars = v

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))

	// Extract values from the plan and compare to the input variables.
	name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias"].AttributeValues["name"]
	bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias"].AttributeValues["body"]

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
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	v["subscription_alias_management_group_id"] = "testdeploy"

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))

	// Extract values from the plan and compare to the input variables.
	name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias"].AttributeValues["name"]
	bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias"].AttributeValues["body"]

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

// TestCreateNewAliasInvalidBillingScope tests the validation function of the subscription_alias_billing_scope variable.
func TestSubscriptionCreateNewAliasInvalidBillingScope(t *testing.T) {
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	v["subscription_alias_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	terraformOptions.Vars = v
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestCreateNewAliasInvalidWorkload tests the validation function of the subscription_alias_workload variable.
func TestSubscriptionCreateNewAliasInvalidWorkload(t *testing.T) {
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	v["subscription_alias_workload"] = "PRoduction"
	terraformOptions.Vars = v
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The workload type can be either Production or DevTest and is case sensitive.")
}

// TestCreateNewAliasInvalidManagementGroupIdInvalidChars tests the validation function of the
// subscription_alias_management_group_id variable.
func TestCreateNewAliasInvalidManagementGroupIdInvalidChars(t *testing.T) {
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	v["subscription_alias_management_group_id"] = "invalid/chars"
	terraformOptions.Vars = v
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}

// TestCreateNewAliasInvalidManagementGroupIdLength tests the validation function of the
// subscription_alias_management_group_id variable.
func TestCreateNewAliasInvalidManagementGroupIdLength(t *testing.T) {
	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	v := getMockInputVariables()
	v["subscription_alias_management_group_id"] = "tooooooooooooooooooooooooooloooooooooooooooooooooonnnnnnnnnnnnnnnnnnngggggggggggggggggggggg"
	terraformOptions.Vars = v
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.).")
}
