package roleassignment

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeployRoleAssignmentDefinitionName tests the deployment of a role assignment
// by defining the role definition name
func TestDeployRoleAssignmentDefinitionName(t *testing.T) {
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	name, err := utils.RandomHex(4)
	require.NoErrorf(t, err, "could not generate random hex")
	terraformOptions.Vars = map[string]interface{}{
		"random_hex":      name,
		"role_definition": "Storage Blob Data Contributor",
	}

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	require.Lenf(t, plan.ResourcePlannedValuesMap, 2, "expected 2 resources to be planned")

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)
}

// TestDeployRoleAssignmentDefinitionId tests the deployment of a role assignment
// by defining a role definition id
func TestDeployRoleAssignmentDefinitionId(t *testing.T) {
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	name, err := utils.RandomHex(4)
	require.NoErrorf(t, err, "could not generate random hex")
	rd := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/ba92f5b4-2d11-453d-a403-e96b0029c9fe", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	terraformOptions.Vars = map[string]interface{}{
		"random_hex":      name,
		"role_definition": rd,
	}

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	require.Lenf(t, plan.ResourcePlannedValuesMap, 2, "expected 2 resources to be planned")

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)
}
