package resourceprovider

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/resourceprovider"
)

func TestSubscriptionRPRegistration(t *testing.T) {
	t.Parallel()

	v := make(map[string]any)
	v["resource_provider"] = "My.Rp"
	v["features"] = []any{"feature1", "feature2"}
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(3).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_registration").Exists().ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_feature_registration[\"feature2\"]").Exists().ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_feature_registration[\"feature1\"]").Exists().ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_registration").Key("action").HasValue("providers/My.Rp/register")
	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_registration").Key("resource_id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000")

	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_feature_registration[\"feature2\"]").Key("action").HasValue("register")
	check.InPlan(test.PlanStruct).That("azapi_resource_action.resource_provider_feature_registration[\"feature2\"]").Key("resource_id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Features/providers/My.Rp/features/feature2")
}
