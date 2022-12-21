package virtualnetwork

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeployVirtualNetworkValid tests the deployment of virtual networks
// with valid input variables.
func TestDeployVirtualNetworkValid(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	numres := 8
	require.Lenf(t, plan.ResourcePlannedValuesMap, numres, "expected %d resources, got %d", numres, len(plan.ResourcePlannedValuesMap))

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		terraform.RequirePlannedValuesMapKeyExists(t, plan, r)
	}

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// check there two outputs for the virtual network resource ids
	vnri, err := terraform.OutputMapE(t, terraformOptions, "virtual_network_resource_ids")
	require.NoErrorf(t, err, "could not get virtual_network_resource_ids output, %s", err)
	assert.Lenf(t, vnri, 2, "expected 2 virtual networks, got %d", len(vnri))
}

// TestDeployVirtualNetworkValidCustomDns tests the deployment of virtual networks
// with valid input variables and custom DNS servers.
func TestDeployVirtualNetworkValidCustomDns(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["secondary"]
	primaryvnet["dns_servers"] = []string{"192.168.0.250", "192.168.0.251"}
	secondaryvnet["dns_servers"] = []string{"192.168.1.250", "192.168.1.251"}
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	numres := 8
	require.Lenf(t, plan.ResourcePlannedValuesMap, numres, "expected %d resources, got %d", numres, len(plan.ResourcePlannedValuesMap))

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		terraform.RequirePlannedValuesMapKeyExists(t, plan, r)
	}

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// check there two outputs for the virtual network resource ids
	vnri, err := terraform.OutputMapE(t, terraformOptions, "virtual_network_resource_ids")
	require.NoErrorf(t, err, "could not get virtual_network_resource_ids output, %s", err)
	assert.Lenf(t, vnri, 2, "expected 2 virtual networks, got %d", len(vnri))
}

// TestDeployVirtualNetworkValidVnetPeering tests the deployment of a virtual network
// with bidirectional peering to a hub virtual network.
func TestDeployVirtualNetworkValidVnetPeering(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["secondary"]
	primaryvnet["hub_peering_enabled"] = true
	secondaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_use_remote_gateways"] = false
	secondaryvnet["hub_peering_use_remote_gateways"] = false

	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	numres := 14
	require.Lenf(t, plan.ResourcePlannedValuesMap, numres, "expected %d resources, got %d", numres, len(plan.ResourcePlannedValuesMap))

	resources := []string{
		"module.virtualnetwork_test.azapi_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vnet[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_inbound[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_inbound[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_outbound[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_outbound[\"secondary\"]",
		"module.virtualnetwork_test.azapi_update_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_update_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		terraform.RequirePlannedValuesMapKeyExists(t, plan, r)
	}

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)
}

// TestDeployVirtualNetworkValidVhubConnection tests the deployment of a virtual network
// with a virtual WAN connection.
func TestDeployVirtualNetworkValidVhubConnection(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["secondary"]
	primaryvnet["vwan_connection_enabled"] = true
	secondaryvnet["vwan_connection_enabled"] = true
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	numres := 13
	require.Lenf(t, plan.ResourcePlannedValuesMap, numres, "expected %d resources, got %d", numres, len(plan.ResourcePlannedValuesMap))

	resources := []string{
		"module.virtualnetwork_test.azapi_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vnet[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.vhubconnection[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vhubconnection[\"secondary\"]",
		"module.virtualnetwork_test.azapi_update_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_update_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		terraform.RequirePlannedValuesMapKeyExists(t, plan, r)
	}

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	// Vhubs cannot be destroyed whilst the routing service is still provisioning
	// hence extended delay
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 5*time.Minute, 5)
}

// TestDeployVirtualNetworkSubnetIdempotency tests that we can make changes
// to the subnet configuration outside the module and that subsequent runs of terraform apply
// are idempotent. See main.tf file in the testdata directory for more details.
func TestDeployVirtualNetworkSubnetIdempotency(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	require.NoError(t, err)

	// test an update to vnet address space, then check for subnet still existing
	primaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["address_space"] = []string{"192.168.0.0/23"}
	_, err = terraform.PlanE(t, terraformOptions)
	require.NoError(t, err)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)
	name := primaryvnet["name"].(string)
	subnets, err := azureutils.ListSubnets(name, name, uuid.MustParse(os.Getenv("AZURE_SUBSCRIPTION_ID")))
	require.NoErrorf(t, err, "failed to list subnets")
	assert.Lenf(t, subnets, 1, "expected 1 subnet, got %d", len(subnets))
}

// TestDeployVirtualNetworkValidMeshPeering tests the deployment of virtual networks
// with mesh peering enables.
func TestDeployVirtualNetworkValidMeshPeering(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]interface{})["secondary"]
	primaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_enabled"] = true
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoError(t, err)
	numres := 10
	require.Lenf(t, plan.ResourcePlannedValuesMap, numres, "expected %d resources, got %d", numres, len(plan.ResourcePlannedValuesMap))

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.peering_mesh[\"primary-secondary\"]",
		"azapi_resource.peering_mesh[\"secondary-primary\"]",
	}
	for _, r := range resources {
		terraform.RequirePlannedValuesMapKeyExists(t, plan, r)
	}

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	// Vhubs cannot be destroyed whilst the routing service is still provisioning
	// hence extended delay
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 5*time.Second, 6)
}

func getValidInputVariables() (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	name2 := name + "-2"

	return map[string]interface{}{
		"subscription_id": os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"virtual_networks": map[string]map[string]interface{}{
			"primary": {
				"name":                name,
				"address_space":       []string{"192.168.0.0/24"},
				"location":            "westeurope",
				"resource_group_name": name,
			},
			"secondary": {
				"name":                name2,
				"address_space":       []string{"192.168.1.0/24"},
				"location":            "northeurope",
				"resource_group_name": name2,
			},
		},
	}, nil
}
