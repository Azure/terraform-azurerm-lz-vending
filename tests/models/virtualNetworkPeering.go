package models

// VirtualNetworkPeeringBody is the body of the virtual network peering.
type VirtualNetworkPeeringBody struct {
	Properties *VirtualNetworkPeeringProperties `json:"properties,omitempty"`
}

// VirtualNetworkPeeringProperties represents the property bag of a virtual network peering.
type VirtualNetworkPeeringProperties struct {
	RemoteVirtualNetwork      *VirtualNetworkPeeringPropertiesRemoteVirtualNetwork `json:"remoteVirtualNetwork,omitempty"`
	AllowVirtualNetworkAccess *bool                                                `json:"allowVirtualNetworkAccess,omitempty"`
	AllowForwardedTraffic     *bool                                                `json:"allowForwardedTraffic,omitempty"`
	AllowGatewayTransit       *bool                                                `json:"allowGatewayTransit,omitempty"`
	UseRemoteGateways         *bool                                                `json:"useRemoteGateways,omitempty"`
}

// VirtualNetworkPeeringPropertiesRemoteVirtualNetwork represents a reference to a remote virtual network.
type VirtualNetworkPeeringPropertiesRemoteVirtualNetwork struct {
	Id string `json:"id,omitempty"`
}
