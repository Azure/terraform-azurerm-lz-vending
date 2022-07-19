<!-- markdownlint-disable MD041 -->
Microsoft can identify the deployments of this module with the deployed Azure resources.
Microsoft can correlate these resources used to support the deployments.
Microsoft collects this information to provide the best experiences with their products and to operate their business.
The telemetry is collected through customer usage attribution. The data is collected and governed by Microsoft's privacy policies, located at the trust center.

## Disabling telemetry

To disable this tracking, we have included a variable with the name disable_telemetry with a simple boolean flag.
The default value is false which does not disable the telemetry.
If you would like to disable this tracking, then simply set this value to true and this module will not create the telemetry tracking resources and therefore telemetry tracking will be disabled.

For example, to disable telemetry tracking, you can add this variable to the module declaration:

```terraform
module "lz-vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "..."

  # ... other module variables

  disable_telemetry = true
}
```

## Telemetry details

Telemetry is comprised of an empty ARM deployment per created subscription.
Each deployment contains a unique id (known as the PID) that is used to identity the particular module that is in use.
A [bit field][bit_field] is also used to identity module features that are in use.

### ARM deployment naming

The ARM deployment name is constructed as follows:

`pid-<UUID>_<module_version>_<bit_field>`

| Field | Description |
| - | - |
| `UUID` | A unique id to identify the Terraform module in use |
| `module_version` | The version of the module in use |
| `bitfield` | A bit field of 32 bits (eight hexadecimal digits) that exposes module features in use. See [next section](#bit-field-composition) for details |

The UUID for this module is `50a8a460-d517-4b11-b86c-6de447806b67`

### Bit field composition

The bit field is composed of the following bits:

| Bit | Value (hex) | Description |
| - | - | - |
| 1 (LSB) | 00000001 | `subscription_alias_enabled` is `true` |
| 2 | 00000002 | `subscription_management_group_association_enabled` is `true` |
| 3 | 00000004 | `subscripton_tags` is not an empty object |
| 4 | 00000008 | reserved |
| 5 | 00000010 | reserved |
| 6 | 00000020 | reserved |
| 7 | 00000040 | reserved |
| 8 | 00000080 | reserved |
| 9 | 00000100 | `virtual_network_enabled` is `true` |
| 10 | 00000200 | `virtual_network_peering_enabled` is `true` |
| 11 | 00000400 | `virtual_network_vwan_connection_enabled` is `true` |
| 12 | 00000800 | `virtual_network_resource_lock_enabled` is `true` |
| 13 | 00001000 | Either `virtual_network_vwan_propagated_routetables_labels` OR `virtual_network_vwan_propagated_routetables_resource_ids` are not an empty lists OR `virtual_network_vwan_routetable_resource_id` is not an empty string |
| 14 | 00002000 | reserved |
| 15 | 00004000 | reserved |
| 16 | 00008000 | reserved |
| 17 | 00010000 | `role_assignment_enabled` is `true` |
| 18 | 00020000 | reserved|
| 19 | 00040000 | reserved|
| 20 | 00080000 | reserved |
| 21 | 00100000 | reserved |
| 22 | 00200000 | reserved |
| 23 | 00400000 | reserved |
| 24 | 00800000 | reserved |
| 25 | 01000000 | reserved |
| 26 | 02000000 | reserved |
| 27 | 04000000 | reserved |
| 28 | 08000000 | reserved |
| 29 | 10000000 | reserved |
| 30 | 20000000 | reserved |
| 31 | 40000000 | reserved |
| 32 | 80000000 | reserved |

[comment]: # (Link labels below, please sort a-z, thanks!)

[bit_field]: https://en.wikipedia.org/wiki/Bit_field
