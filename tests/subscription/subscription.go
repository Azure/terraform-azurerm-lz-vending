package subscription

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/azureutils"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"gopkg.in/matryer/try.v1"
)

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_alias_name":          name,
		"subscription_alias_display_name":  name,
		"subscription_alias_billing_scope": billingScope,
		"subscription_alias_workload":      "DevTest",
	}, nil
}

// cancelSubscription cancels the supplied Azure subscription.
// it retries a few times as the subscription api is eventually consistent.
func cancelSubscription(t *testing.T, id uuid.UUID) error {
	const (
		max      = 4
		delaysec = 20
	)

	if exists, err := subscriptionExists(id); err != nil || !exists {
		return fmt.Errorf("subscription %s does not exist or cannot successfully check, %s", id, err)
	}

	client, err := azureutils.NewSubscriptionClient()
	if err != nil {
		return fmt.Errorf("cannot create subscription client, %s", err)
	}
	ctx := context.Background()
	err = try.Do(func(attempt int) (bool, error) {
		_, err := client.Cancel(ctx, id.String(), nil)
		if err != nil {
			t.Logf("subscription id %s cancel failed, attempt %d/%d: %v", id, attempt, max, err)
			time.Sleep(delaysec * time.Second)
		}
		return attempt < max, err
	})
	if err != nil {
		return fmt.Errorf("cannot cancel subscription %s, %v", id, err)
	}
	return nil
}

// subscriptionExists checks if the supplied subscription exists.
func subscriptionExists(id uuid.UUID) (bool, error) {
	client, err := azureutils.NewSubscriptionsClient()
	if err != nil {
		return false, fmt.Errorf("cannot create subscriptions client, %s", err)
	}
	ctx := context.Background()
	if _, err := client.Get(ctx, id.String(), nil); err != nil {
		return false, fmt.Errorf("cannot get subscription, %s", err)
	}
	return true, nil
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_alias_name":          "test-subscription-alias",
		"subscription_alias_display_name":  "test-subscription-alias",
		"subscription_alias_billing_scope": "/providers/Microsoft.Billing/billingAccounts/test-billing-account",
		"subscription_alias_workload":      "Production",
	}
}

// isSubscriptionInManagementGroup returns true if the subscription is a management group.
func isSubscriptionInManagementGroup(t *testing.T, id uuid.UUID, mg string) error {
	const (
		max      = 8
		delaysec = 20
	)

	if exists, err := subscriptionExists(id); err != nil || !exists {
		return fmt.Errorf("subscription %s does not exist, or could not successfully check, %s", id, err)
	}

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

func setSubscriptionManagementGroup(id uuid.UUID, mg string) error {
	client, err := azureutils.NewManagementGroupSubscriptionsClient()
	if err != nil {
		return fmt.Errorf("cannot create mg subscriptions client, %s", err)
	}
	cc := "no-cache"
	opts := armmanagementgroups.ManagementGroupSubscriptionsClientCreateOptions{
		CacheControl: &cc,
	}
	if _, err := client.Create(context.Background(), mg, id.String(), &opts); err != nil {
		return fmt.Errorf("cannot create subscription %s in management group %s, %s", id.String(), mg, err)
	}
	return nil
}
