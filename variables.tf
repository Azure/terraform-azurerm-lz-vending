variable "location" {
  type        = string
  description = <<DESCRIPTION
The location of resources deployed by this module.
DESCRIPTION
  default     = ""
}

variable "disable_telemetry" {
  type        = bool
  description = <<DESCRIPTION
To disable tracking, we have included this variable with a simple boolean flag.
The default value is `false` which does not disable the telemetry.
If you would like to disable this tracking, then simply set this value to true and this module will not create the telemetry tracking resources and therefore telemetry tracking will be disabled.

For more information, see the [wiki](https://aka.ms/lz-vending/tf/telemetry)

E.g.

```terraform
module "lz-vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "..."

  # ... other module variables

  disable_telemetry = true
}
```
DESCRIPTION
  default     = false
}
