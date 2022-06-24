package virtualnetwork

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeployVirtualNetworkValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, "")
	defer utils.RemoveTestDir(t, filepath.Dir(tmp))

	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 3)
}

func TestDeployVirtualNetworkValidVnetPeering(t *testing.T) {
	utils.PreCheckDeployTests(t)

	dir := utils.GetTestDir(t)
	dir += "/testdata/" + t.Name()
	terraformOptions := utils.GetDefaultTerraformOptions(t, dir)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 3)
}

func getValidInputVariables() (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_id":                     os.Getenv("AZURE_EXISTING_SUBSCRIPTION_ID"),
		"virtual_network_address_space":       []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_location":            "northeurope",
		"virtual_network_name":                name,
		"virtual_network_resource_group_name": name,
	}, nil
}
