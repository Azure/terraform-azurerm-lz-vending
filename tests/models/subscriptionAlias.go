package models

// SubscriptionAliasBody represents the JSON body of the subscription alias
type SubscriptionAliasBody struct {
	Properties *SubscriptionAliasBodyProperties `json:"properties"`
}

// SubscriptionAliasBodyProperties represents the JSON property bag of a subscription alias.
type SubscriptionAliasBodyProperties struct {
	DisplayName    *string `json:"displayName,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
	Workload       *string `json:"workload,omitempty"`
	BillingScope   *string `json:"billingScope,omitempty"`
}
