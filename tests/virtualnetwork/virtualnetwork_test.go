package virtualnetwork

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/models"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/virtualnetwork"
)

// TestVirtualNetworkCreateValid tests the creation of a plan that
// creates two virtual networks in the specified resource groups.
func TestVirtualNetworkCreateValid(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	terraformOptions.Vars = vars

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
	}
	require.Equalf(t, len(resources), len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", len(resources), len(plan.ResourcePlannedValuesMap))

	for _, r := range resources {
		require.Containsf(t, plan.ResourcePlannedValuesMap, r, "plan does not contain expected resource %s", r)
	}

	// Loop through each virtual network and check the values
	vns := vars["virtual_networks"].(map[string]map[string]interface{})
	for k, v := range vns {
		rg := plan.ResourcePlannedValuesMap[fmt.Sprintf("azapi_resource.rg[\"%s-rg\"]", k)]
		vnet := plan.ResourcePlannedValuesMap[fmt.Sprintf("azapi_resource.vnet[\"%s\"]", k)]

		require.Containsf(t, rg.AttributeValues, "name", "resource group %s does not contain name", k)
		assert.Equal(t, v["resource_group_name"].(string), rg.AttributeValues["name"])

		require.Containsf(t, vnet.AttributeValues, "name", "virtual network %s does not contain name", k)
		assert.Equal(t, v["name"].(string), vnet.AttributeValues["name"])

		require.Containsf(t, rg.AttributeValues, "location", "resource group %s does not contain location", k)
		assert.Equalf(t, v["location"].(string), rg.AttributeValues["location"], "resource group %s location does not match %s", k, v["location"].(string))

		require.Containsf(t, vnet.AttributeValues, "location", "virtual network %s does not contain location", k)
		assert.Equalf(t, v["location"].(string), vnet.AttributeValues["location"], "virtual network %s location does not match %s", k, v["location"].(string))

		var vnb models.VirtualNetworkBody
		require.Containsf(t, vnet.AttributeValues, "body", "virtual network %s does not contain body", k)
		err = json.Unmarshal([]byte(vnet.AttributeValues["body"].(string)), &vnb)
		require.NoErrorf(t, err, "Could not unmarshal virtual network body")
		assert.Equalf(t, v["address_space"], vnb.Properties.AddressSpace.AddressPrefixes, "virtual network %s address space does not match", k)
	}
}

// TestVirtualNetworkCreateValidWithMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering.
func TestVirtualNetworkCreateValidWithMeshPeering(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	secondaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["secondary"]
	primaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_allow_forwarded_traffic"] = true

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 10
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	res := "azapi_resource.peering_mesh[\"primary-secondary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "body", "virtual network peering %s does not contain body", res)
	var body models.VirtualNetworkPeeringBody
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Falsef(t, *body.Properties.AllowForwardedTraffic, "expected allow forwarded traffic to be false")
	assert.Truef(t, *body.Properties.AllowVirtualNetworkAccess, "expected allow virtual network access to be true")
	assert.Falsef(t, *body.Properties.AllowGatewayTransit, "expected allow gateway transit to be false")
	assert.Falsef(t, *body.Properties.UseRemoteGateways, "expected use remote gateways to be true")
	rvnid := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/secondary-rg/providers/Microsoft.Network/virtualNetworks/secondary-vnet"
	assert.Equalf(t, body.Properties.RemoteVirtualNetwork.ID, rvnid, "expected remote virtual network id to be %s", rvnid)

	res = "azapi_resource.peering_mesh[\"secondary-primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp = plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "body", "virtual network peering %s does not contain body", res)
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Truef(t, *body.Properties.AllowForwardedTraffic, "expected allow forwarded traffic to be true")
	assert.Truef(t, *body.Properties.AllowVirtualNetworkAccess, "expected allow virtual network access to be true")
	assert.Falsef(t, *body.Properties.AllowGatewayTransit, "expected allow gateway transit to be false")
	assert.Falsef(t, *body.Properties.UseRemoteGateways, "expected use remote gateways to be true")
	rvnid = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/virtualNetworks/primary-vnet"
	assert.Equalf(t, body.Properties.RemoteVirtualNetwork.ID, rvnid, "expected remote virtual network id to be %s", rvnid)
}

// TestVirtualNetworkCreateValidInvalidMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering
// enabled on only one of the two vnets.
func TestVirtualNetworkCreateValidInvalidMeshPeering(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["mesh_peering_enabled"] = true

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 8 resources here, as only one of the two vnets has mesh peering enabled, then no peerings should be created
	numres := 8
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))
}

