variable "subscription_id" {
  type        = string
  description = "The ID of the subscription to deploy resources into. E.g. `00000000-0000-0000-0000-000000000000`"
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a subscription id in the format of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "resource_groups_to_create" {
  type = map(object({
    name     = string
    location = string
    tags     = optional(map(string), {})
  }))
  default     = {}
  description = <<DESCRIPTION
A map of the resource groups to create. THe value is an object with the following attributes:

- `name` - the name of the resource group
- `location` - the location of the resource group
- `tags` - (optional) a map of type string
DESCRIPTION
}
