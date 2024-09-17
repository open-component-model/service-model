NAME      = examples
PROVIDER  ?= acme.org
GITHUBORG = mandelsoft
OCMREPO   ?= ghcr.io/$(GITHUBORG)/ocm

REPO_ROOT                                     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../..

SRCS=$(shell find $(REPO_ROOT)/api -type f) $(REPO_ROOT)/examples/descriptors/*.yaml

CREDS ?=
OCM = ocm $(CREDS)

GEN = $(REPO_ROOT)/gen
CTF = $(GEN)/$(NAME).ctf

.PHONY: build
build: $(CTF)

$(CTF): $(GEN)/.exists $(SRCS) Makefile components.yaml
	$(OCM) add comp -c -t tgz  -F $(CTF) components.yaml


.PHONY: ctf
ctf: $(CTF)

$(GEN)/.exists:
	@mkdir -p $(GEN)
	@touch $@

.PHONY: info
info:
	@echo "ROOT:     $(REPO_ROOT)"
	@echo "CTF:      $(CTF)"

.PHONY: describe
describe: $(CTF)
	$(OCM) get resources -o treewide $(CTF)

.PHONY: descriptor
descriptor: $(CTF)
	$(OCM) get component -S v3alpha1 -o yaml $(CTF)

.PHONY: clean
clean:
	rm -f $(CTF)
