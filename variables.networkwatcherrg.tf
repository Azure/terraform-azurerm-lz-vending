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
