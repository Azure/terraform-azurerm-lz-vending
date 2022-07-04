package azureutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/google/uuid"
	"github.com/matryer/try"
	"golang.org/x/sync/errgroup"
)

// CancelSubscription cancels the supplied Azure subscription.
// it retries a few times as the subscription api is eventually consistent.
func CancelSubscription(t *testing.T, id *uuid.UUID) error {
	const (
		max      = 4
		delaysec = 20
	)

	t.Logf("cancelling subscription %s", id.String())

	sub, err := GetSubscription(*id)
	if err != nil {
		return fmt.Errorf("subscription %s does not exist or cannot successfully check, %s", id, err)
	}

	client, err := NewSubscriptionClient()
	if err != nil {
		return fmt.Errorf("cannot create subscription client, %s", err)
	}
	ctx := context.TODO()
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	rgs, err := ListResourceGroup(ctx, *id)
	if err != nil {
		return fmt.Errorf("cannot list resource groups for subscription %s, %v", id, err)
	}

	t.Logf("removing %d resource groups for subscription %s", len(rgs), id)

	for _, rg := range rgs {
		rg := rg // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			t.Logf("removing resource group %s for subscription %s", *rg.Name, id.String())
			return DeleteResourceGroup(ctx, *rg.Name, *id)
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("cannot delete resource groups for subscription %s, %v", id, err)
	}
	t.Logf("removed %d resource groups for subscription %s", len(rgs), id)

	// If the sub is already in warned or disabled state then do not try and cancel again.
	if *sub.State == "Disabled" || *sub.State == "Warned" {
		t.Logf("subscription %s is already cancelled", id.String())
		return nil
	}

	ctx = context.TODO()
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
	t.Logf("cancelled subscription %s", id.String())
	return nil
}

// SubscriptionExists checks if the supplied subscription exists
func SubscriptionExists(id uuid.UUID) (bool, error) {
	client, err := NewSubscriptionsClient()
	if err != nil {
		return false, fmt.Errorf("cannot create subscriptions client, %s", err)
	}
	ctx := context.Background()
	if _, err := client.Get(ctx, id.String(), nil); err != nil {
		return false, fmt.Errorf("cannot get subscription, %s", err)
	}
	return true, nil
}

// GetSubscription checks if the supplied subscription exists and returns it
func GetSubscription(id uuid.UUID) (armsubscription.SubscriptionsClientGetResponse, error) {
	client, err := NewSubscriptionsClient()
	var resp armsubscription.SubscriptionsClientGetResponse
	if err != nil {
		return resp, fmt.Errorf("cannot create subscriptions client, %s", err)
	}
	ctx := context.Background()
	resp, err = client.Get(ctx, id.String(), nil)
	if err != nil {
		return resp, fmt.Errorf("cannot get subscription, %s", err)
	}
	return resp, nil
}

// IsSubscriptionInManagementGroup returns true if the subscription is a management group.
func IsSubscriptionInManagementGroup(t *testing.T, id uuid.UUID, mg string) error {
	// constants for retry loop in try.Do
	const (
		max      = 8
		delaysec = 20
	)

	if exists, err := SubscriptionExists(id); err != nil || !exists {
		return fmt.Errorf("subscription %s does not exist, or could not successfully check, %s", id, err)
	}

	client, err := NewManagementGroupSubscriptionsClient()
	if err != nil {
		return fmt.Errorf("cannot create mg subscriptions client, %s", err)
	}

	var mgopts armmanagementgroups.ManagementGroupSubscriptionsClientGetSubscriptionOptions
	cc := "no-cache"
	mgopts.CacheControl = &cc

	err = try.Do(func(attempt int) (bool, error) {
		_, err := client.GetSubscription(context.Background(), mg, id.String(), &mgopts)
		if err != nil {
			t.Logf("failed to get subscription %s in management group %s, attempt %d/%d", id.String(), mg, attempt, max)
			time.Sleep(delaysec * time.Second)
		}
		return attempt < max, err
	})
	if err != nil {
		return fmt.Errorf("failed determine if subscription %s in management group %s: %v", id.String(), mg, err)
	}
	return nil
}

// SetSubscriptionManagementGroup moves the subscription to the management group.
func SetSubscriptionManagementGroup(id uuid.UUID, mg string) error {
	client, err := NewManagementGroupSubscriptionsClient()
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
