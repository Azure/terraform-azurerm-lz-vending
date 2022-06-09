package alzLandingZoneTfModuleTest

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestCreateNewAliasValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestCreateNewAliasValid(t *testing.T) {
	v := getMockInputVariables()
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       getLogger(),
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
	if err != nil {
		t.Errorf("Could not unmarshal subscription alias resource body, %s", err)
	}
	bodyProperties := body["properties"].(map[string]interface{})
	assert.Equal(t, v["subscription_alias_billing_scope"], bodyProperties["billingScope"])
	assert.Equal(t, v["subscription_alias_display_name"], bodyProperties["displayName"])
	assert.Equal(t, v["subscription_alias_workload"], bodyProperties["workload"])
	assert.Equal(t, v["subscription_alias_name"], name)
}

// TestCreateNewAliasExistingSubscriptionId tests the validation functions with valid data for supplying an existing subscription id.
func TestCreateNewAliasExistingSubscriptionId(t *testing.T) {
	v := map[string]interface{}{
		"subscription_id":            "00000000-0000-0000-0000-000000000000",
		"subscription_alias_enabled": false,
	}
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       getLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	assert.NoError(t, err)
}

// TestCreateNewAliasInvalidBillingScope tests the validation function of the subscription_alias_billing_scope variable.
func TestCreateNewAliasInvalidBillingScope(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_billing_scope"] = "/PRoviders/Microsoft.Billing/billingAccounts/test-billing-account"
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       getLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := sanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive.")
}

// TestCreateNewAliasInvalidWorkload tests the validation function of the subscription_alias_workload variable.
func TestCreateNewAliasInvalidWorkload(t *testing.T) {
	v := getMockInputVariables()
	v["subscription_alias_workload"] = "PRoduction"
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       getLogger(),
	}
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	errMessage := sanitiseErrorMessage(err)
	assert.Contains(t, errMessage, "The workload type can be either Production or DevTest and is case sensitive.")
}

// sanitiseErrorMessage replaces the newline characters in an error.Error() output with a single space to allow us to check for the entire error message.
// We need to do this because Terraform adds newline characters depending on the width of the console window.
// TODO: Test on Windows if we get \r\n instead of just \n.
func sanitiseErrorMessage(err error) string {
	return strings.Replace(err.Error(), "\n", " ", -1)
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

// getLogger returns a logger that can be used for testing.
// The default logger will discard the Terraform output.
// Set TERRATEST_LOGGER to a non empty value to enable verbose logging.
func getLogger() *logger.Logger {
	if os.Getenv("TERRATEST_LOGGER") != "" {
		return logger.Terratest
	}
	return logger.Discard
}
