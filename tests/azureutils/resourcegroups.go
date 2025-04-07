package azureutils

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
)

// ListResourceGroup returns all resource groups in the subscription
func ListResourceGroup(ctx context.Context, subID uuid.UUID) ([]*armresources.ResourceGroup, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subID.String(), cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource group client: %v", err)
	}

	resultPager := resourceGroupClient.NewListPager(nil)

	resourceGroups := make([]*armresources.ResourceGroup, 0)
	for resultPager.More() {
		pageResp, err := resultPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		resourceGroups = append(resourceGroups, pageResp.Value...)
	}
	return resourceGroups, nil
}

// DeleteResourceGroup deletes a resource group by name and subscription id
func DeleteResourceGroup(ctx context.Context, rgname string, subID uuid.UUID) error {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %v", err)
	}
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subID.String(), cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create resource group client: %v", err)
	}

	pollerResp, err := resourceGroupClient.BeginDelete(ctx, rgname, nil)
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
