TESTTIMEOUT=60m
TESTFILTER=
TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'utils')
TESTARGS='-v'

default:
	@echo "==> Type make <thing> to run tasks"
	@echo
	@echo "Thing is one of:"
	@echo "docs fmt fmtcheck fumpt lint test testdeploy tfclean tools"

docs:
	@echo "==> Updating documentation..."
	find . | egrep ".md" | grep -v README.md | sort | while read f; do terrafmt fmt $$f; done
	terraform-docs -c .tfdocs-config.yml .

fmt:
	@echo "==> Fixing source code with gofmt..."
	find ./tests -name '*.go' | grep -v vendor | xargs gofmt -s -w
	@echo "==> Fixing Terraform code with terraform fmt..."
	terraform fmt -recursive
	@echo "==> Fixing embedded Terraform with terrafmt..."
	find . | egrep ".md|.tf" | grep -v README.md | sort | while read f; do terrafmt fmt $$f; done

fmtcheck:
	@echo "==> Checking source code with gofmt..."
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"
	@echo "==> Checking source code with terraform fmt..."
	terraform fmt -check -recursive

fumpt:
	@echo "==> Fixing source code with Gofumpt..."
	find ./tests -name '*.go' | grep -v vendor | xargs gofumpt -w

lint:
	cd tests && golangci-lint run

tftest-unit:
	@echo "==> Running unit tests in root module..."
	@if [ -d "$(CURDIR)/tests/unit" ]; then \
		terraform init -test-directory "$(CURDIR)/tests/unit"; \
		terraform test -test-directory "$(CURDIR)/tests/unit"; \
	fi
	@echo "==> Running unit tests in submodules..."
	@for dir in $$(find ./modules -maxdepth 1 -type d); do \
		cd "$$dir"; \
		echo "==> Running unit tests in submodule $$dir..."; \
		if [ -d "tests/unit" ]; then \
			terraform init -test-directory "tests/unit"; \
			terraform test -test-directory "tests/unit"; \
		fi; \
		cd -; \
	done

test: fmtcheck
	cd tests && go test $(TEST) $(TESTARGS) -run ^Test$(TESTFILTER) -timeout=$(TESTTIMEOUT)

testdeploy: fmtcheck
	cd tests &&	TERRATEST_DEPLOY=1 go test $(TEST) $(TESTARGS) -run ^TestDeploy$(TESTFILTER) -timeout $(TESTTIMEOUT)

tfclean:
	@echo "==> Cleaning terraform files..."
	find . -type d -name '.terraform' | xargs rm -vrf
	find . -type f -name 'tfplan' | xargs rm -vf
	find . -type f -name 'terraform.tfstate*' | xargs rm -vf
	find . -type f -name '.terraform.lock.hcl' | xargs rm -vf

tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/katbyte/terrafmt@latest
	go install github.com/terraform-docs/terraform-docs@v0.19
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.52.2

# Makefile targets are files, but we aren't using it like this,
# so have to declare PHONY targets
.PHONY: docs fmt fmtcheck fumpt lint test testdeploy tfclean tools
