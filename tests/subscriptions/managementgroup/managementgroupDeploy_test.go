package managementgroup

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/azureutils"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/subscriptions"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/matryer/try.v1"
)

// TestDeploySubscriptionAliasValidWithManagementGroup tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasValidWithManagementGroup(t *testing.T) {
	utils.PreCheckDeployTests(t)

	dir := utils.GetTestDir(t)
	terraformOptions := utils.GetDefaultTerraformOptions(t, dir)
	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := subscriptions.GetValidInputVariables(billingScope)
	require.NoError(t, err)
	v["subscription_id_billing_scope"] = billingScope
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 12)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)

	u, err := uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	err = isSubscriptionInManagementGroup(t, u, v["subscription_alias_management_group_id"].(string))
	assert.NoError(t, err)

	// cancel the newly created sub
	if err := subscriptions.CancelSubscription(t, u); err != nil {
		t.Logf("could not cancel subscription: %v", err)
	} else {
		t.Logf("subscription %s cancelled", sid)
	}
}

// isSubscriptionInManagementGroup returns true if the subscription is a management group.
func isSubscriptionInManagementGroup(t *testing.T, id uuid.UUID, mg string) error {
	const (
		max      = 12
		delaysec = 20
	)

	client, err := azureutils.NewManagementGroupSubscriptionsClient()
	if err != nil {
		return fmt.Errorf("cannot create mg subscriptions client, %s", err)
	}

	var mgopts armmanagementgroups.ManagementGroupSubscriptionsClientGetSubscriptionOptions
	cc := "no-cache"
	mgopts.CacheControl = &cc

	err = try.Do(func(attempt int) (bool, error) {
		_, err := client.GetSubscription(context.Background(), mg, id.String(), &mgopts)
		if err != nil {
			t.Logf("failed to get subscription %s in management group %s, attempt %d/%d: %v", id.String(), mg, attempt, max, err)
			time.Sleep(delaysec * time.Second)
		}
		return attempt < max, err
	})
	if err != nil {
		return fmt.Errorf("failed determine if subscription %s in management group %s: %v", id.String(), mg, err)
	}
	return nil
}
