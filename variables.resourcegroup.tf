variable "network_watcher_resource_group_enabled" {
  type        = bool
  description = <<DESCRIPTION
Create `NetworkWatcherRG` in the subscription.

Although this resource group is created automatically by Azure,
it is not managed by Terraform and therefore can impede the subscription cancellation process.

Enabling this variable will create the resource group in the subscription and allow Terraform to manage it,
which includes destroying the resource (and all resources within it).
DESCRIPTION
  default     = false
}

variable "resource_group_creation_enabled" {
  type        = bool
  description = "Whether to create additional resource groups in the target subscription. Requires `var.resource_groups_to_create`."
  default     = false
}

variable "resource_groups" {
  type = map(object({
    name     = string
    location = string
    tags     = optional(map(string), {})
  }))
  description = <<DESCRIPTION
A map of the resource groups to create. THe value is an object with the following attributes:

- `name` - the name of the resource group
- `location` - the location of the resource group
- `tags` - (optional) a map of type string

> Do not include the `NetworkWatcherRG` resource group in this map if you have enabled `var.network_watcher_resource_group_enabled`.
DESCRIPTION
  nullable    = false
  default     = {}
}
