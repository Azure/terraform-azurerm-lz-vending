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
// creates a virtual network in the specified resource group.
func TestVirtualNetworkCreateValid(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["virtual_network_resource_lock_enabled"] = true
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equalf(t, 4, len(plan.ResourcePlannedValuesMap), "expected 4 resources to be created, got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"azapi_resource.rg",
		"azapi_resource.vnet",
		"azapi_update_resource.vnet",
		"azapi_resource.rg_lock[0]",
	}
	for _, r := range resources {
		require.Contains(t, plan.ResourcePlannedValuesMap, r)
	}

	rg := plan.ResourcePlannedValuesMap["azapi_resource.rg"]
	vnet := plan.ResourcePlannedValuesMap["azapi_resource.vnet"]
	require.Contains(t, rg.AttributeValues, "name")
	assert.Equal(t, v["virtual_network_resource_group_name"].(string), rg.AttributeValues["name"])
	require.Contains(t, vnet.AttributeValues, "name")
	assert.Equal(t, v["virtual_network_name"].(string), vnet.AttributeValues["name"])
	var vnb models.VirtualNetworkBody
	require.Contains(t, vnet.AttributeValues, "body")
	err = json.Unmarshal([]byte(vnet.AttributeValues["body"].(string)), &vnb)
	require.NoErrorf(t, err, "Could not unmarshal virtual network body")
	assert.Equal(t, v["virtual_network_address_space"], vnb.Properties.AddressSpace.AddressPrefixes)
}

