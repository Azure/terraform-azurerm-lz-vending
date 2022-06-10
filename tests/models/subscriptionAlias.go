package models

// SubscriptionAliasBody represents the JSON body of the subscription alias
type SubscriptionAliasBody struct {
	Properties *SubscriptionAliasBodyProperties `json:"properties"`
}

// SubscriptionAliasBodyProperties represents the JSON property bag of a subscription alias.
type SubscriptionAliasBodyProperties struct {
	AdditionalProperties SubscriptionAliasBodyAdditionalProperties `json:"additionalProperties,omitempty"`
	DisplayName          *string                                   `json:"displayName,omitempty"`
	SubscriptionId       *string                                   `json:"subscriptionId,omitempty"`
	Workload             *string                                   `json:"workload,omitempty"`
	BillingScope         *string                                   `json:"billingScope,omitempty"`
}

// SubscriptionAliasBodyAdditionalProperties represents the JSON additional properties bag of a subscription alias.
// Only currently used for the ManagementGroupId property.
type SubscriptionAliasBodyAdditionalProperties struct {
	ManagementGroupId *string `json:"managementGroupId,omitempty"`
}
