package models

// VirtualNetworkBody represents the JSON body of the virtualNetwork
type VirtualNetworkBody struct {
	Properties *VirtualNetworkBodyProperties `json:"properties"`
}

// VirtualNetworkBodyProperties represents the JSON property bag of a virtualNetwork.
type VirtualNetworkBodyProperties struct {
	AddressSpace *VirtualNetworkBodyPropertiesAddressSpace `json:"addressSpace"`
	DhcpOptions  *VirtualNetworkBodyPropertiesDhcpOptions  `json:"dhcpOptions"`
}

// VirtualNetworkBodyPropertiesAddressSpace represents the JSON addressSpace of a virtualNetwork.
type VirtualNetworkBodyPropertiesAddressSpace struct {
	AddressPrefixes []string `json:"addressPrefixes"`
}

// VirtualNetworkBodyPropertiesDhcpOptions represents the JSON dhcpOptions of a virtualNetwork.
type VirtualNetworkBodyPropertiesDhcpOptions struct {
	DnsServers []string `json:"dnsServers"`
}
