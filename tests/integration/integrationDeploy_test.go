package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/azureutils"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
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

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoErrorf(t, err, "failed to init and plan")

	require.Lenf(t, plan.ResourcePlannedValuesMap, 8, "expected 8 resources to be planned")
	resources := []string{
		"azurerm_resource_group.hub",
		"azurerm_virtual_network.hub",
		"module.alz_landing_zone.module.subscription[0].azurerm_subscription.this[0]",
		"module.alz_landing_zone.module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.alz_landing_zone.module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.alz_landing_zone.module.virtualnetwork[0].azapi_resource.rg",
		"module.alz_landing_zone.module.virtualnetwork[0].azapi_resource.vnet",
		"module.alz_landing_zone.module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	for _, r := range resources {
		require.Contains(t, plan.ResourcePlannedValuesMap, r, "expected resource %s to be planned", r)
	}

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		if err != nil {
			t.Logf("failed to cancel subscription: %v", err)
		}
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 3)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	id, err := terraform.OutputRequiredE(t, terraformOptions, "subscription_id")
	assert.NoErrorf(t, err, "failed to get subscription id output")
	u, err = uuid.Parse(id)
	assert.NoErrorf(t, err, "cannot parse subscription id as uuid: %s", id)
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
