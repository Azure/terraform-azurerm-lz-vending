variable "location" {
  type        = string
  description = <<DESCRIPTION
The default location of resources created by this module.
Virtual networks will be created in this location unless overridden by the `location` attribute.
DESCRIPTION
  nullable    = false
}

variable "enable_telemetry" {
  type        = bool
  default     = true
  description = <<DESCRIPTION
This variable controls whether or not telemetry is enabled for the module.
For more information see <https://aka.ms/avm/telemetryinfo>.
If it is set to false, then no telemetry will be collected.
DESCRIPTION
  nullable    = false
}
