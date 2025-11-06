variable "resource_group_creation_enabled" {
  type        = bool
  default     = false
  description = "Whether to create additional resource groups in the target subscription. Requires `var.resource_groups`."
}

variable "resource_groups" {
  type = map(object({
    name         = string
    location     = optional(string)
    tags         = optional(map(string), {})
    lock_enabled = optional(bool, false)
    lock_name    = optional(string, "")
  }))
  default     = {}
  description = <<DESCRIPTION
A map of the resource groups to create. The value is an object with the following attributes:

- `name` - the name of the resource group
- `location` - the location of the resource group
- `tags` - (optional) a map of type string

We recommend that you include an entry to create the NetworkWatcherRG resource group so that this is managed by Terraform.
DESCRIPTION
  nullable    = false
}
