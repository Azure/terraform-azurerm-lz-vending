package tests

import (
	"encoding/json"
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateNewAliasValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestSubscriptionCreateNewAliasValid(t *testing.T) {
	v := getMockInputVariables()
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       utils.GetLogger(),
		PlanFilePath: "../tfplan",
	}

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(plan.ResourcePlannedValuesMap))

	// Extract values from the plan and compare to the input variables.
	name := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["name"]
	bodyText := plan.ResourcePlannedValuesMap["azapi_resource.subscription_alias[0]"].AttributeValues["body"]
	body := make(map[string]interface{})
	err = json.Unmarshal([]byte(bodyText.(string)), &body)
	require.NoErrorf(t, err, "Failed to unmarshal body JSON: %s", bodyText)
	bodyProperties := body["properties"].(map[string]interface{})
	assert.Equal(t, v["subscription_alias_billing_scope"], bodyProperties["billingScope"])
	assert.Equal(t, v["subscription_alias_display_name"], bodyProperties["displayName"])
	assert.Equal(t, v["subscription_alias_workload"], bodyProperties["workload"])
	assert.Equal(t, v["subscription_alias_name"], name)
}

// TestCreateNewAliasExistingSubscriptionId tests the validation functions with valid data for supplying an existing subscription id.
func TestSubscriptionCreateNewAliasExistingSubscriptionId(t *testing.T) {
	v := map[string]interface{}{
		"subscription_id":            "00000000-0000-0000-0000-000000000000",
		"subscription_alias_enabled": false,
	}
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       utils.GetLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)
}

// TestCreateNewAliasInvalidBillingScope tests the validation function of the subscription_alias_billing_scope variable.
func TestSubscriptionCreateNewAliasInvalidBillingScope(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       utils.GetLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestCreateNewAliasInvalidWorkload tests the validation function of the subscription_alias_workload variable.
func TestSubscriptionCreateNewAliasInvalidWorkload(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_workload"] = "PRoduction"
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       utils.GetLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := utils.SanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The workload type can be either Production or DevTest and is case sensitive.")
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
