package roleassignment

import (
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/roleassignment"
)

func TestRoleAssignmentValidWithRoleName(t *testing.T) {
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
	require.Contains(t, plan.ResourcePlannedValuesMap, "azurerm_role_assignment.this")
	ra := plan.ResourcePlannedValuesMap["azurerm_role_assignment.this"]
	assert.Equalf(t, v["role_assignment_definition"], ra.AttributeValues["role_definition_name"], "role_definition_name incorrect")
	assert.Nilf(t, ra.AttributeValues["role_definition_id"], "role_definition_id should be nil")
	assert.Equalf(t, v["role_assignment_principal_id"], ra.AttributeValues["principal_id"], "role_definition_principal_id incorrect")
}

func TestRoleAssignmentValidWithRoleDefId(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["role_assignment_definition"] = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
	require.Contains(t, plan.ResourcePlannedValuesMap, "azurerm_role_assignment.this")
	ra := plan.ResourcePlannedValuesMap["azurerm_role_assignment.this"]
	assert.Equalf(t, v["role_assignment_definition"], ra.AttributeValues["role_definition_id"], "role_definition_id incorrect")
	assert.Nilf(t, ra.AttributeValues["role_definition_name"], "role_definition_name should be nil")
	assert.Equalf(t, v["role_assignment_principal_id"], ra.AttributeValues["principal_id"], "role_definition_principal_id incorrect")
}

// TestRoleAssignmentInvalidScopes tests that the module will not accept a tenant
// or management group scope for the role assignment.
func TestRoleAssignmentInvalidScopes(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	errString := "Must begin with a subscription scope, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`. All letters must be lowercase in the subscription id."

	// test tenant scope error
	v["role_assignment_scope"] = "/"
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	assert.Contains(t, utils.SanitiseErrorMessage(err), errString)

	// test management group scope error
	v["role_assignment_scope"] = "/providers/Microsoft.Management/managementGroups/myMg"
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	assert.Contains(t, utils.SanitiseErrorMessage(err), errString)
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"role_assignment_principal_id": "00000000-0000-0000-0000-000000000000",
		"role_assignment_scope":        "/subscriptions/00000000-0000-0000-0000-000000000000",
		"role_assignment_definition":   "Owner",
	}
}
