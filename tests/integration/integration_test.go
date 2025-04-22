package integration

// Integration tests apply to the root module and test the deployment of resources
// in specific scenarios.

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../"
)

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpoke(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	primaryvnet["resource_group_lock_enabled"] = true
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
		`module.virtualnetwork[0].azapi_resource.rg_lock["primary-rg"]`,
		`module.virtualnetwork[0].azapi_resource.rg["primary-rg"]`,
		`module.virtualnetwork[0].module.peering_hub_inbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.peering_hub_outbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.this`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000b05").ErrorIsNil(t)
}

// TestIntegrationVwan tests the resource plan when creating a new subscription,
// with a new virtual network and vwan connection to a supplied vhub.
// RG resource lock is disabled
func TestIntegrationVwan(t *testing.T) {
	t.Parallel()

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
		`module.virtualnetwork[0].azapi_resource.rg["primary-rg"]`,
		`module.virtualnetwork[0].azapi_resource.vhubconnection["primary"]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.this`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000505").ErrorIsNil(t)
}

// TestIntegrationSubscriptionAndRoleAssignmentOnly tests the resource plan when creating a new subscription,
// with a role assignments, but no networking.
// This tests that the depends_on property of the roleassignments module is working
// when a dependent resource is disabled through the use of count.
func TestIntegrationSubscriptionAndRoleAssignmentOnly(t *testing.T) {
	t.Parallel()

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

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00010005").ErrorIsNil(t)
}

// TestIntegrationHubAndSpokeExistingSubscription tests the resource plan when supplying an existing subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpokeExistingSubscription(t *testing.T) {
	t.Parallel()

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
		`module.virtualnetwork[0].azapi_resource.rg["primary-rg"]`,
		`module.virtualnetwork[0].module.peering_hub_inbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.peering_hub_outbound["primary"].azapi_resource.this[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].azapi_resource.vnet`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.azurerm_client_config.this`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].data.modtm_module_source.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].modtm_telemetry.telemetry[0]`,
		`module.virtualnetwork[0].module.virtual_networks["primary"].random_uuid.telemetry[0]`,
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}

	check.InPlan(test.PlanStruct).That("azapi_resource.telemetry_root[0]").Key("name").ContainsString("00000300").ErrorIsNil(t)
}

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationDisableTelemetry(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["subscription_alias_enabled"] = true
	v["disable_telemetry"] = true

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

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}
}

func TestIntegrationResourceGroups(t *testing.T) {
	t.Parallel()

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

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}
}

func TestIntegrationUmiRoleAssignment(t *testing.T) {
	t.Parallel()

	v := map[string]any{
		"subscription_id":   "00000000-0000-0000-0000-000000000000",
		"location":          "westeurope",
		"disable_telemetry": true,
		"umi_enabled":       true,
		"user_managed_identities": map[string]any{
			"default": map[string]any{
				"name":                "umi",
				"resource_group_name": "rg-umi",
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
		`module.usermanagedidentity["default"].azapi_resource.rg[0]`,
		`module.usermanagedidentity["default"].azapi_resource.rg_lock[0]`,
		`module.usermanagedidentity["default"].azapi_resource.umi`,
		`module.roleassignment_umi["default/owner"].azapi_resource.this`,
		`module.roleassignment_umi["default/owner"].data.azapi_resource_list.role_definitions[0]`,
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}
}

func TestIntegrationMultipleUmiRoleAssignments(t *testing.T) {
	t.Parallel()

	v := map[string]any{
		"subscription_id":   "00000000-0000-0000-0000-000000000000",
		"location":          "westeurope",
		"disable_telemetry": true,
		"umi_enabled":       true,
		"user_managed_identities": map[string]any{
			"default": map[string]any{
				"name":                "umi",
				"resource_group_name": "rg-umi",
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
				"name":                "umi-backup",
				"resource_group_name": "rg-umi",
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
		`module.usermanagedidentity["backup"].azapi_resource.rg_lock[0]`,
		`module.usermanagedidentity["backup"].azapi_resource.rg[0]`,
		`module.usermanagedidentity["backup"].azapi_resource.umi`,
		`module.usermanagedidentity["default"].azapi_resource.rg_lock[0]`,
		`module.usermanagedidentity["default"].azapi_resource.rg[0]`,
		`module.usermanagedidentity["default"].azapi_resource.umi`,
		`time_sleep.wait_for_umi_before_umi_role_assignment_operations[0]`,
	}

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNil(t)
	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}
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

		// virtualnetwork variables
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":                        "primary-vnet",
				"address_space":               []string{"192.168.0.0/24"},
				"location":                    "westeurope",
				"resource_group_name":         "primary-rg",
				"resource_group_lock_enabled": false,
			},
		},
	}
}
