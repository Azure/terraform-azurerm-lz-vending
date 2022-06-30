package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployIntegrationHubAndSpoke(t *testing.T) {
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "failed to create terraform providers file")

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables")
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 3)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)

	assert.NoError(t, err)
}

func getValidInputVariables() (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"location":                            "northeurope",
		"subscription_alias_name":             name,
		"subscription_display_name":           name,
		"subscription_billing_scope":          os.Getenv("AZURE_BILLING_SCOPE"),
		"subscription_workload":               "DevTest",
		"subscription_alias_enabled":          true,
		"virtual_network_enabled":             true,
		"virtual_network_address_space":       []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_name":                name,
		"virtual_network_resource_group_name": name,
		"virtual_network_peering_enabled":     true,
		"virtual_network_use_remote_gateways": false,
	}, nil
}
