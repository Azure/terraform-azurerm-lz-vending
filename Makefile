TESTTIMEOUT=60m
TESTFILTER=
TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'utils')

docs:
	@echo "==> Updating documentation"
	terraform-docs -c .tfdocs-config.yml .
	find . | egrep ".md" | sort | while read f; do terrafmt fmt $f; done

fmt:
	@echo "==> Fixing source code with gofmt..."
	find ./tests -name '*.go' | grep -v vendor | xargs gofmt -s -w

fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

fumpt:
	@echo "==> Fixing source code with Gofumpt..."
	find ./tests -name '*.go' | grep -v vendor | xargs gofumpt -w

lint:
	cd tests && golangci-lint run

test: fmtcheck
	cd tests &&  go test $(TEST) $(TESTARGS) -timeout=$(TESTTIMEOUT) -run ^$(TESTFILTER)

testdeploy: fmtcheck
	cd tests &&	TERRATEST_DEPLOY=1 go test $(TEST) $(TESTARGS) -run ^TestDeploy$(TESTFILTER) -timeout $(TESTTIMEOUT)

tfclean:
	@echo "==> Cleaning terraform files..."
	find ./ -type d -name '.terraform' | xargs rm -vrf
	find ./ -type f -name 'tfplan' | xargs rm -vf
	find ./ -type f -name 'terraform.tfstate*' | xargs rm -vf
	find ./ -type f -name '.terraform.lock.hcl' | xargs rm -vf

tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/katbyte/terrafmt@latest
	go install github.com/terraform-docs/terraform-docs@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.46.2

# Makefile targets are files, but we aren't using it like this,
# so have to declare PHONY targets
.PHONY: docs fmt fmtcheck fumpt lint test testdeploy tfclean tools
