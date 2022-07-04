package azureutils

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
)

func ListResourceGroup(ctx context.Context, subId uuid.UUID) ([]*armresources.ResourceGroup, error) {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %v", err)
	}
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subId.String(), cred, nil)
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
		resourceGroups = append(resourceGroups, pageResp.ResourceGroupListResult.Value...)
	}
	return resourceGroups, nil
}

func DeleteResourceGroup(ctx context.Context, rgname string, subId uuid.UUID) error {
	cred, err := newDefaultAzureCredential()
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %v", err)
	}
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subId.String(), cred, nil)
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
