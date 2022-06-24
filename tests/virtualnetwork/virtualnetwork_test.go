package virtualnetwork

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/models"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/virtualnetwork"
)

// TestVirtualNetworkCreateValid tests the creation of a plan that
// creates a virtual network in the specified resource group.
func TestVirtualNetworkCreateValid(t *testing.T) {
	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, "")
	defer utils.RemoveTestDir(t, filepath.Dir(tmp))
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	// Create plan and ensure only two resources are created.

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 2, len(plan.ResourcePlannedValuesMap))
	rg := plan.ResourcePlannedValuesMap["azapi_resource.rg"]
	vnet := plan.ResourcePlannedValuesMap["azapi_resource.vnet"]
	assert.Equal(t, v["virtual_network_resource_group_name"].(string), rg.AttributeValues["name"])
	assert.Equal(t, v["virtual_network_name"].(string), vnet.AttributeValues["name"])
	var vnb models.VirtualNetworkBody
	require.Contains(t, vnet.AttributeValues, "body")
	err = json.Unmarshal([]byte(vnet.AttributeValues["body"].(string)), &vnb)
	require.NoErrorf(t, err, "Could not unmarshal virtual network body")
	assert.Equal(t, v["virtual_network_address_space"], vnb.Properties.AddressSpace.AddressPrefixes)
}

// TestVirtualNetworkCreateValidWithPeering tests the creation of a plan that
// creates a virtual network in the specified resource group.
func TestVirtualNetworkCreateValidWithPeering(t *testing.T) {
	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, "")
	defer utils.RemoveTestDir(t, filepath.Dir(tmp))
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	// Create plan and ensure only two resources are created.
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 4, len(plan.ResourcePlannedValuesMap))

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
	assert.True(t, *body.Properties.AllowGatewayTransit)

	res = "azapi_resource.peering[\"inbound\"]"
	require.Contains(t, plan.ResourcePlannedValuesMap, res)
	vnp = plan.ResourcePlannedValuesMap[res]
	require.Contains(t, vnp.AttributeValues, "parent_id")
	assert.Equal(t, v["hub_network_resource_id"], vnp.AttributeValues["parent_id"])
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_id":                     "00000000-0000-0000-0000-000000000000",
		"virtual_network_address_space":       []string{"10.0.0.0", "172.16.0.0"},
		"virtual_network_location":            "northeurope",
		"virtual_network_name":                "testvnet",
		"virtual_network_resource_group_name": "testrg",
	}
}
