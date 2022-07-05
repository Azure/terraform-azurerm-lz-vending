package roleassignment

import (
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/roleassignment"
)

func TestRoleAssignmentValid(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	// Create plan and ensure only one resource is created.

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"role_assignment_principal_id": "00000000-0000-0000-0000-000000000000",
		"role_assignment_scope":        "/subscriptions/00000000-0000-0000-0000-000000000000",
		"role_assignment_definition":   "Owner",
	}
}
