SHELL := /bin/bash
AVM_MAKEFILE_REF := main

$(shell curl -H 'Cache-Control: no-cache, no-store' -sSL "https://raw.githubusercontent.com/Azure/avm-terraform-governance/$(AVM_MAKEFILE_REF)/Makefile" -o avmmakefile)
-include avmmakefile
