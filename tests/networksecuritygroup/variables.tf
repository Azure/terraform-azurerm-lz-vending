variable "name" {
  description = "The name of the network security group"
  type        = string
}

variable "location" {
  description = "The location of the network security group"
  type        = string
}

variable "parent_id" {
  description = "The parent resource ID (resource group)"
  type        = string
}

variable "tags" {
  description = "Tags to apply to the network security group"
  type        = map(string)
  default     = {}
}

variable "security_rules" {
  description = "Security rules for the network security group"
  type = map(object({
    access                                     = string
    direction                                  = string
    priority                                   = number
    protocol                                   = string
    name                                       = string
    source_port_range                          = optional(string)
    source_port_ranges                         = optional(list(string))
    destination_port_range                     = optional(string)
    destination_port_ranges                    = optional(list(string))
    source_address_prefix                      = optional(string)
    source_address_prefixes                    = optional(list(string))
    destination_address_prefix                 = optional(string)
    destination_address_prefixes               = optional(list(string))
    source_application_security_group_ids      = optional(list(string))
    destination_application_security_group_ids = optional(list(string))
  }))
  default = {}
}
