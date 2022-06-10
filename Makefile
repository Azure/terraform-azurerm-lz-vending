TESTTIMEOUT=180m
TESTFILTER=

test:
	cd tests && go test -v -run ^Test$(TESTFILTER) -timeout=$(TESTTIMEOUT)

testdeploy:
	cd tests &&	TERRATEST_DEPLOY=1 go test -v -run ^TestDeploy$(TESTFILTER) -timeout $(TESTTIMEOUT)

# Makefile targets are files, but we aren't using it like this,
# so have to declare PHONY targets
.PHONY: test testdeploy
