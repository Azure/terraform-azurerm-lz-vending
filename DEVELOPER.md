# Developer Requirements

* [Terraform (Core)](https://www.terraform.io/downloads.html) - version 1.x or above
* [Go](https://golang.org/doc/install) version 1.18.x (to run the tests)

## On Windows

If you're on Windows you'll also need:

* [Git Bash for Windows](https://git-scm.com/download/win)
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".

Or, use [Windows Subsystem for Linux](https://docs.microsoft.com/windows/wsl/install)

## Terratest

We use [Terratest](https://terratest.gruntwork.io/) to run the unit and deployment testing for the module. Therefore, if you wish to work on the module, you'll first need [Go](http://www.golang.org) installed on your machine.
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

### Unit Testing

These tests do not deploy resources to an Azure environment, but may require access in order to run `terraform plan`.

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

These tests wil resources to an Azure environment, so ensure you are prepared to incur any costs.

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
* `TERRATEST_DEPLOY` - set to a non-empty value to run the deployemnt tests. `make testdeploy` will do this for you.
