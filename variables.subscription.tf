variable "subscription_alias_enabled" {
  default     = false
  type        = bool
  description = <<DESCRIPTION
Whether to create a new subscription using the subscription alias resource.

If enabled, the following must also be supplied:

- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`

Optionally, supply the following to enable the placement of the subscription into a management group:

- `subscription_management_group_id`
- `subscription_management_group_association_enabled`

If disabled, supply the `subscription_id` variable instead.
DESCRIPTION
}

variable "subscription_alias_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and _.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_display_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_billing_scope" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

E.g.

For CustomerLed and FieldLed, e.g. MCA - `/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
For PartnerLed, e.g. MPA - `/billingAccounts/{billingAccountName}/customers/{customerName}`
For Legacy EA - `/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_workload" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_management_group_id" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  The destination management group ID for the new subscription.

**Note:** Do not supply the display name.
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.
DESCRIPTION
}

variable "subscription_management_group_association_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create the `azurerm_management_group_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.
DESCRIPTION
}

variable "subscription_id" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
An existing subscription id.

Use this when you do not want the module to create a new subscription.
But do want to manage the management group membership.

A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
All letters must be lowercase.

When using this, `subscription_management_group_association_enabled` should be enabled,
and `subscription_management_group_id` should be supplied.

You may also supply an empty string if you want to create a new subscription alias.
In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:

- `subscription_alias_name`
- `subscription_alias_display_name`
- `subscription_alias_billing_scope`
- `subscription_alias_workload`
DESCRIPTION
}

variable "subscription_tags" {
  type        = map(string)
  description = <<DESCRIPTION
A map of tags to assign to the newly created subscription.
Only valid when `subsciption_alias_enabled` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```
DESCRIPTION
  default     = {}
}
