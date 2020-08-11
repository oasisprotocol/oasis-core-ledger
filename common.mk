SHELL := /bin/bash

# Check if we're running in an interactive terminal.
ISATTY := $(shell [ -t 0 ] && echo 1)

# If running interactively, use terminal colors.
ifdef ISATTY
	MAGENTA := \e[35;1m
	CYAN := \e[36;1m
	RED := \e[0;31m
	OFF := \e[0m
	# Use external echo command since the built-in echo doesn't support '-e'.
	ECHO_CMD := /bin/echo -e
else
	MAGENTA := ""
	CYAN := ""
	RED := ""
	OFF := ""
	ECHO_CMD := echo
endif

# Output messages to stderr instead stdout.
ECHO := $(ECHO_CMD) 1>&2

# Name of git remote pointing to the canonical upstream git repository, i.e.
# git@github.com:oasisprotocol/oasis-core-ledger.git.
OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE ?= origin

# Name of the branch where to tag the next release.
RELEASE_BRANCH ?= master

# Determine the project's version from git.
GIT_VERSION_LATEST_TAG := $(shell git describe --tags --match 'v*' --abbrev=0 2>/dev/null $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) || echo undefined)
GIT_VERSION_LATEST := $(subst v,,$(GIT_VERSION_LATEST_TAG))
GIT_VERSION_IS_TAG := $(shell git describe --tags --match 'v*' --exact-match 2>/dev/null $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) && echo YES || echo NO)
ifeq ($(and $(GIT_VERSION_LATEST_TAG),$(GIT_VERSION_IS_TAG)),YES)
	VERSION := $(GIT_VERSION_LATEST)
else
    # The current commit is not exactly a tag, append commit and dirty info to
    # the version.
    VERSION := $(GIT_VERSION_LATEST)-git$(shell git describe --always --match '' --dirty=+dirty 2>/dev/null)
endif
export VERSION

# Go binary to use for all Go commands.
OASIS_GO ?= go

# Go command prefix to use in all Go commands.
GO := env -u GOPATH $(OASIS_GO)

# NOTE: The -trimpath flag strips all host dependent filesystem paths from
# binaries which is required for deterministic builds.
GOFLAGS ?= -trimpath -v

# Add the plugin version as a linker string value definition.
ifneq ($(VERSION),)
	export GOLDFLAGS ?= "-X github.com/oasisprotocol/oasis-core-ledger/common.SoftwareVersion=$(VERSION)"
endif

# Helper that ensures the git workspace is clean.
define ENSURE_GIT_CLEAN =
	if [[ ! -z `git status --porcelain` ]]; then \
		$(ECHO) "$(RED)Error: Git workspace is dirty.$(OFF)"; \
		exit 1; \
	fi
endef

# Helper that checks if the go mod tidy command was run.
# NOTE: go mod tidy doesn't implement a check mode yet.
# For more details, see: https://github.com/golang/go/issues/27005.
define CHECK_GO_MOD_TIDY =
    $(GO) mod tidy; \
    if [[ ! -z `git status --porcelain go.mod go.sum` ]]; then \
        $(ECHO) "$(RED)Error: The following changes detected after running 'go mod tidy':$(OFF)"; \
        git diff go.mod go.sum; \
        exit 1; \
    fi
endef

# Helper that checks commits with gitlilnt.
# NOTE: gitlint internally uses git rev-list, where A..B is asymmetric
# difference, which is kind of the opposite of how git diff interprets
# A..B vs A...B.
define CHECK_GITLINT =
	BRANCH=$(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH); \
	COMMIT_SHA=`git rev-parse $$BRANCH` && \
	$(ECHO) "$(CYAN)*** Running gitlint for commits from $$BRANCH ($${COMMIT_SHA:0:7})... $(OFF)"; \
	gitlint --commits $$BRANCH..HEAD
endef

# List of non-trivial Change Log fragments.
CHANGELOG_FRAGMENTS_NON_TRIVIAL := $(filter-out $(wildcard .changelog/*trivial*.md),$(wildcard .changelog/[0-9]*.md))

# Helper that checks Change Log fragments with markdownlint-cli and gitlint.
# NOTE: Non-zero exit status is recorded but only set at the end so that all
# markdownlint or gitlint errors can be seen at once.
define CHECK_CHANGELOG_FRAGMENTS =
	exit_status=0; \
	$(ECHO) "$(CYAN)*** Running markdownlint-cli for Change Log fragments... $(OFF)"; \
	npx markdownlint-cli --config .changelog/.markdownlint.yml .changelog/ || exit_status=$$?; \
	$(ECHO) "$(CYAN)*** Running gitlint for Change Log fragments: $(OFF)"; \
	for fragment in $(CHANGELOG_FRAGMENTS_NON_TRIVIAL); do \
		$(ECHO) "- $$fragment"; \
		gitlint --msg-filename $$fragment -c title-max-length.line-length=78 || exit_status=$$?; \
	done; \
	exit $$exit_status
endef