// TestVirtualNetworkCreateValidSameRg tests the creation of a plan that
// creates two virtual networks in the same resource group.
func TestVirtualNetworkCreateValidSameRg(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["resource_group_name"] = "secondary-rg"
	primaryvnet["resource_group_creation_enabled"] = false

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 6 resources here, as the two vnets have the same rg, then 2 fewer resources than
	// TestVirtualNetworkCreateValid (rg + rg lock)
	numres := 6
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))
}

// TestVirtualNetworkCreateValidWithPeering tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub.
func TestVirtualNetworkCreateValidWithHubPeering(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 10
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	// We can only check the body of the outbound peering as the inbound values
	// are not known until apply
	res := "azapi_resource.peering_hub_outbound[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "body", "virtual network peering %s does not contain body", res)
	var body models.VirtualNetworkPeeringBody
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Truef(t, *body.Properties.AllowForwardedTraffic, "expected allow forwarded traffic to be true")
	assert.Truef(t, *body.Properties.AllowVirtualNetworkAccess, "expected allow virtual network access to be true")
	assert.Falsef(t, *body.Properties.AllowGatewayTransit, "expected allow gateway transit to be false")
	assert.Truef(t, *body.Properties.UseRemoteGateways, "expected use remote gateways to be true")
	assert.Equalf(t, body.Properties.RemoteVirtualNetwork.ID, primaryvnet["hub_network_resource_id"], "expected remote virtual network id to be %s", primaryvnet["hub_network_resource_id"])

	// More limited checks on the inbound peering
	res = "azapi_resource.peering_hub_inbound[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp = plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "parent_id", "virtual network peering %s does not contain parent_id", res)
	assert.Equalf(t, primaryvnet["hub_network_resource_id"], vnp.AttributeValues["parent_id"], "expected parent_id to be %s", primaryvnet["hub_network_resource_id"])
}

// TestVirtualNetworkCreateValidWithPeeringCustomNames tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithHubPeeringCustomNames(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 10
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	// Check outbound peering name
	res := "azapi_resource.peering_hub_outbound[\"primary\"]"
	require.Containsf(t, plan.ResourcePlannedValuesMap, res, "virtual network peering %s does not exist", res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "name", "virtual network peering %s does not contain name", res)
	assert.Equalf(t, primaryvnet["hub_peering_name_tohub"], vnp.AttributeValues["name"], "expected name to be %s", primaryvnet["hub_peering_name_tohub"])

	// Check inbound peering name
	res = "azapi_resource.peering_hub_inbound[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp = plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "name", "virtual network peering %s does not contain name", res)
	assert.Equalf(t, primaryvnet["hub_peering_name_fromhub"], vnp.AttributeValues["name"], "expected name to be %s", primaryvnet["hub_peering_name_fromhub"])
}

// TestVirtualNetworkCreateValidWithPeeringUseRemoteGatewaysDisabled
// tests the creation of a plan that configured the outbound peering
// with useRemoteGateways disabled.
func TestVirtualNetworkCreateValidWithPeeringUseRemoteGatewaysDisabled(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_use_remote_gateways"] = false
	terraformOptions.Vars = vars

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 10
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	// We can only check the body of the outbound peering as the inbound values
	// not known until apply
	res := "azapi_resource.peering_hub_outbound[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Containsf(t, vnp.AttributeValues, "body", "virtual network peering %s does not contain body", res)
	var body models.VirtualNetworkPeeringBody
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Falsef(t, *body.Properties.UseRemoteGateways, "expected use remote gateways to be false")
}

