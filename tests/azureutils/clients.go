package azureutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/google/uuid"
)

// NewSubnetClient creates a new subnet client using
// armnetwork.NewSubnetsClient
func NewSubnetClient(id uuid.UUID) (*armnetwork.SubnetsClient, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := arm.ClientOptions{
		DisableRPRegistration: true,
	}

	client, err := armnetwork.NewSubnetsClient(id.String(), cred, &clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subnet client: %v", err)
	}
	return client, nil
}

// NewSubscriptionsClient creates a new subscriptions client using
// azidentity.NewDefaultAzureCredential.
func NewSubscriptionsClient() (*armsubscription.SubscriptionsClient, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := &arm.ClientOptions{
		DisableRPRegistration: true,
	}

	client, err := armsubscription.NewSubscriptionsClient(cred, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions client: %v", err)
	}
	return client, nil
}

// NewSubscriptionClient creates a new subscription client using
// azidentity.NewDefaultAzureCredential.
func NewSubscriptionClient() (*armsubscription.Client, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := &arm.ClientOptions{
		DisableRPRegistration: true,
	}

	client, err := armsubscription.NewClient(cred, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription client: %v", err)
	}
	return client, nil
}

// NewManagementGroupSubscriptionsClient creates a new management group subscriptions client using
// azidentity.NewDefaultAzureCredential.
func NewManagementGroupSubscriptionsClient() (*armmanagementgroups.ManagementGroupSubscriptionsClient, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	clientOpts := &arm.ClientOptions{
		DisableRPRegistration: true,
	}

	client, err := armmanagementgroups.NewManagementGroupSubscriptionsClient(cred, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create management group subscription client: %v", err)
	}
	return client, nil
}

// newDefaultAzureCredential creates a new default AzureCredential using
// OIDC or azidentity.NewDefaultAzureCredential.
// OIDC is used if the environment variable USE_OIDC or ARM_USE_OIDC is set to non-empty.
func newDefaultAzureCredential() (azcore.TokenCredential, error) {
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

	useoidc := multiEnvDefault("", "USE_OIDC", "ARM_USE_OIDC")
	if useoidc != "" {
		return NewOidcCredential(&OidcCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: cloudConfig,
			},
			TenantID:      multiEnvDefault("", "ARM_TENANT_ID", "AZURE_TENANT_ID"),
			ClientID:      multiEnvDefault("", "ARM_CLIENT_ID", "AZURE_CLIENT_ID"),
			RequestToken:  multiEnvDefault("", "ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"),
			RequestUrl:    multiEnvDefault("", "ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"),
			Token:         multiEnvDefault("", "ARM_OIDC_TOKEN"),
			TokenFilePath: multiEnvDefault("", "ARM_OIDC_TOKEN_FILE_PATH"),
		})
	}

	// Get default credentials, this will look for the well-known environment variables,
	// managed identity credentials, and az cli credentials
	return azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
	})
}

func multiEnvDefault(dv string, envs ...string) string {
	for _, e := range envs {
		if v := os.Getenv(e); v != "" {
			return v
		}
	}
	return dv
}
