package integration

// Integration tests apply to the root module and test the deployment of resources
// in specific scenarios.

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../"
)

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpoke(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`azapi_resource.telemetry_root[0]`,
		`module.subscription[0].azapi_resource_action.subscription_cancel[0]`,
		`module.subscription[0].azapi_resource_action.subscription_rename[0]`,
		`module.subscription[0].azapi_resource.subscription[0]`,
		`module.subscription[0].azapi_update_resource.subscription_tags[0]`,
		`module.subscription[0].time_sleep.wait_for_subscription_before_subscription_operations[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.resourcegroup["secondary"].azapi_resource.rg`,
		`module.virtualnetwork[0].module.peering_hub_inbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.peering_hub_outbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azapi_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000305").ErrorIsNil(t)
}

// TestIntegrationVirtualNetworkMissingResourceGroupReference ensures validation fails when
// neither resource_group_key nor resource_group_name_existing is provided for a virtual network.
func TestIntegrationVirtualNetworkMissingResourceGroupReference(t *testing.T) {

	v := getMockInputVariables()
	// Disable the virtual network submodule so locals and module inputs are not evaluated,
	// ensuring variable validation triggers first.
	v["virtual_network_enabled"] = true
	v["resource_group_creation_enabled"] = true
	// Provide a dummy subscription ID so provider init doesn't fail before var validation
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["subscription_alias_enabled"] = false

	// Remove both RG reference fields from the vnet to trigger validation
	vnets := v["virtual_networks"].(map[string]map[string]any)
	vn := vnets["primary"]
	delete(vn, "resource_group_key")
	delete(vn, "resource_group_name_existing")
	vnets["primary"] = vn

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Each virtual network must specify either 'resource_group_key' or")
	assert.ErrorContains(t, err, "'resource_group_name_existing'")
}

// TestIntegrationVwan tests the resource plan when creating a new subscription,
// with a new virtual network and vwan connection to a supplied vhub.
// RG resource lock is disabled
func TestIntegrationVwan(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualHubs/testhub"
	primaryvnet["vwan_connection_enabled"] = true
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`azapi_resource.telemetry_root[0]`,
		`module.subscription[0].azapi_resource_action.subscription_cancel[0]`,
		`module.subscription[0].azapi_resource_action.subscription_rename[0]`,
		`module.subscription[0].azapi_resource.subscription[0]`,
		`module.subscription[0].azapi_update_resource.subscription_tags[0]`,
		`module.subscription[0].time_sleep.wait_for_subscription_before_subscription_operations[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.resourcegroup["secondary"].azapi_resource.rg`,
		`module.virtualnetwork[0].azapi_resource.vhubconnection["primary"]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azapi_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000505").ErrorIsNil(t)
}

// TestIntegrationSubscriptionAndRoleAssignmentOnly tests the resource plan when creating a new subscription,
// with a role assignments, but no networking.
// This tests that the depends_on property of the roleassignments module is working
// when a dependent resource is disabled through the use of count.
func TestIntegrationSubscriptionAndRoleAssignmentOnly(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = false
	v["role_assignment_enabled"] = true
	v["role_assignments"] = map[string]any{
		"ra": map[string]any{
			"principal_id":   "00000000-0000-0000-0000-000000000000",
			"definition":     "Owner",
			"relative_scope": "",
		},
	}
	v["resource_group_creation_enabled"] = false
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`azapi_resource.telemetry_root[0]`,
		`module.roleassignment["ra"].data.azapi_resource_list.role_definitions[0]`,
		`module.roleassignment["ra"].azapi_resource.this`,
		`module.subscription[0].azapi_resource.subscription[0]`,
		`module.subscription[0].azapi_resource_action.subscription_cancel[0]`,
		`module.subscription[0].azapi_resource_action.subscription_rename[0]`,
		`module.subscription[0].azapi_update_resource.subscription_tags[0]`,
		`module.subscription[0].time_sleep.wait_for_subscription_before_subscription_operations[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
}

// TestIntegrationHubAndSpokeExistingSubscription tests the resource plan when supplying an existing subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpokeExistingSubscription(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]

	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	primaryvnet["hub_peering_enabled"] = true
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["virtual_network_enabled"] = true
	delete(v, "subscription_tags")
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`azapi_resource.telemetry_root[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.resourcegroup["secondary"].azapi_resource.rg`,
		`module.virtualnetwork[0].module.peering_hub_inbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.peering_hub_outbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azapi_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000300").ErrorIsNil(t)
}

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationDisableTelemetry(t *testing.T) {

	v := getMockInputVariables()
	v["subscription_alias_enabled"] = true
	v["disable_telemetry"] = true
	v["resource_group_creation_enabled"] = false
	v["virtual_network_enabled"] = false

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`module.subscription[0].azapi_resource.subscription[0]`,
		`module.subscription[0].azapi_resource_action.subscription_cancel[0]`,
		`module.subscription[0].azapi_resource_action.subscription_rename[0]`,
		`module.subscription[0].azapi_update_resource.subscription_tags[0]`,
		`module.subscription[0].time_sleep.wait_for_subscription_before_subscription_operations[0]`,
	}
	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
}

func TestIntegrationResourceGroups(t *testing.T) {

	v := map[string]any{
		"subscription_id":                 "00000000-0000-0000-0000-000000000000",
		"location":                        "westeurope",
		"resource_group_creation_enabled": true,
		"disable_telemetry":               true,
		"resource_groups": map[string]any{
			"NetworkWatcherRG": map[string]any{
				"location": "westeurope",
				"name":     "NetworkWatcherRG",
			},
			"rg1": map[string]any{
				"location": "westeurope",
				"name":     "rg1",
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`module.resourcegroup["rg1"].azapi_resource.rg`,
		`module.resourcegroup["NetworkWatcherRG"].azapi_resource.rg`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
}

func TestIntegrationUmiRoleAssignment(t *testing.T) {

	v := map[string]any{
		"subscription_id":                 "00000000-0000-0000-0000-000000000000",
		"location":                        "westeurope",
		"disable_telemetry":               true,
		"resource_group_creation_enabled": true,
		"resource_groups": map[string]any{
			"primary": map[string]any{
				"location": "westeurope",
				"name":     "rg-umi",
			},
		},
		"umi_enabled": true,
		"user_managed_identities": map[string]any{
			"default": map[string]any{
				"name":               "umi",
				"resource_group_key": "primary",
				"role_assignments": map[string]any{
					"owner": map[string]any{
						"definition":     "Owner",
						"relative_scope": "",
					},
				},
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`time_sleep.wait_for_umi_before_umi_role_assignment_operations[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.usermanagedidentity["default"].azapi_resource.umi`,
		`module.roleassignment_umi["default/owner"].azapi_resource.this`,
		`module.roleassignment_umi["default/owner"].data.azapi_resource_list.role_definitions[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
}

func TestIntegrationMultipleUmiRoleAssignments(t *testing.T) {

	v := map[string]any{
		"subscription_id":                 "00000000-0000-0000-0000-000000000000",
		"location":                        "westeurope",
		"disable_telemetry":               true,
		"resource_group_creation_enabled": true,
		"resource_groups": map[string]any{
			"primary": map[string]any{
				"location": "westeurope",
				"name":     "rg-umi",
			},
		},
		"umi_enabled": true,
		"user_managed_identities": map[string]any{
			"default": map[string]any{
				"name":                         "umi",
				"resource_group_name_existing": "rg-umi",
				"role_assignments": map[string]any{
					"owner": map[string]any{
						"definition":     "Owner",
						"relative_scope": "",
					},
					"blob": map[string]any{
						"definition":     "Storage Blob Data Owner",
						"relative_scope": "",
					},
				},
			},
			"backup": map[string]any{
				"name":               "umi-backup",
				"resource_group_key": "primary",
				"role_assignments": map[string]any{
					"owner": map[string]any{
						"definition":     "Owner",
						"relative_scope": "",
					},
					"blob": map[string]any{
						"definition":     "Storage Blob Data Owner",
						"relative_scope": "",
					},
				},
			},
		},
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`module.roleassignment_umi["backup/blob"].azapi_resource.this`,
		`module.roleassignment_umi["backup/blob"].data.azapi_resource_list.role_definitions[0]`,
		`module.roleassignment_umi["backup/owner"].azapi_resource.this`,
		`module.roleassignment_umi["backup/owner"].data.azapi_resource_list.role_definitions[0]`,
		`module.roleassignment_umi["default/blob"].azapi_resource.this`,
		`module.roleassignment_umi["default/blob"].data.azapi_resource_list.role_definitions[0]`,
		`module.roleassignment_umi["default/owner"].azapi_resource.this`,
		`module.roleassignment_umi["default/owner"].data.azapi_resource_list.role_definitions[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.usermanagedidentity["backup"].azapi_resource.umi`,
		`module.usermanagedidentity["default"].azapi_resource.umi`,
		`time_sleep.wait_for_umi_before_umi_role_assignment_operations[0]`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
}

// TestIntegrationVirtualNetworkRouteTable tests the resource plan when creating a new subscription,
// with a new virtual network with route table.
func TestIntegrationVirtualNetworkRouteTable(t *testing.T) {

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["subnets"] = map[string]map[string]any{
		"primary": {
			"name":             "primary-subnet",
			"address_prefixes": []string{"192.168.0.0/25"},
			"route_table": map[string]string{
				"key_reference": "primary",
			},
		},
		"secondary": {
			"name":             "secondary-subnet",
			"address_prefixes": []string{"192.168.0.128/25"},
			"route_table": map[string]string{
				"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/routeTables/primary-route-table",
			},
		},
	}
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["route_table_enabled"] = true
	v["route_tables"] = map[string]any{
		"primary": map[string]string{
			"name":               "primary-route-table",
			"resource_group_key": "primary",
			"location":           "westeurope",
		},
		"default": map[string]string{
			"name":                         "default-route-table",
			"resource_group_name_existing": "primary-rg",
			"location":                     "westeurope",
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	resources := []string{
		`azapi_resource.telemetry_root[0]`,
		`module.subscription[0].azapi_resource_action.subscription_cancel[0]`,
		`module.subscription[0].azapi_resource_action.subscription_rename[0]`,
		`module.subscription[0].azapi_resource.subscription[0]`,
		`module.subscription[0].azapi_update_resource.subscription_tags[0]`,
		`module.subscription[0].time_sleep.wait_for_subscription_before_subscription_operations[0]`,
		`module.resourcegroup["primary"].azapi_resource.rg`,
		`module.resourcegroup["secondary"].azapi_resource.rg`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].module.subnet["primary"].azapi_resource.subnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].module.subnet["secondary"].azapi_resource.subnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azapi_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
		`module.routetable["primary"].azapi_resource.route_table`,
		`module.routetable["default"].azapi_resource.route_table`,
	}

	check.InPlan(test.PlanStruct).PlannedResourcesAre(resources...).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(`module.virtualnetwork[0].module.virtual_networks["primary"].module.subnet["secondary"].azapi_resource.subnet`).Key("body").Query("properties.routeTable.id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/routeTables/primary-route-table").ErrorIsNil(t)
}

func getMockInputVariables() map[string]any {
	return map[string]any{
		"location": "northeurope",
		// subscription variables
		"subscription_billing_scope":                            "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_register_resource_providers_and_features": map[string][]any{},
		"subscription_display_name":                             "test-subscription-alias",
		"subscription_alias_name":                               "test-subscription-alias",
		"subscription_workload":                                 "Production",
		"subscription_tags": map[string]any{
			"test-tag":   "test-value",
			"test-tag-2": "test-value-2",
		},

		// resource_group variables
		"resource_group_creation_enabled": true,
		"resource_groups": map[string]map[string]any{
			"primary": {
				"name":     "primary-rg",
				"location": "westeurope",
			},
			"secondary": {
				"name":     "secondary-rg",
				"location": "westeurope",
			},
		},

		// virtualnetwork variables
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":               "primary-vnet",
				"address_space":      []string{"192.168.0.0/24"},
				"location":           "westeurope",
				"resource_group_key": "primary",
			},
		},
	}
}
