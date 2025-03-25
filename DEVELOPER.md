# Developer Requirements

* [Terraform (Core)](https://www.terraform.io/downloads.html) - version 1.x or above
* [Go](https://golang.org/doc/install) version 1.20.x (to run the tests)

## On Windows

If you're on Windows you'll also need:

* [Git Bash for Windows](https://git-scm.com/download/win)
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".

Or, use [Windows Subsystem for Linux](https://docs.microsoft.com/windows/wsl/install)

### Setup on WSL with Ubuntu 22.04

#### Install Go and Make

1. Run `curl -L https://go.dev/dl/go1.21.0.linux-amd64.tar.gz -o go1.21.0.linux-amd64.tar.gz` (replace with a different version of go if desired)
2. Run `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz`
3. Run `sudo nano ~/.profile`
4. Add the following lines at the end of the file:
```bash
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go/bin
export PATH=$PATH:$GOPATH/bin
```
5. Type <kbd>Ctrl</kbd> + <kbd>x</kbd> to save, then enter <kbd>y</kbd> and hit <kbd>enter</kbd>
5. Run `source ~/.profile`
6. Run `sudo apt-get update && apt-get install make`

#### Install Terraform

1. Follow these instructions `https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli`.

#### Setup Project Tools

1. Clone the repository
2. Navigate to the root of the repository
3. Run `make tools`

## Terratest

We use [Terratest](https://terratest.gruntwork.io/) to run the unit and deployment testing for the module. Therefore, if you wish to work on the module, you'll first need [Go](http://www.golang.org) installed on your machine.
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

### Unit Testing

#### Unit test environment variables

These tests do not deploy resources to an Azure environment, but may require access in order to run `terraform plan`. 

Example error when running virtual network tests without environment variables: `Error: subscription_id is a required provider property when performing a plan/apply operation`

Please make sure to set the following environment variables:

* `AZURE_TENANT_ID` - set to the tenant id of the Azure account.
* `AZURE_SUBSCRIPTION_ID` - set to the subscription id to use for deployment testing.

**NOTE:** You may login to your Azure account using `az login -t <tenant-id>`  and selecting the subscription from the cli. If you are not prompted you can run the `az account set --subscription <subscription-id>` command.

Use cases for unit testing:

* Validating variable inputs, e.g. validation rules are correct.
* Ensuring plan is generated successfully.
* Ensuring plan contents are correct.

To run the unit tests, run the following command:

```bash
make test
```

To run only a partial set of tests, add the TESTFILTER variable:

> The TESTFILTER is appended to the `-run ^Test` flag of `go test`.
> This will run the tests that match that regex.

```bash
make test TESTFILTER=Subscription
```

### Deployment Testing

These tests will deploy resources to an Azure environment, so ensure you are prepared to incur any costs.

Use cases for deployment testing:

* Validating the deployment is successful.
* Validating the deployment is idempotent.
* Destroying the deployment.

To run the unit tests, run the following command:

```bash
make testdeploy
```

To run only a partial set of tests, add the TESTFILTER variable:

> The TESTFILTER is appended to the `-run ^TestDeploy` flag of `go test`.
> This will run the tests that match that regex.

```bash
make testdeploy TESTFILTER=Subscription
```

#### Deployment environment variables

The following environment variables are required for deployment testing:

* `AZURE_BILLING_SCOPE` - set to the resource id of the billing scope to use for the deployment.
* `AZURE_SUBSCRIPTION_ID` - set to the subscription id to use for deployment testing.
* `AZURE_TENANT_ID` - set to the tenant id of the Azure account.
* `TERRATEST_DEPLOY` - set to a non-empty value to run the deployment tests. `make testdeploy` will do this for you.

## PR Naming

We have adopted [conventional commit](https://www.conventionalcommits.org/) naming standards for PRs.

E.g.:

```text
feat(roleassignment)!: add `relative_scope` value.
^    ^              ^ ^
|    |              | |__ Subject
|    |_____ Scope   |____ Breaking change flag
|__________ Type
```

### Type

The following types are permitted:

* `chore` - Other changes that do not modify src or test files
* `ci` - changes to the CI system
* `docs` - documentation only changes
* `feat` - a new feature (this correlates with `MINOR` in Semantic Versioning)
* `fix` - a bug fix (this correlates with `PATCH` in Semantic Versioning)
* `refactor` - a code change that neither fixes a bug or adds a feature
* `revert` - revert to a previous commit
* `style` - changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* `test` - adding or correcting tests

### Scope (Optional)

The following scopes are permitted:

* resourcegroup - pertaining to the resourcegroup sub-module
* roleassignment - pertaining to the roleassignment sub-module
* root - pertaining to the root module
* subscription - pertaining to the subscription sub-module
* usermanagedidentity - pertaining to the user-assigned managed identity sub-module
* virtualnetwork - pertaining to the virtual network sub-module

### Breaking Changes

An exclamation mark `!` is appended to the type/scope of a breaking change PR (this correlates with `MAJOR` in Semantic Versioning).