// TestVirtualNetworkCreateValidWithPeering tests the creation of a plan that
// creates a virtual network with bidirectional peering.
func TestVirtualNetworkCreateValidWithPeering(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["virtual_network_peering_enabled"] = true
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	terraformOptions.Vars = v
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equalf(t, 5, len(plan.ResourcePlannedValuesMap), "expected 5 resources to be created, got %d", len(plan.ResourcePlannedValuesMap))

	// We can only check the body of the outbound peering as the inbound values
	// not known until apply
	res := "azapi_resource.peering[\"outbound\"]"
	require.Contains(t, plan.ResourcePlannedValuesMap, res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Contains(t, vnp.AttributeValues, "body")
	var body models.VirtualNetworkPeeringBody
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.True(t, *body.Properties.AllowForwardedTraffic)
	assert.True(t, *body.Properties.AllowVirtualNetworkAccess)
	assert.False(t, *body.Properties.AllowGatewayTransit)
	assert.True(t, *body.Properties.UseRemoteGateways)
	assert.Equal(t, body.Properties.RemoteVirtualNetwork.ID, v["hub_network_resource_id"])

	// More limited checks on the inbound peering
	res = "azapi_resource.peering[\"inbound\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp = plan.ResourcePlannedValuesMap[res]
	require.Contains(t, vnp.AttributeValues, "parent_id")
	assert.Equal(t, v["hub_network_resource_id"], vnp.AttributeValues["parent_id"])
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
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/tes.-tvnet2"
	v["virtual_network_peering_enabled"] = true
	v["virtual_network_use_remote_gateways"] = false
	terraformOptions.Vars = v

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equalf(t, 5, len(plan.ResourcePlannedValuesMap), "expected 5 resources to be created, got %d", len(plan.ResourcePlannedValuesMap))

	// We can only check the body of the outbound peering as the inbound values
	// not known until apply
	res := "azapi_resource.peering[\"outbound\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, res)
	vnp := plan.ResourcePlannedValuesMap[res]
	require.Contains(t, vnp.AttributeValues, "body")
	var body models.VirtualNetworkPeeringBody
	err = json.Unmarshal([]byte(vnp.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.False(t, *body.Properties.UseRemoteGateways)
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
	v := getMockInputVariables()
	v["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	v["virtual_network_vwan_connection_enabled"] = true
	terraformOptions.Vars = v
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 4, len(plan.ResourcePlannedValuesMap))
	vhcname := "vhubcon-1b4db7eb-4057-5ddf-91e0-36dec72071f5"
	vhcres := fmt.Sprintf("azapi_resource.vhubconnection[\"%s\"]", vhcname)
	terraform.RequirePlannedValuesMapKeyExists(t, plan, vhcres)
	vhc := plan.ResourcePlannedValuesMap[vhcres]
	require.Contains(t, vhc.AttributeValues, "name")
	assert.Equal(t, vhcname, vhc.AttributeValues["name"])
	require.Contains(t, vhc.AttributeValues, "parent_id")
	assert.Equal(t, v["vwan_hub_resource_id"], vhc.AttributeValues["parent_id"])

	require.Contains(t, vhc.AttributeValues, "body")
	var body models.HubVirtualNetworkConnectionBody
	err = json.Unmarshal([]byte(vhc.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	drt := v["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable"
	assert.Equal(t, drt, body.Properties.RoutingConfiguration.AssociatedRouteTable.ID)
	assert.EqualValues(t, []string{"default"}, body.Properties.RoutingConfiguration.PropagatedRouteTables.Labels)
	assert.Len(t, body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs, 1)
	for _, rt := range body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs {
		assert.Contains(t, drt, rt.ID)
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
	v := getMockInputVariables()
	v["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	v["virtual_network_vwan_connection_enabled"] = true
	v["virtual_network_vwan_propagated_routetables_labels"] = []string{"testlabel", "testlabel2"}
	v["virtual_network_vwan_propagated_routetables_resource_ids"] = []string{
		v["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable",
		v["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable2",
	}
	v["virtual_network_vwan_routetable_resource_id"] = v["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable3"
	terraformOptions.Vars = v
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 4, len(plan.ResourcePlannedValuesMap))

	vhcname := "vhubcon-1b4db7eb-4057-5ddf-91e0-36dec72071f5"
	vhcres := fmt.Sprintf("azapi_resource.vhubconnection[\"%s\"]", vhcname)
	terraform.RequirePlannedValuesMapKeyExists(t, plan, vhcres)
	vhc := plan.ResourcePlannedValuesMap[vhcres]
	require.Contains(t, vhc.AttributeValues, "name")
	assert.Equal(t, vhcname, vhc.AttributeValues["name"])
	require.Contains(t, vhc.AttributeValues, "parent_id")
	assert.Equal(t, v["vwan_hub_resource_id"], vhc.AttributeValues["parent_id"])

	require.Contains(t, vhc.AttributeValues, "body")
	var body models.HubVirtualNetworkConnectionBody
	err = json.Unmarshal([]byte(vhc.AttributeValues["body"].(string)), &body)
	require.NoErrorf(t, err, "Could not unmarshal virtual network peering body")
	assert.Equal(t, v["virtual_network_vwan_routetable_resource_id"], body.Properties.RoutingConfiguration.AssociatedRouteTable.ID)
	assert.EqualValues(t, v["virtual_network_vwan_propagated_routetables_labels"], body.Properties.RoutingConfiguration.PropagatedRouteTables.Labels)
	assert.Len(t, body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs, 2)
	for _, rt := range body.Properties.RoutingConfiguration.PropagatedRouteTables.IDs {
		assert.Contains(t, v["virtual_network_vwan_propagated_routetables_resource_ids"], rt.ID)
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
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/testrg/providers/Microsoft.Network/virtualNetworks/tes.-tvnet2"
	terraformOptions.Vars = v
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "Value must be an Azure virtual network resource id")
}

// TestVirtualNetworkCreateInvalidHubNetResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidVhubResId(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["vwan_hub_resource_id"] = "/subscription/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	terraformOptions.Vars = v
	_, err = terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.ErrorContains(t, err, "Value must be an Azure vwan hub resource id")
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_id":                       "00000000-0000-0000-0000-000000000000",
		"virtual_network_address_space":         []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_location":              "northeurope",
		"virtual_network_name":                  "testvnet",
		"virtual_network_resource_group_name":   "testrg",
		"virtual_network_resource_lock_enabled": false,
	}
}
