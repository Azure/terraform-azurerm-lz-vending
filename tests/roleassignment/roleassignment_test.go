package roleassignment

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/check"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/roleassignment"
)

// TestRoleAssignmentValidWithRoleName tests that the module will accept a role by name
func TestRoleAssignmentValidWithRoleName(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	test := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, test.Err)
	defer test.Cleanup()

	check.InPlan(test.Plan).NumberOfResourcesEquals(1).IfNotFail(t)

	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("role_definition_name").HasValue(v["role_assignment_definition"]).IfNotFail(t)
	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("role_definition_id").DoesNotExist().IfNotFail(t)
	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("principal_id").HasValue(v["role_assignment_principal_id"]).IfNotFail(t)
}

// TestRoleAssignmentValidWithRoleDefId tests that the module will accept a role by id
func TestRoleAssignmentValidWithRoleDefId(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["role_assignment_definition"] = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"
	test := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, test.Err)
	defer test.Cleanup()

	check.InPlan(test.Plan).NumberOfResourcesEquals(1).IfNotFail(t)

	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("role_definition_name").DoesNotExist().IfNotFail(t)
	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("role_definition_id").HasValue(v["role_assignment_definition"]).IfNotFail(t)
	check.InPlan(test.Plan).That("azurerm_role_assignment.this").Key("principal_id").HasValue(v["role_assignment_principal_id"]).IfNotFail(t)
}

// TestRoleAssignmentInvalidScopes tests that the module will not accept a tenant
// or management group scope for the role assignment.
func TestRoleAssignmentInvalidScopes(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["role_assignment_scope"] = "/"
	testTenant := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, testTenant.Err)
	defer testTenant.Cleanup()

	errString := "Must begin with a subscription scope, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`. All letters must be lowercase in the subscription id."
	assert.Contains(t, testTenant.Err, errString)

	// test management group scope error
	v["role_assignment_scope"] = "/providers/Microsoft.Management/managementGroups/myMg"
	testMg := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, testMg.Err)
	defer testMg.Cleanup()
	assert.Contains(t, testMg.Err, errString)
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"role_assignment_principal_id": "00000000-0000-0000-0000-000000000000",
		"role_assignment_scope":        "/subscriptions/00000000-0000-0000-0000-000000000000",
		"role_assignment_definition":   "Owner",
	}
}
