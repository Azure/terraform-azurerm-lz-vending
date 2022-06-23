package subscriptions

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/azureutils"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"gopkg.in/matryer/try.v1"
)

// GetValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func GetValidInputVariables(billingScope string) (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("Cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_alias_name":          name,
		"subscription_alias_display_name":  name,
		"subscription_alias_billing_scope": billingScope,
		"subscription_alias_workload":      "DevTest",
	}, nil
}

// CancelSubscription cancels the supplied Azure subscription.
// it retries a few times as the subscription api is eventually consistent.
func CancelSubscription(t *testing.T, id uuid.UUID) error {
	const (
		max      = 12
		delaysec = 20
	)

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

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"subscription_alias_name":          "test-subscription-alias",
		"subscription_alias_display_name":  "test-subscription-alias",
		"subscription_alias_billing_scope": "/providers/Microsoft.Billing/billingAccounts/test-billing-account",
		"subscription_alias_workload":      "Production",
	}
}
