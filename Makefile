TESTTIMEOUT=30m
TESTFILTER=
TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'utils')

fmt:
	@echo "==> Fixing source code with gofmt..."
	find ./tests -name '*.go' | grep -v vendor | xargs gofmt -s -w

fumpt:
	@echo "==> Fixing source code with Gofumpt..."
	find . -name '*.go' | grep -v vendor | xargs gofumpt -w

fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

tools:
	go install mvdan.cc/gofumpt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.46.2

lint:
	cd tests && golangci-lint run

test: fmtcheck
	cd tests &&  go test $(TEST) $(TESTARGS) -timeout=$(TESTTIMEOUT) -run ^$(TESTFILTER)

testdeploy: fmtcheck
	cd tests &&	TERRATEST_DEPLOY=1 go test $(TEST) $(TESTARGS) -run ^TestDeploy$(TESTFILTER) -timeout $(TESTTIMEOUT)

# Makefile targets are files, but we aren't using it like this,
# so have to declare PHONY targets
.PHONY: test testdeploy lint tools fmt fumpt fmtcheck
