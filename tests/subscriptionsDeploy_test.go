package alzLandingZoneTfModuleTest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestDeploySubscriptionAliasValid(t *testing.T) {
	preCheckDeployTests(t)

	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	if err != nil {
		t.Fatalf("Cannot generate valid input variables, %s", err)
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		NoColor:      true,
		Vars:         v,
		Logger:       getLogger(),
		PlanFilePath: "../tfplan",
	}

	_, err = terraform.InitAndApplyAndIdempotentE(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	sid := terraform.Output(t, terraformOptions, "subscription_id")
	if _, err := uuid.Parse(sid); err != nil {
		t.Errorf("subscription id output is not valid uuid: %s", sid)
		t.FailNow()
	}
	if err := cancelSubscription(sid); err != nil {
		t.Logf("could not cancel subscription: %v", err)
	}
	t.Logf("subscription %s cancelled", sid)
}

func preCheckDeployTests(t *testing.T) {
	variables := []string{
		"TERRATEST_DEPLOY",
		"AZURE_BILLING_SCOPE",
		"AZURE_TENANT_ID",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for deployment tests!", variable)
		}
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func cancelSubscription(id string) error {
	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzurePublic,
		},
		TenantID: os.Getenv("AZURE_TENANT_ID"),
	})
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := &arm.ClientOptions{
		DisableRPRegistration: true,
	}

	client, err := armsubscription.NewClient(cred, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create subscription client: %v", err)
	}
	ctx := context.Background()
	if _, err = client.Cancel(ctx, id, nil); err != nil {
		return fmt.Errorf("failed to cancel subscription: %v", err)
	}

	return nil
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]interface{}, error) {
	r, err := randomHex(4)
	if err != nil {
		fmt.Errorf("Cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_alias_name":          name,
		"subscription_alias_display_name":  name,
		"subscription_alias_billing_scope": billingScope,
		"subscription_alias_workload":      "DevTest",
	}, nil
}
