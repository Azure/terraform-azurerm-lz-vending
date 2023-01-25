package roleassignment

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/check"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
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

	v := map[string]interface{}{
		"random_hex":      name,
		"role_definition": "Storage Blob Data Contributor",
	}
	testDir := filepath.Join("testdata", t.Name())
	test := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, test.Err)
	defer test.Cleanup()

	check.InPlan(test.Plan).NumberOfResourcesEquals(2).IfNotFail(t)
	defer test.DestroyWithRetry(t, setuptest.DefaultRetry)
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
	v := map[string]interface{}{
		"random_hex":      name,
		"role_definition": rd,
	}

	testDir := filepath.Join("testdata/", t.Name())
	test := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	require.NoError(t, test.Err)
	check.InPlan(test.Plan).NumberOfResourcesEquals(2).IfNotFailNow(t)

	defer test.DestroyWithRetry(t, setuptest.DefaultRetry)

	err = test.ApplyIdempotent(t)
	assert.NoError(t, err)
}
