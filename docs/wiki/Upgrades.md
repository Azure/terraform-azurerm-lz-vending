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

## Upgrading from v1.x to v2.x

The v2 release of the module is a major release, therefore has breaking changes.
See the [release notes](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for more information.

v2 of the module makes large-scale changes to the virtual networking capabilities of the module.
We therefore recommend that you keep any existing instances of the module at v1, and use v2 going forward for new instances.
If you would like multiple vnets in the same subscription using v1 of the module you can use the pattern [described here](https://github.com/Azure/terraform-azurerm-lz-vending/issues/97#issuecomment-1240712419)
