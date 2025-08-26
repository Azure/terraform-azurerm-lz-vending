package virtualnetwork

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
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

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// Loop through each virtual network and check the values
	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		vnetres := fmt.Sprintf("module.virtual_networks[\"%s\"].azapi_resource.vnet", k)
		check.InPlan(test.PlanStruct).That(vnetres).Key("name").HasValue(v["name"]).ErrorIsNil(t)
		check.InPlan(test.PlanStruct).That(vnetres).Key("location").HasValue(v["location"]).ErrorIsNil(t)
		check.InPlan(test.PlanStruct).That(vnetres).Key("body").Query("properties.addressSpace.addressPrefixes").HasValue(v["address_space"]).ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValid tests the creation of a plan that
// creates two virtual networks in the specified resource groups with custom DNS servers.
func TestVirtualNetworkCreateValidWithCustomDns(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["dns_servers"] = []any{"1.2.3.4", "4.3.2.1"}
	secondaryvnet["dns_servers"] = []any{
		"8.8.8.8",
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// want 6 resources, like TestVirtualNetworkCreateValid
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(2).ErrorIsNilFatal(t)

	// Loop through each virtual network and check the values
	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"].azapi_resource.vnet", k)
		check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.dhcpOptions.dnsServers").HasValue(v["dns_servers"]).ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithTags tests the creation of a plan that
// creates two virtual networks in the specified resource groups with tags on vnet and rg.
func TestVirtualNetworkCreateValidWithTags(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["tags"] = map[string]any{
		"tag1": "value1",
		"tag2": "2",
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, same as TestVirtualNetworkCreateValid test
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(2).ErrorIsNilFatal(t)

	check.InPlan(test.PlanStruct).That("module.virtual_networks[\"primary\"].azapi_resource.vnet").Key("tags").HasValue(primaryvnet["tags"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering.
func TestVirtualNetworkCreateValidWithMeshPeering(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_allow_forwarded_traffic"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(4).ErrorIsNilFatal(t)

	peer1 := "module.peering_mesh[\"primary-secondary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowForwardedTraffic").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
	peer1Remote := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/secondary-rg/providers/Microsoft.Network/virtualNetworks/secondary-vnet"
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(peer1Remote).ErrorIsNil(t)

	peer2 := "module.peering_mesh[\"secondary-primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowForwardedTraffic").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
	peer2Remote := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/virtualNetworks/primary-vnet"
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(peer2Remote).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidInvalidMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering
// enabled on only one of the two vnets.
func TestVirtualNetworkCreateValidInvalidMeshPeering(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["mesh_peering_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, as only one of the two vnets has mesh peering enabled, then no peerings should be created
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidSameRg tests the creation of a plan that
// creates two virtual networks in the same resource group.
func TestVirtualNetworkCreateValidSameRg(t *testing.T) {

	v := getMockInputVariables()

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 2 resources here, as the two vnets have the same rg
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidSameRgSameLocation tests the creation of a plan that
// creates two virtual networks in the same resource group in the same location.
func TestVirtualNetworkCreateValidSameRgSameLocation(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["location"] = "northeurope"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 2 resources here, the two vnets have the same rg and same location
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithSubnet tests the creation of a plan that
// creates a virtual network with a subnet.
func TestVirtualNetworkCreateValidSubnet(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("name").HasValue(subnetVal["name"]).ErrorIsNil(t)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.addressPrefixes").HasValue(subnetVal["address_prefixes"]).ErrorIsNil(t)
			}
		}
	}
}

// TestVirtualNetworkCreateSubnetZeroLengthAddressPrefixes tests the length of address_space > 0
func TestVirtualNetworkCreateSubnetZeroLengthAddressPrefixes(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []string{},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "At least 1 subnet address prefix must be specified")
}

// TestVirtualNetworkCreateValidWithMultiplSubnets tests the creation of a plan that
// creates a virtual network with a single subnet in each.
func TestVirtualNetworkCreateValidWithMultiplSubnets(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("name").HasValue(subnetVal["name"]).ErrorIsNil(t)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.addressPrefixes").HasValue(subnetVal["address_prefixes"]).ErrorIsNil(t)
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithMultiplSubnetsInSingleVnet tests the creation of a plan that
// creates a virtual network with multiple subnets and another virtual network without subnets.
func TestVirtualNetworkCreateValidWithMultiplSubnetsInSingleVnet(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
		},
		"privateendpoint": {
			"name":             "snet-privateendpoint",
			"address_prefixes": []any{"192.168.0.64/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"privateendpoint\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("name").HasValue(subnetVal["name"]).ErrorIsNil(t)
				check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.addressPrefixes").HasValue(subnetVal["address_prefixes"]).ErrorIsNil(t)
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetNatGateway tests the creation of a plan that
// creates a virtual network with a subnet and a nat gateway.
func TestVirtualNetworkCreateValidWithSubnetNatGateway(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"nat_gateway": map[string]any{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/natGateways/testvnatgw",
			},
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				natGw, ngExists := subnetVal["nat_gateway"].(map[string]any)
				if ngExists == false || natGw == nil || reflect.ValueOf(natGw).IsNil() {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.natGateway").DoesNotExist().ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.natGateway.id").HasValue(subnetVal["nat_gateway"].(map[string]any)["id"].(string)).ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateSubnetInvalidNetworkSecurityGroup test the resource id value is correct
func TestVirtualNetworkCreateSubnetInvalidNatGateway(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []string{"192.168.0.0/26"},
			"nat_gateway": map[string]any{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/routeTable/testrt",
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Nat Gateway resource id must be valid")
}

// TestVirtualNetworkCreateValidWithSubnetNetworkSecurityGroup tests the creation of a plan that
// creates a virtual network with a subnet and a network security group.
func TestVirtualNetworkCreateValidWithSubnetNetworkSecurityGroup(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"network_security_group": map[string]any{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/testvnsg",
			},
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetKey := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetKey)
				nsg, nsgExists := subnetVal["network_security_group"].(map[string]any)
				if nsgExists == false || nsg == nil || reflect.ValueOf(nsg).IsNil() {
					check.InPlan(test.PlanStruct).That(subnetKey).Key("body").Query("properties.networkSecurityGroup").DoesNotExist().ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetKey).Key("body").Query("properties.networkSecurityGroup.id").HasValue(subnetVal["network_security_group"].(map[string]any)["id"].(string)).ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateSubnetInvalidNetworkSecurityGroup test the resource id value is correct
func TestVirtualNetworkCreateSubnetInvalidNetworkSecurityGroup(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []string{"192.168.0.0/26"},
			"network_security_group": map[string]any{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/routeTable/testrt",
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Network security group resource id must be valid")
}

// TestVirtualNetworkCreateValidWithSubnetPrivateEndpointNetworkPolicy tests the creation of a plan that
// creates a virtual network with a subnet and a private endpoint network policy enabled and disabled.
func TestVirtualNetworkCreateValidWithSubnetPrivateEndpointNetworkPolicy(t *testing.T) {
	// run bith scenarios here

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":                              "snet-default",
			"address_prefixes":                  []any{"192.168.0.0/26"},
			"private_endpoint_network_policies": "Disabled",
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				pe, peExists := subnetVal["private_endpoint_network_policies"]
				if peExists == false {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.privateEndpointNetworkPolicies").HasValue("Enabled").ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.privateEndpointNetworkPolicies").HasValue(pe).ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetPrivateLinkServiceNetworkPolicy tests the creation of a plan that
// creates a virtual network with a subnet and a private link service network policy enabled and disabled.
func TestVirtualNetworkCreateValidWithSubnetPrivateLinkServiceNetworkPolicy(t *testing.T) {
	// run bith scenarios here

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"private_link_service_network_policies_enabled": false,
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				_, plsExists := subnetVal["private_link_service_network_policies_enabled"]
				if plsExists == false {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.privateLinkServiceNetworkPolicies").HasValue("Enabled").ErrorIsNil(t)
				} else {
					if subnetVal["private_link_service_network_policies_enabled"] == true {
						check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.privateLinkServiceNetworkPolicies").HasValue("Enabled").ErrorIsNil(t)
					} else {
						check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.privateLinkServiceNetworkPolicies").HasValue("Disabled").ErrorIsNil(t)
					}
				}
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetRouteTable tests the creation of a plan that
// creates a virtual network with a subnet and a route table associated with it.
func TestVirtualNetworkCreateValidWithSubnetRouteTable(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"route_table": map[string]any{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/routeTables/testrt",
			},
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				rt, rtExists := subnetVal["route_table"].(map[string]any)
				if rtExists == false || rt == nil || reflect.ValueOf(rt).IsNil() {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.routeTable").DoesNotExist().ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.routeTable.id").HasValue(subnetVal["route_table"].(map[string]any)["id"].(string)).ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetDefaultOutboundAccess tests the creation of a plan that
// creates a virtual network with a subnet with default outbound access enabled and disabled.
func TestVirtualNetworkCreateValidWithSubnetDefaultOutboundAccess(t *testing.T) {
	// run bith scenarios here

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":                            "snet-default",
			"address_prefixes":                []any{"192.168.0.0/26"},
			"default_outbound_access_enabled": true,
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				doa, doaExists := subnetVal["default_outbound_access_enabled"]
				if doaExists == false {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.defaultOutboundAccess").HasValue(false).ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.defaultOutboundAccess").HasValue(doa).ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetSingleServiceEndpoint tests the creation of a plan that
// creates a virtual network with a subnet and a single service endpoint assigned.
func TestVirtualNetworkCreateValidWithSubnetSingleServiceEndpoint(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":              "snet-default",
			"address_prefixes":  []any{"192.168.0.0/26"},
			"service_endpoints": []any{"Microsoft.Storage"},
		},
	}
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.1.0/26"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].module.subnet[\"default\"].azapi_resource.subnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("module.virtual_networks[\"%s\"]", k)
		subnets, exists := v["subnets"].(map[string]map[string]any)
		if exists {
			for subnetKey, subnetVal := range subnets {
				subnetRes := fmt.Sprintf("%s.module.subnet[\"%s\"].azapi_resource.subnet", res, subnetKey)
				terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, subnetRes)
				_, seExists := subnetVal["service_endpoints"].([]any)
				if seExists {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.serviceEndpoints.0.service").HasValue("Microsoft.Storage").ErrorIsNil(t)
				} else {
					check.InPlan(test.PlanStruct).That(subnetRes).Key("body").Query("properties.serviceEndpoints").DoesNotExist().ErrorIsNil(t)
				}
			}
		}
	}
}

// TestVirtualNetworkCreateValidWithSubnetMultipleServiceEndpoints tests the creation of a plan that
// creates a virtual network with a subnet and multiple service endpoint assigned.
func TestVirtualNetworkCreateValidWithSubnetMultipleServiceEndpoints(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":              "snet-default",
			"address_prefixes":  []any{"192.168.0.0/26"},
			"service_endpoints": []any{"Microsoft.Storage", "Microsoft.KeyVault"},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.serviceEndpoints.1.service").HasValue("Microsoft.Storage").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.serviceEndpoints.0.service").HasValue("Microsoft.KeyVault").ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithSubnetSingleServiceEndpointPolicy tests the creation of a plan that
// creates a virtual network with a subnet and a single service endpoint policy assigned.
func TestVirtualNetworkCreateValidWithSubnetSingleServiceEndpointPolicy(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	sepResID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/serviceEndpointPolicies/testsep"
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"service_endpoint_policies": map[string]map[string]any{
				"policy1": {
					"id": sepResID,
				},
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.serviceEndpointPolicies.0.id").HasValue(sepResID).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithSubnetMultipleServiceEndpointPolicies tests the creation of a plan that
// creates a virtual network with a subnet and multiple service endpoint policies assigned.
func TestVirtualNetworkCreateValidWithSubnetMultipleServiceEndpointPolicies(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	storageSepResID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/serviceEndpointPolicies/testsepsto"
	kvSepResID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/serviceEndpointPolicies/testsepkv"
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"service_endpoint_policies": map[string]map[string]any{
				"policy1": {
					"id": storageSepResID,
				},
				"policy2": {
					"id": kvSepResID,
				},
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.serviceEndpointPolicies.0.id").HasValue(storageSepResID).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.serviceEndpointPolicies.1.id").HasValue(kvSepResID).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithSubnetSingleDelegation tests the creation of a plan that
// creates a virtual network with a subnet and a single delegation.
func TestVirtualNetworkCreateValidWithSubnetSingleDelegation(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	expectedDelegations := []map[string]any{}
	expectedDelegations = append(expectedDelegations, map[string]any{
		"name": "Microsoft.Web/serverFarms",
		"service_delegation": map[string]any{
			"name": "Microsoft.Web/serverFarms",
		},
	})
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"delegation":       expectedDelegations,
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.0.name").HasValue("Microsoft.Web/serverFarms").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.0.properties.serviceName").HasValue("Microsoft.Web/serverFarms").ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithSubnetMultipleDelegations tests the creation of a plan that
// creates a virtual network with a subnet and multiple delegations.
func TestVirtualNetworkCreateValidWithSubnetMultipleDelegations(t *testing.T) {

	v := getMockInputVariables()

	// Enable primary vnet subnet in test mock input variables
	expectedDelegations := []map[string]any{}
	expectedDelegations = append(expectedDelegations, map[string]any{
		"name": "Microsoft.Web/serverFarms",
		"service_delegation": map[string]any{
			"name": "Microsoft.Web/serverFarms",
		},
	})
	expectedDelegations = append(expectedDelegations, map[string]any{
		"name": "Microsoft.ContainerInstance/containerGroups",
		"service_delegation": map[string]any{
			"name": "Microsoft.ContainerInstance/containerGroups",
		},
	})
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"default": {
			"name":             "snet-default",
			"address_prefixes": []any{"192.168.0.0/26"},
			"delegation":       expectedDelegations,
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.virtual_networks[\"primary\"].module.subnet[\"default\"].azapi_resource.subnet"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.0.name").HasValue("Microsoft.Web/serverFarms").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.0.properties.serviceName").HasValue("Microsoft.Web/serverFarms").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.1.name").HasValue("Microsoft.ContainerInstance/containerGroups").ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.delegations.1.properties.serviceName").HasValue("Microsoft.ContainerInstance/containerGroups").ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithPeering tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub.
func TestVirtualNetworkCreateValidWithHubPeering(t *testing.T) {

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]",
		"module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// We can only check the body of the outbound peering as the inbound values
	// are not known until apply
	outbound := "module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]"
	terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, outbound)

	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowForwardedTraffic").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.useRemoteGateways").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(primaryvnet["hub_network_resource_id"]).ErrorIsNil(t)

	// More limited checks on the inbound peering
	inbound := "module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(inbound).Key("parent_id").HasValue(primaryvnet["hub_network_resource_id"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithPeeringCustomNames tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithHubPeeringCustomNames(t *testing.T) {

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]",
		"module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// Check outbound peering name
	outbound := "module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(outbound).Key("name").HasValue(primaryvnet["hub_peering_name_tohub"]).ErrorIsNil(t)

	// Check inbound peering name
	inbound := "module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(inbound).Key("name").HasValue(primaryvnet["hub_peering_name_fromhub"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithOnlyToHubPeering tests the creation of a plan that
// creates a virtual network with unidirectional peering to a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithOnlyToHubPeering(t *testing.T) {

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_direction"] = "tohub"
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional one is the outbound peering
	resources := []string{
		"module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithOnlyFromHubPeering tests the creation of a plan that
// creates a virtual network with unidirectional peering from a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithOnlyFromHubPeering(t *testing.T) {

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_direction"] = "fromhub"
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional one is the inbound peering
	resources := []string{
		"module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithPeeringUseRemoteGatewaysDisabled
// tests the creation of a plan that configured the outbound peering
// with useRemoteGateways disabled.
func TestVirtualNetworkCreateValidWithPeeringUseCustomOptions(t *testing.T) {

	v := getMockInputVariables()
	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_options_tohub"] = map[string]any{
		"allow_forwarded_traffic":      false,
		"allow_virtual_network_access": false,
		"allow_gateway_transit":        false,
		"use_remote_gateways":          false,
	}
	primaryvnet["hub_peering_options_fromhub"] = map[string]any{
		"allow_forwarded_traffic":      true,
		"allow_virtual_network_access": true,
		"allow_gateway_transit":        true,
		"use_remote_gateways":          true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 4 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]",
		"module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	res := "module.peering_hub_inbound[\"primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.useRemoteGateways").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowForwardedTraffic").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowGatewayTransit").HasValue(true).ErrorIsNil(t)

	res = "module.peering_hub_outbound[\"primary\"].azapi_resource.this[0]"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowForwardedTraffic").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhub tests the creation of a plan that
// creates a virtual network with a vhub connection.
func TestVirtualNetworkCreateValidWithVhub(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true

	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	secondaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	secondaryvnet["vwan_connection_enabled"] = true
	secondaryvnet["vwan_security_configuration"] = map[string]any{
		"routing_intent_enabled": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 4 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection and the other is routing intent
	resources := []string{
		"azapi_resource.vhubconnection_routing_intent[\"secondary\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("parent_id").HasValue(primaryvnet["vwan_hub_resource_id"]).ErrorIsNil(t)

	drt := primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.associatedRouteTable.id").HasValue(drt).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.labels").HasValue([]any{"default"}).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.#").HasValue(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").HasValue(drt).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubCustomRouting tests the creation of a plan that
// creates a virtual network with a vhub connection with custom routing.
func TestVirtualNetworkCreateValidWithVhubCustomRouting(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_propagated_routetables_labels"] = []any{"testlabel", "testlabel2"}
	primaryvnet["vwan_propagated_routetables_resource_ids"] = []any{
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable",
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable2",
	}
	primaryvnet["vwan_associated_routetable_resource_id"] = primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable3"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.vhubconnection[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("parent_id").HasValue(primaryvnet["vwan_hub_resource_id"]).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.associatedRouteTable.id").HasValue(primaryvnet["vwan_associated_routetable_resource_id"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.labels").HasValue(primaryvnet["vwan_propagated_routetables_labels"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.#.id").HasValue(primaryvnet["vwan_propagated_routetables_resource_ids"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecureInternetTraffic tests that secure_internet_traffic == true
func TestVirtualNetworkCreateValidWithVhubSecureInternetTraffic(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_internet_traffic": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.vhubconnection[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.enableInternetSecurity").HasValue(true).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecurePrivateTraffic that managed vnets propagate to "noneRouteTable" with labels "none"
func TestVirtualNetworkCreateValidWithVhubSecurePrivateTraffic(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_private_traffic": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.vhubconnection[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.associatedRouteTable.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.labels").
		HasValue([]any{"none"}).
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/noneRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.#").
		HasValue(1).
		ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecureInternetAndPrivateTraffic tests secure_internet_traffic == true
// and that managed vnets propagate to "noneRouteTable" with labels "none"
func TestVirtualNetworkCreateValidWithVhubSecureInternetAndPrivateTraffic(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_internet_traffic": true,
		"secure_private_traffic":  true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 7 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.vhubconnection[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.enableInternetSecurity").HasValue(true).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.associatedRouteTable.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.labels").
		HasValue([]any{"none"}).
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/noneRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.#").
		HasValue(1).
		ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubRoutingIntentEnabled tests that routingConfiguration is null when
// routing intent is enabled
func TestVirtualNetworkCreateValidWithVhubRoutingIntentEnabled(t *testing.T) {

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"routing_intent_enabled": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 3 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.vhubconnection_routing_intent[\"primary\"]",
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration").DoesNotExist().ErrorIsNil(t)
}

// TestVirtualNetworkCreateInvalidHubNetResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidHubNetResId(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/testrg/providers/Microsoft.Network/virtualNetworks/tes.-tvnet2"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Hub network resource id must be an Azure virtual network resource id")
}

// TestVirtualNetworkCreateInvalidVhubResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidVhubResId(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_hub_resource_id"] = "/subscription/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "vWAN hub resource id must be an Azure vWAN hub network resource id")
}

// TestVirtualNetworkCreateZeroLengthAddressSpace tests the length of address_space > 0
func TestVirtualNetworkCreateZeroLengthAddressSpace(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["address_space"] = []string{}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "At least 1 address space must be specified")
}

// TestVirtualNetworkCreateInvalidAddressSpace tests a valid CIDR address space is used
func TestVirtualNetworkCreateInvalidAddressSpace(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["address_space"] = []string{"10.37.242/35"}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Address space entries must be specified in IPv4 or IPv6 CIDR notation")
}

func TestVirtualNetworkDdosProtection(t *testing.T) {

	// We want 6 resources here
	resources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
		"module.virtual_networks[\"primary\"].data.azurerm_client_config.this",
		"module.virtual_networks[\"secondary\"].azapi_resource.vnet",
		"module.virtual_networks[\"secondary\"].data.azurerm_client_config.this",
	}

	vnetresources := []string{
		"module.virtual_networks[\"primary\"].azapi_resource.vnet",
	}

	t.Run("DdosEnabled", func(t *testing.T) {
		v := getMockInputVariables()
		primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
		primaryvnet["ddos_protection_enabled"] = true
		primaryvnet["ddos_protection_plan_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/ddosProtectionPlans/test-ddos-plan"

		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		require.NoError(t, err)

		check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)
		for _, r := range vnetresources {
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.enableDdosProtection").HasValue(true).ErrorIsNil(t)
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.ddosProtectionPlan.id").HasValue(primaryvnet["ddos_protection_plan_id"]).ErrorIsNil(t)
		}
	})

	t.Run("DdosDisabled", func(t *testing.T) {
		v := getMockInputVariables()

		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		require.NoError(t, err)

		check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)
		for _, r := range vnetresources {
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.enableDdosProtection").HasValue(false).ErrorIsNil(t)
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.ddosProtectionPlan").DoesNotExist().ErrorIsNil(t)
		}
	})
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]any {
	return map[string]any{
		"subscription_id":  "00000000-0000-0000-0000-000000000000",
		"enable_telemetry": false,
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":                "primary-vnet",
				"address_space":       []any{"192.168.0.0/24"},
				"location":            "westeurope",
				"resource_group_name": "primary-rg",
			},
			"secondary": {
				"name":                "secondary-vnet",
				"address_space":       []any{"192.168.1.0/24"},
				"location":            "northeurope",
				"resource_group_name": "secondary-rg",
			},
		},
	}
}
