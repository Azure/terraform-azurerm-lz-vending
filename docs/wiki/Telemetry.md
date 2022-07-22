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

| Bit | Value (hex) | Value (denary) | Description |
| - | - | - | - |
| 0 (LSB) | 00000001 | 1 | `subscription_alias_enabled` is `true` |
| 1 | 00000002 | 2 | `subscription_management_group_association_enabled` is `true` |
| 2 | 00000004 | 4 | `subscription_tags` is not an empty object |
| 3 | 00000008 | 8 | reserved |
| 4 | 00000010 | 16 | reserved |
| 5 | 00000020 | 32 | reserved |
| 6 | 00000040 | 64 | reserved |
| 7 | 00000080 | 128 | reserved |
| 8 | 00000100 | 256 | `virtual_network_enabled` is `true` |
| 9 | 00000200 | 512 | `virtual_network_peering_enabled` is `true` |
| 10 | 00000400 | 1024 | `virtual_network_vwan_connection_enabled` is `true` |
| 11 | 00000800 | 2048 | `virtual_network_resource_lock_enabled` is `true` |
| 12 | 00001000 | 4096 |Either `virtual_network_vwan_propagated_routetables_labels` OR `virtual_network_vwan_propagated_routetables_resource_ids` are not empty lists OR `virtual_network_vwan_associated_routetable_resource_id` is not an empty string |
| 13 | 00002000 | 8192 | reserved |
| 14 | 00004000 | 16384 | reserved |
| 15 | 00008000 | 32768 | reserved |
| 16 | 00010000 | 65536 | `role_assignment_enabled` is `true` |
| 17 | 00020000 | 131072 | reserved|
| 18 | 00040000 | 262144 | reserved|
| 19 | 00080000 | 524288 | reserved |
| 20 | 00100000 | 1048576 | reserved |
| 21 | 00200000 | 2097152 | reserved |
| 22 | 00400000 | 4194304 | reserved |
| 23 | 00800000 | 8388608 | reserved |
| 24 | 01000000 | 16777216 | reserved |
| 25 | 02000000 | 33554432 | reserved |
| 26 | 04000000 | 67108864 | reserved |
| 27 | 08000000 | 134217728 | reserved |
| 28 | 10000000 | 268435456 | reserved |
| 29 | 20000000 | 536870912 | reserved |
| 30 | 40000000 | 1073741824 | reserved |
| 31 | 80000000 | 2147483648 | reserved |

[comment]: # (Link labels below, please sort a-z, thanks!)

[bit_field]: https://en.wikipedia.org/wiki/Bit_field
