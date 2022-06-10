package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	if err != nil {
		t.Fatalf("Cannot generate valid input variables, %s", err)
	}

	terraformOptions := utils.GetDefaultTerraformOptions(v)

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)
	require.NoError(t, err)

	sid := terraform.Output(t, terraformOptions, "subscription_id")
	u, err := uuid.Parse(sid)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// cancel the newly created sub
	if err := cancelSubscription(u); err != nil {
		t.Logf("could not cancel subscription: %v", err)
	} else {
		t.Logf("subscription %s cancelled", sid)
	}
}

// Creating an alias for an existing subscription is not currently supported.
// Need use case data to justify the effort in testing support.
//
// // TestDeploySubscriptionAliasExistingSubscription tests the creation
// // of a subscription alias for an existing subscription
// func TestDeploySubscriptionAliasExistingSubscription(t *testing.T) {
// 	utils.PreCheckDeployTests(t)

// 	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
// 	v, err := getValidInputVariables(billingScope)
// 	if err != nil {
// 		t.Fatalf("Cannot generate valid input variables, %s", err)
// 	}

// 	existingSub, err := uuid.Parse(os.Getenv("AZURE_EXISTING_SUBSCRIPTION_ID"))
// 	if err != nil {
// 		t.Fatalf("Cannot parse AZURE_EXISTING_SUBSCRIPTION_ID as uuid, %s", err)
// 	}
// 	v["subscription_id"] = existingSub.String()

// 	terraformOptions := utils.GetDefaultTerraformOptions(v)

// 	_, err = terraform.InitAndPlanE(t, terraformOptions)
// 	require.NoError(t, err)

// 	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
// 	defer terraform.Destroy(t, terraformOptions)
// 	require.NoError(t, err)

// 	sid := terraform.Output(t, terraformOptions, "subscription_id")
// 	_, err = uuid.Parse(sid)
// 	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)
// 	// DO NOT CANCEL THIS SUBSCRIPTION
// }

// cancelSubscription cancels the supplied Azure subscription.
func cancelSubscription(id uuid.UUID) error {
	// Select the Azure cloud from the AZURE_ENVIRONMENT env var
	var cloudConfig cloud.Configuration
	env := os.Getenv("AZURE_ENVIRONMENT")
	switch strings.ToLower(env) {
	case "public":
		cloudConfig = cloud.AzurePublic
	case "usgovernment":
		cloudConfig = cloud.AzureGovernment
	case "china":
		cloudConfig = cloud.AzureChina
	default:
		cloudConfig = cloud.AzurePublic
	}

	// Get default credentials, this will look for the well-known environment variables,
	// managed identity credentials, and az cli credentials
	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
		TenantID: os.Getenv("AZURE_TENANT_ID"),
	})
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := &arm.ClientOptions{
		DisableRPRegistration: true,
	}

	// Create the subscriptions API client and cancel the subscription
	client, err := armsubscription.NewClient(cred, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create subscription client: %v", err)
	}
	ctx := context.Background()
	if _, err = client.Cancel(ctx, id.String(), nil); err != nil {
		return fmt.Errorf("failed to cancel subscription: %v", err)
	}

	return nil
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]interface{}, error) {
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
