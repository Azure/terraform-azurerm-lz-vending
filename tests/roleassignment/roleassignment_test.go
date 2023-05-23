package roleassignment

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
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
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("role_definition_name").HasValue(v["role_assignment_definition"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("role_definition_id").DoesNotExist().ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("principal_id").HasValue(v["role_assignment_principal_id"]).ErrorIsNil(t)
}

// TestRoleAssignmentValidWithRoleDefId tests that the module will accept a role by id
func TestRoleAssignmentValidWithRoleDefId(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	v["role_assignment_definition"] = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("role_definition_name").DoesNotExist().ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("role_definition_id").HasValue(v["role_assignment_definition"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azurerm_role_assignment.this").Key("principal_id").HasValue(v["role_assignment_principal_id"]).ErrorIsNil(t)
}

// TestRoleAssignmentInvalidScopes tests that the module will not accept a tenant
// or management group scope for the role assignment.
func TestRoleAssignmentInvalidScopes(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	errString := "Must begin with a subscription scope, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`. All letters must be lowercase in the subscription id."

	t.Run("tenant", func(t *testing.T) {
		v := v
		v["role_assignment_scope"] = "/"
		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		assert.Contains(t, utils.SanitiseErrorMessage(err), errString)
	})

	t.Run("managementGroup", func(t *testing.T) {
		v := v
		v["role_assignment_scope"] = "/providers/Microsoft.Management/managementGroups/myMg"
		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		assert.Contains(t, utils.SanitiseErrorMessage(err), errString)
	})
}

func getMockInputVariables() map[string]any {
	return map[string]any{
		"role_assignment_principal_id": "00000000-0000-0000-0000-000000000000",
		"role_assignment_scope":        "/subscriptions/00000000-0000-0000-0000-000000000000",
		"role_assignment_definition":   "Owner",
	}
}