// TestVirtualNetworkCreateValidWithVhub tests the creation of a plan that
// creates a virtual network with a vhub connection.
func TestVirtualNetworkCreateValidWithVhub(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 9
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, vhcres)
	vhc := plan.ResourcePlannedValuesMap[vhcres]
	require.Contains(t, vhc.AttributeValues, "parent_id")
	assert.Equal(t, primaryvnet["vwan_hub_resource_id"], vhc.AttributeValues["parent_id"])

	require.Contains(t, vhc.AttributeValues, "body")
	var body models.HubVirtualNetworkConnectionBody
	err = json.Unmarshal([]byte(vhc.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	drt := primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable"
	assert.Equalf(t, drt, body.Properties.RoutingConfiguration.AssociatedRouteTable.ID, "expected default route table to be %s", drt)
	assert.EqualValuesf(t, []string{"default"}, body.Properties.RoutingConfiguration.PropagatedRouteTables.Labels, "expected propagated route tables to be %v", []string{"default"})
	assert.Lenf(t, body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs, 1, "expected length of propageted route tables to be 1")
	for _, rt := range body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs {
		assert.Containsf(t, drt, rt.ID, "expected propagated route tables to contain %s", rt.ID)
	}
}

// TestVirtualNetworkCreateValidWithVhubCustomRouting tests the creation of a plan that
// creates a virtual network with a vhub connection with custom routing.
func TestVirtualNetworkCreateValidWithVhubCustomRouting(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_propagated_routetables_labels"] = []string{"testlabel", "testlabel2"}
	primaryvnet["vwan_propagated_routetables_resource_ids"] = []string{
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable",
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable2",
	}
	primaryvnet["vwan_associated_routetable_resource_id"] = primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable3"

	terraformOptions.Vars = vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoErrorf(t, err, "failed to init and plan")

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	numres := 9
	require.Equalf(t, numres, len(plan.ResourcePlannedValuesMap), "expected %d resources to be created, got %d", numres, len(plan.ResourcePlannedValuesMap))

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, vhcres)
	vhc := plan.ResourcePlannedValuesMap[vhcres]
	require.Containsf(t, vhc.AttributeValues, "parent_id", "expected parent_id to be set")
	assert.Equalf(t, primaryvnet["vwan_hub_resource_id"], vhc.AttributeValues["parent_id"], "expected parent_id to be %s", primaryvnet["vwan_hub_resource_id"])

	require.Containsf(t, vhc.AttributeValues, "body", "expected body to be set")
	var body models.HubVirtualNetworkConnectionBody
	err = json.Unmarshal([]byte(vhc.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Equalf(t, primaryvnet["vwan_associated_routetable_resource_id"], body.Properties.RoutingConfiguration.AssociatedRouteTable.ID, "expected associated route table to be %s", primaryvnet["vwan_associated_routetable_resource_id"])
	assert.EqualValuesf(t, primaryvnet["vwan_propagated_routetables_labels"], body.Properties.RoutingConfiguration.PropagatedRouteTables.Labels, "expected propagated route tables to be %v", primaryvnet["vwan_propagated_routetables_labels"])
	assert.Lenf(t, body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs, 2, "expected length of propagated route tables to be 2")
	for _, rt := range body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs {
		assert.Containsf(t, primaryvnet["vwan_propagated_routetables_resource_ids"], rt.ID, "expected propagated route tables to contain %s", rt.ID)
	}
}

// TestVirtualNetworkCreateInvalidHubNetResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidHubNetResId(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/testrg/providers/Microsoft.Network/virtualNetworks/tes.-tvnet2"
	terraformOptions.Vars = vars
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "Hub network resource id must be an Azure virtual network resource id")
}

// TestVirtualNetworkCreateInvalidVhubResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidVhubResId(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscription/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	terraformOptions.Vars = vars
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "vWAN hub resource id must be an Azure vWAN hub network resource id")
}

// TestVirtualNetworkCreateZeroLengthAddressSpace tests the length of address_space > 0
func TestVirtualNetworkCreateZeroLengthAddressSpace(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["address_space"] = []string{}
	terraformOptions.Vars = vars
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "At least 1 address space must be specified")
}

// TestVirtualNetworkCreateInvalidAddressSpace tests the length of address_space > 0
func TestVirtualNetworkCreateInvalidAddressSpace(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	vars := getMockInputVariables()
	primaryvnet := vars["virtual_networks"].(map[string]map[string]interface{})["primary"]
	primaryvnet["address_space"] = []string{"10.37.242/35"}
	terraformOptions.Vars = vars
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "Address space entries must be specified in CIDR notation")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_id": "00000000-0000-0000-0000-000000000000",
		"virtual_networks": map[string]map[string]interface{}{
			"primary": {
				"name":                "primary-vnet",
				"address_space":       []string{"192.168.0.0/24"},
				"location":            "westeurope",
				"resource_group_name": "primary-rg",
			},
			"secondary": {
				"name":                "secondary-vnet",
				"address_space":       []string{"192.168.1.0/24"},
				"location":            "northeurope",
				"resource_group_name": "secondary-rg",
			},
		},
	}
}
