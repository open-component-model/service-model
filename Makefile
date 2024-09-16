NAME      = serviceplugin
PROVIDER  ?= ocm.software
GITHUBORG  ?= open-component-model
COMPONENT = $(PROVIDER)/plugins/$(NAME)
OCMREPO   ?= ghcr.io/$(GITHUBORG)/ocm
PLATFORMS = linux/amd64 linux/arm64 darwin/arm64 darwin/amd64 windows/amd64
PLUGINTARGET = $${HOME}/.ocm/plugins

REPO_ROOT                                     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        = $(shell go run ./api/version/generate/release_generate.go print-rc-version $(CANDIDATE))
COMMIT                                         = $(shell git rev-parse HEAD)
EFFECTIVE_VERSION                              = $(VERSION)+$(COMMIT)
GIT_TREE_STATE                                := $(shell [ -z "$$(git status --porcelain 2>/dev/null)" ] && echo clean || echo dirty)

CMDSRCS=$(shell find $(REPO_ROOT)/plugins/$(NAME) -type f) Makefile
OCMSRCS=$(shell find $(REPO_ROOT)/api -type f) $(REPO_ROOT)/go.*
EXAMPLESRCS=$(shell find $(REPO_ROOT)/examples -type f) $(REPO_ROOT)/go.*

CREDS ?=
OCM = ocm $(CREDS)

GEN = $(REPO_ROOT)/gen

NOW         := $(shell date -u +%FT%T%z)
BUILD_FLAGS := "-s -w \
 -X github.com/open-component-model/service-model/api/version.gitVersion=$(EFFECTIVE_VERSION) \
 -X github.com/open-component-model/service-model/api/version.gitTreeState=$(GIT_TREE_STATE) \
 -X github.com/open-component-model/service-model/api/version.gitCommit=$(COMMIT) \
 -X github.com/open-component-model/service-model/api/version.buildDate=$(NOW)"


.PHONY: build
build: $(GEN)/.exists $(GEN)/$(NAME)/$(NAME)

gen/serviceplugin/serviceplugin:
	echo doit

$(GEN)/$(NAME)/$(NAME): $(CMDSRCS) $(OCMSRCS)
	CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME)/$(NAME) ./plugins/$(NAME)


.PHONY: test
test:
	go test ./...

.PHONY: install
install: $(GEN)/$(NAME)/$(NAME)
	cp $(GEN)/$(NAME)/$(NAME) "$(PLUGINTARGET)"

.PHONY: ctf
ctf:
	cd components/serviceplugin; make ctf

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: push
push:
	cd components/serviceplugin; make push

$(GEN)/.exists:
	@mkdir -p $(GEN)
	@touch $@

.PHONY: info
info:
	@echo "ROOT:     $(REPO_ROOT)"
	@echo "VERSION:  $(VERSION)"
	@echo "COMMIT;   $(COMMIT)"

.PHONY: clean
clean:
	rm -rf $(GEN)
