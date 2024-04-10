package budget

import (
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/budget"
)

// TestRoleAssignmentValidWithRoleName tests that the module will accept a role by name
func TestBudgetScopeSubscription(t *testing.T) {
	t.Parallel()

	v := map[string]interface{}{
		"budget_name":       "budget",
		"budget_scope":      "/subscriptions/00000000-0000-0000-0000-000000000000",
		"budget_amount":     1000,
		"budget_time_grain": "Monthly",
		"budget_time_period": map[string]interface{}{
			"start_date": time.Now().Format("2006-01") + "-01T00:00:00Z",
			"end_date":   time.Now().AddDate(0, 1, 0).Format("2006-01-02T15:04:05Z"),
		},
		"budget_notifications": map[string]interface{}{
			"notification1": map[string]interface{}{
				"enabled":        true,
				"operator":       "GreaterThanOrEqualTo",
				"threshold":      50,
				"threshold_type": "Actual",
				"contact_emails": []string{"email1@example.com", "email2@example.com"},
			},
			"notification2": map[string]interface{}{
				"enabled":        true,
				"operator":       "GreaterThan",
				"threshold":      75,
				"threshold_type": "Actual",
				"contact_roles":  []string{"role1", "role2"},
			},
		},
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
}
