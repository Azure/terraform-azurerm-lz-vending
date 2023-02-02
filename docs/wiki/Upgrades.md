<!-- markdownlint-disable MD041 -->
This module uses [Semantic Versioning (SemVer)](https://semver.org/s) versioning.

Given a version number `MAJOR.MINOR.PATCH`, we increment the:

* `MAJOR` version when we make incompatible/breaking changes,
* `MINOR` version when we add functionality in a backwards compatible manner, and
* `PATCH` version when we make backwards compatible bug fixes.

## Upgrade process

If you are upgrading to a new `MINOR` or `PATCH` release, you will not see any breaking changes.
If you are using the Terraform registry, you can update the version number for the module.

For a new `MAJOR` release, you will see breaking changes.
We will publish guidance in the release notes on GitHub.

See the [release notes](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for more information.

## Upgrading from v1.x to v2.x

v2 of the module makes large-scale changes to the virtual networking capabilities of the module.
We therefore recommend that you keep any existing instances of the module at v1, and use v2 going forward for new instances.
If you would like multiple vnets in the same subscription using v1 of the module you can use the pattern [described here](https://github.com/Azure/terraform-azurerm-lz-vending/issues/97#issuecomment-1240712419)

## Upgrading from v2.x to v3.x

v3 of the module makes changes the the `role_assignments` variable, changing the format of the variable from a list of objects `list(object({...}))` to a map of objects `map(object({...}))`.
This change fixes [#153](https://github.com/Azure/terraform-azurerm-lz-vending/issues/153).

By way of explanation, in order to run the `for_each` loop on the role assignments, we need to use either a set or a map.
If using a map, Terraform needs to know all keys at plan time, they cannot be 'known after apply'.
Previously we converted the list of objects into a map of objects, using the `uuidv5()` function to generate predictable map keys from the inputs.
Unfortunately this caused issues when any of the inputs to the `uuidv5()` function were not known at plan time, in this case it was the principal id.

Rather than revert to a set, where ordering can be an issue, we decided to change the input variable to be a map from the outset.
This does mean a small change is required, you must specify a map key. This can be anything but do not use a reference to other object.

### v2.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = [
    {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  ]
}
```

### v3.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = {
    contrib_to_group = {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  }
}
```
