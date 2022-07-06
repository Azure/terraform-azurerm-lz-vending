package azureutils

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/google/uuid"
)

func ListSubnets(rg, vnet string, subid uuid.UUID) ([]*armnetwork.Subnet, error) {
	ctx := context.Background()
	subnets := make([]*armnetwork.Subnet, 0)
	client, err := NewSubnetClient(subid)
	if err != nil {
		return nil, fmt.Errorf("failed to create subnet client: %v", err)
	}
	pager := client.NewListPager(rg, vnet, nil)
	for pager.More() {
		pageResp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list subnets: %v", err)
		}
		subnets = append(subnets, pageResp.SubnetListResult.Value...)
	}
	return subnets, nil
}
