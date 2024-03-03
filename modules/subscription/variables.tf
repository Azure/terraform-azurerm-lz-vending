variable "subscription_alias_enabled" {
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

If disabled, supply the `subscription_id` variable to use an existing subscription instead.

> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.
DESCRIPTION

  default = false
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
  validation {
    condition     = length(var.subscription_alias_name) <= 64 && !can(regex("[<>;|]", var.subscription_alias_name))
    error_message = "Subscription Alias must either \"\", or be less or equal to 64 characters in length and cannot contain the characters `<`, `>`, `;`, or `|`"
  }
}

variable "subscription_display_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
The maximum length is 64 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
  validation {
    condition     = length(var.subscription_display_name) > 0 && length(var.subscription_display_name) <= 64 && !can(regex("[<>;|]", var.subscription_display_name))
    error_message = "Subscription Name must be between 1 and 64 characters in length and cannot contain the characters `<`, `>`, `;`, or `|`"
  }
}

variable "subscription_billing_scope" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

E.g.

- For CustomerLed and FieldLed, e.g. MCA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
- For PartnerLed, e.g. MPA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}`
- For Legacy EA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
  validation {
    condition     = can(regex("^$|^/providers/Microsoft.Billing/billingAccounts/.*$", var.subscription_billing_scope))
    error_message = "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive."
  }
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
  validation {
    condition     = can(regex("^$|^(Production|DevTest)$", var.subscription_workload))
    error_message = "The workload type can be either Production or DevTest and is case sensitive."
  }
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
  validation {
    condition     = can(regex("^$|^[().a-zA-Z0-9_-]{1,90}$", var.subscription_management_group_id))
    error_message = "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.)."
  }
}

variable "subscription_management_group_association_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create the `azurerm_management_group_subscription_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.
DESCRIPTION
}

variable "subscription_id" {
  type        = string
  description = <<DESCRIPTION
DESCRIPTION
  default     = ""
  validation {
    condition     = can(regex("^$|^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must be empty, or a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
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
  validation {
    error_message = "Tag values must be between 0-256 characters."
    condition = alltrue(
      [for _, v in var.subscription_tags : can(regex("^.{0,256}$", v))]
    )
  }
  validation {
    error_message = "Tag name must contain neither `<>%&\\?/` nor control characters, and must be between 0-512 characters."
    condition = alltrue(
      [for k, _ in var.subscription_tags : can(regex("^[^<>%&\\?/[:cntrl:]]{0,512}$", k))]
    )
  }
}

variable "subscription_use_azapi" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to use the azapi_resource resource to create the subscription alias. This includes the subscription alias in the management group.
DESCRIPTION
}

variable "subscription_update_existing" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to update an existing subscription with the supplied tags and display name.
If enabled, the following must also be supplied:
- `subscription_id`
DESCRIPTION
}

variable "wait_for_subscription_before_subscription_operations" {
  type = object({
    create  = optional(string, "30s")
    destroy = optional(string, "0s")
  })
  default     = {}
  description = <<DESCRIPTION
The duration to wait after vending a subscription before performing subscription operations.
DESCRIPTION
}

variable "subscription_dfc_contact_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to enable Microsoft Defender for Cloud (DFC) contact settings on the subscription. [optional - default `false`]
If enabled, provide settings in var.subscription_dfc_contact
DESCRIPTION
}

variable "subscription_dfc_contact" {
  type = object({
    notifications_by_role = optional(list(string), [])
    emails                = optional(string, "")
    phone                 = optional(string, "")
    alert_notifications   = optional(string, "Off")
  })
  default     = {}
  description = <<DESCRIPTION
Microsoft Defender for Cloud (DFC) contact and notification configurations

### Security Contact Information - Determines who'll get email notifications from Defender for Cloud 

- `notifications_by_role`: All users with these specific RBAC roles on the subscription will get email notifications. [optional - allowed values are: `AccountAdmin`, `ServiceAdmin`, `Owner` and `Contributor` - default empty]"
- `emails`: List of additional email addresses which will get notifications. Multiple emails can be provided in a ; separated list. Example: "john@microsoft.com;jane@microsoft.com". [optional - default empty]
- `phone`: The security contact's phone number. [optional - default empty]
> **Note**: At least one role or email address must be provided to enable alert notification.

### Alert Notifications

- `alert_notifications`: Enables email notifications and defines the minimal alert severity. [optional - allowed values are: `Off`, `High`, `Medium` or `Low` - default `Off`]

DESCRIPTION

  # validate email addresses
  validation {
    condition     = (var.subscription_dfc_contact.emails == "" || can(regex("^([\\w+-.%]+@[\\w.-]+\\.[A-Za-z]{2,4})(;[\\w+-.%]+@[\\w.-]+\\.[A-Za-z]{2,4})*$", var.subscription_dfc_contact.emails)))
    error_message = "Invalid email address(es) provided. Multiple emails must be separated with a `;`"
  }

  # validate phone number
  validation {
    condition     = (var.subscription_dfc_contact.phone == "" || can(regex("^[\\+0-9-]+$", var.subscription_dfc_contact.phone)))
    error_message = "Invalid phone number provided. Valid characters are 0-9, '-', and '+'. An example for a valid phone number is: +1-555-555-5555"
  }

  # validate alert notifications
  validation {
    condition     = contains(["Off", "High", "Medium", "Low"], var.subscription_dfc_contact.alert_notifications)
    error_message = "Invalid alert_notifications_state. Valid options are Off, High, Medium, Low."
  }

  # validate notifications by role
  validation {
    condition     = alltrue([for role in var.subscription_dfc_contact.notifications_by_role : contains(["Owner", "AccountAdmin", "Contributor", "ServiceAdmin"], role)])
    error_message = "Invalid notifications_by_role. The supported RBAC roles are: AccountAdmin, ServiceAdmin, Owner, Contributor."
  }

  # validate that when alert notifications are enabled, an email or role is also provided 
  validation {
    condition     = (var.subscription_dfc_contact.alert_notifications == "Off" ? true : var.subscription_dfc_contact.emails != "" || length(var.subscription_dfc_contact.notifications_by_role) > 0)
    error_message = "To enable alert notifications, either an email address or role must be provided."
  }

}
