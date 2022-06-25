package models

// HubVirtualNetworkConnectionBody represents the JSON body of the hub virtual network connection.
type HubVirtualNetworkConnectionBody struct {
	Properties *HubVirtualNetworkConnectionBodyProperties `json:"properties,omitempty"`
}

// HubVirtualNetworkConnectionBodyProperties represents the JSON property bag of a hub virtual network connection.
type HubVirtualNetworkConnectionBodyProperties struct {
	RemoteVirtualNetwork *HubVirtualNetworkConnectionBodyPropertiesRemoteVirtualNetwork `json:"remoteVirtualNetwork,omitempty"`
	RoutingConfiguration *HubVirtualNetworkConnectionBodyPropertiesRoutingConfiguration `json:"routingConfiguration,omitempty"`
}

// HubVirtualNetworkConnectionBodyPropertiesRemoteVirtualNetwork represents a reference to the remote virtual network.
type HubVirtualNetworkConnectionBodyPropertiesRemoteVirtualNetwork struct {
	ID string `json:"id,omitempty"`
}

// HubVirtualNetworkConnectionBodyPropertiesRoutingConfiguration represents the routing configuration of a hub virtual network connection.
type HubVirtualNetworkConnectionBodyPropertiesRoutingConfiguration struct {
	AssociatedRouteTable  *HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationAssociatedRouteTable  `json:"associatedRouteTable,omitempty"`
	PropagatedRouteTables *HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationPropagatedRouteTables `json:"propagatedRouteTables,omitempty"`
}

// HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationAssociatedRouteTable represents a reference to the associated route table for the hub virtual network connection.
type HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationAssociatedRouteTable struct {
	ID string `json:"id,omitempty"`
}

// HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationPropagatedRouteTables represents the routing propagation configuration of a hub virtual network connection.
type HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationPropagatedRouteTables struct {
	IDs    []HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationPropagatedRouteTablesIDs `json:"ids,omitempty"`
	Labels []string                                                                                `json:"labels,omitempty"`
}

type HubVirtualNetworkConnectionBodyPropertiesRoutingConfigurationPropagatedRouteTablesIDs struct {
	ID string `json:"id,omitempty"`
}
