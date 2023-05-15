package roleassignment

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeployRoleAssignmentDefinitionName tests the deployment of a role assignment
// by defining the role definition name
func TestDeployRoleAssignmentDefinitionName(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	name, err := utils.RandomHex(4)
	require.NoErrorf(t, err, "could not generate random hex")

	v := map[string]any{
		"random_hex":      name,
		"role_definition": "Storage Blob Data Contributor",
	}
	testDir := filepath.Join("testdata", t.Name())
	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.Plan).NumberOfResourcesEquals(2).ErrorIsNil(t)
	defer test.DestroyRetry(t, setuptest.DefaultRetry)
	err = test.ApplyIdempotent(t)
	assert.NoError(t, err)
}

// TestDeployRoleAssignmentDefinitionId tests the deployment of a role assignment
// by defining a role definition id
func TestDeployRoleAssignmentDefinitionId(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	name, err := utils.RandomHex(4)
	require.NoErrorf(t, err, "could not generate random hex")

	rd := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/ba92f5b4-2d11-453d-a403-e96b0029c9fe", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	v := map[string]any{
		"random_hex":      name,
		"role_definition": rd,
	}

	testDir := filepath.Join("testdata/", t.Name())
	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	require.NoError(t, err)
	check.InPlan(test.Plan).NumberOfResourcesEquals(2).ErrorIsNilFatal(t)

	defer test.DestroyRetry(t, setuptest.DefaultRetry)

	err = test.ApplyIdempotent(t)
	assert.NoError(t, err)
}
