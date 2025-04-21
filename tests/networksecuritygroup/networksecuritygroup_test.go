package networksecuritygroup

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/networksecuritygroup"
)

func TestNetworkSecurityGroup(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRulePrimary(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                     "Allow",
			"direction":                  "Inbound",
			"priority":                   100,
			"protocol":                   "Tcp",
			"source_port_range":          "*",
			"destination_port_range":     "*",
			"name":                       "test-rule",
			"source_address_prefix":      "*",
			"destination_address_prefix": "*",
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceAddressPrefix").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_address_prefix"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationAddressPrefix").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_address_prefix"]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRuleSourcePrefixes(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                     "Allow",
			"direction":                  "Inbound",
			"priority":                   100,
			"protocol":                   "Tcp",
			"destination_port_range":     "*",
			"destination_address_prefix": "*",
			"name":                       "test-rule",
			"source_port_ranges":         []string{"*"},
			"source_address_prefixes":    []string{"*"},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRanges.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_ranges"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceAddressPrefixes.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_address_prefixes"].([]string)[0]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRuleDestinationPrefixes(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                       "Allow",
			"direction":                    "Outbound",
			"priority":                     100,
			"protocol":                     "Tcp",
			"source_port_range":            "*",
			"source_address_prefix":        "*",
			"name":                         "test-rule",
			"destination_port_ranges":      []string{"*"},
			"destination_address_prefixes": []string{"*"},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRanges.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_ranges"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationAddressPrefixes.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_address_prefixes"].([]string)[0]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRulePrefixesOnly(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                       "Allow",
			"direction":                    "Inbound",
			"priority":                     100,
			"protocol":                     "Tcp",
			"name":                         "test-rule",
			"source_port_ranges":           []string{"*"},
			"destination_port_ranges":      []string{"*"},
			"source_address_prefixes":      []string{"*"},
			"destination_address_prefixes": []string{"*"},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRanges.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_ranges"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceAddressPrefixes.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_address_prefixes"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRanges.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_ranges"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationAddressPrefixes.0").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_address_prefixes"].([]string)[0]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRuleSourceAsgs(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                     "Allow",
			"direction":                  "Inbound",
			"priority":                   100,
			"protocol":                   "Tcp",
			"source_port_range":          "*",
			"destination_port_range":     "*",
			"destination_address_prefix": "*",
			"name":                       "test-rule",
			"source_application_security_group_ids": []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/sourceASG",
			},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationAddressPrefix").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_address_prefix"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceApplicationSecurityGroups.0.id").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_application_security_group_ids"].([]string)[0]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRuleDestinationAsgs(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                 "Allow",
			"direction":              "Inbound",
			"priority":               100,
			"protocol":               "Tcp",
			"source_port_range":      "*",
			"destination_port_range": "*",
			"source_address_prefix":  "*",
			"name":                   "test-rule",
			"destination_application_security_group_ids": []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/destinationASG",
			},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceAddressPrefix").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_address_prefix"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationApplicationSecurityGroups.0.id").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_application_security_group_ids"].([]string)[0]).ErrorIsNil(t)
}

func TestNetworkSecurityGroupSecurityRuleAsgsOnly(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["security_rules"] = map[string]map[string]any{
		"primary": {
			"access":                 "Allow",
			"direction":              "Inbound",
			"priority":               100,
			"protocol":               "Tcp",
			"source_port_range":      "*",
			"destination_port_range": "*",
			"name":                   "test-rule",
			"source_application_security_group_ids": []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/sourceASG",
			},
			"destination_application_security_group_ids": []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/destinationASG",
			},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("name").HasValue(v["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("location").HasValue(v["location"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.name").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["name"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.access").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["access"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.direction").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["direction"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.priority").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["priority"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.protocol").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["protocol"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourcePortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationPortRange").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_port_range"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.sourceApplicationSecurityGroups.0.id").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["source_application_security_group_ids"].([]string)[0]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.network_security_group").Key("body").Query("properties.securityRules.0.properties.destinationApplicationSecurityGroups.0.id").HasValue(v["security_rules"].(map[string]map[string]any)["primary"]["destination_application_security_group_ids"].([]string)[0]).ErrorIsNil(t)

}

func getMockInputVariables() map[string]any {
	return map[string]any{
		"name":                "test",
		"location":            "westeurope",
		"resource_group_name": "rg-test",
		"subscription_id":     "00000000-0000-0000-0000-000000000000",
		"enable_telemetry":    false,
	}
}
