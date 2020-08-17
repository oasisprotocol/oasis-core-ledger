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

# Boolean indicating whether to assume the 'yes' answer when confirming actions.
ASSUME_YES ?= 0

# Helper that asks the user to confirm the action.
define CONFIRM_ACTION =
	if [[ $(ASSUME_YES) != 1 ]]; then \
		$(ECHO) -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]; \
	fi
endef

# Name of git remote pointing to the canonical upstream git repository, i.e.
# git@github.com:oasisprotocol/oasis-core-ledger.git.
OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE ?= origin

# Name of the branch where to tag the next release.
RELEASE_BRANCH ?= master

# Determine the project's version from git.
GIT_VERSION_LATEST_TAG := $(shell git describe --tags --match 'v*' --abbrev=0 2>/dev/null $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) || echo undefined)
GIT_VERSION_LATEST := $(subst v,,$(GIT_VERSION_LATEST_TAG))
GIT_VERSION_IS_TAG := $(shell git describe --tags --match 'v*' --exact-match &>/dev/null $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) && echo YES || echo NO)
ifeq ($(GIT_VERSION_IS_TAG),YES)
	VERSION := $(GIT_VERSION_LATEST)
else
    # The current commit is not exactly a tag, append commit and dirty info to
    # the version.
    VERSION := $(GIT_VERSION_LATEST)-git$(shell git describe --always --match '' --dirty=+dirty 2>/dev/null)
endif
export VERSION

_PUNCH_CONFIG_FILE := .punch_config.py
_PUNCH_VERSION_FILE := .punch_version.py
# Obtain project's version as tracked by the Punch tool.
# NOTE: The Punch tool doesn't have the ability fo print project's version to
# stdout yet.
# For more details, see: https://github.com/lgiordani/punch/issues/42.
_PUNCH_VERSION := $(shell python3 -c "exec(open('$(_PUNCH_VERSION_FILE)').read()); print('{}.{}.{}'.format(major, minor, patch))")

# Helper that bumps project's version with the Punch tool.
define PUNCH_BUMP_VERSION =
	if [[ "$(RELEASE_BRANCH)" == master ]]; then \
		if [[ -n "$(CHANGELOG_FRAGMENTS_BREAKING)" ]]; then \
			PART=major; \
		else \
			PART=minor; \
		fi; \
	elif [[ "$(RELEASE_BRANCH)" == stable/* ]]; then \
		if [[ -n "$(CHANGELOG_FRAGMENTS_BREAKING)" ]]; then \
	        $(ECHO) "$(RED)Error: There shouldn't be breaking changes in a release on a stable branch.$(OFF)"; \
			$(ECHO) "List of detected breaking changes:"; \
			for fragment in "$(CHANGELOG_FRAGMENTS_BREAKING)"; do \
				$(ECHO) "- $$fragment"; \
			done; \
			exit 1; \
		else \
			PART=patch; \
		fi; \
    else \
	    $(ECHO) "$(RED)Error: Unsupported release branch: '$(RELEASE_BRANCH)'.$(OFF)"; \
		exit 1; \
	fi; \
	punch --config-file $(_PUNCH_CONFIG_FILE) --version-file $(_PUNCH_VERSION_FILE) --part $$PART --quiet
endef

# Helper that ensures project's version according to the latest Git tag equals
# project's version as tracked by the Punch tool.
define ENSURE_GIT_VERSION_EQUALS_PUNCH_VERSION =
	if [[ "$(GIT_VERSION_LATEST)" != "$(_PUNCH_VERSION)" ]]; then \
		$(ECHO) "$(RED)Error: Project version according to the latest Git tag from \
		    $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) ($(GIT_VERSION_LATEST)) \
			doesn't equal project's version in $(_PUNCH_VERSION_FILE) ($(_PUNCH_VERSION)).$(OFF)"; \
		exit 1; \
	fi
endef

# Go binary to use for all Go commands.
OASIS_GO ?= go

# Go command prefix to use in all Go commands.
GO := env -u GOPATH $(OASIS_GO)

# NOTE: The -trimpath flag strips all host dependent filesystem paths from
# binaries which is required for deterministic builds.
GOFLAGS ?= -trimpath -v

# Project's version as the linker's string value definition.
export GOLDFLAGS_VERSION := -X github.com/oasisprotocol/oasis-core-ledger/common.SoftwareVersion=$(VERSION)

# Go's linker flags.
export GOLDFLAGS ?= "$(GOLDFLAGS_VERSION)"

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

# List of breaking Change Log fragments.
CHANGELOG_FRAGMENTS_BREAKING := $(wildcard .changelog/*breaking*.md)

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

# Helper that builds the Change Log.
define BUILD_CHANGELOG =
	if [[ $(ASSUME_YES) != 1 ]]; then \
		towncrier build --version $(_PUNCH_VERSION); \
	else \
		towncrier build --version $(_PUNCH_VERSION) --yes; \
	fi
endef

# Helper that prints a warning when breaking changes are indicated by Change Log
# fragments.
define WARN_BREAKING_CHANGES =
	if [[ -n "$(CHANGELOG_FRAGMENTS_BREAKING)" ]]; then \
		$(ECHO) "$(RED)Warning: This release contains breaking changes.$(OFF)"; \
		$(ECHO) "$(RED)         Make sure the version is bumped appropriately.$(OFF)"; \
	fi
endef

# Helper that ensures the origin's release branch's HEAD doesn't contain any
# Change Log fragments.
define ENSURE_NO_CHANGELOG_FRAGMENTS =
	if ! CHANGELOG_FILES=`git ls-tree -r --name-only $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) .changelog`; then \
		$(ECHO) "$(RED)Error: Could not obtain Change Log fragments for $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) branch.$(OFF)"; \
		exit 1; \
	fi; \
	if CHANGELOG_FRAGMENTS=`echo "$$CHANGELOG_FILES" | grep --invert-match --extended-regexp '(README.md|template.md.j2|.markdownlint.yml)'`; then \
		$(ECHO) "$(RED)Error: Found the following Change Log fragments on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) branch:"; \
		$(ECHO) "$${CHANGELOG_FRAGMENTS}$(OFF)"; \
		exit 1; \
	fi
endef

# Helper that ensures the origin's release branch's HEAD contains a Change Log
# section for the next release.
define ENSURE_NEXT_RELEASE_IN_CHANGELOG =
	if ! ( git show $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH):CHANGELOG.md | \
			grep --quiet '^## $(_PUNCH_VERSION) (.*)' ); then \
		$(ECHO) "$(RED)Error: Could not locate Change Log section for release $(_PUNCH_VERSION) on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH) branch.$(OFF)"; \
		exit 1; \
	fi
endef

# Tag of the next release.
_RELEASE_TAG := v$(_PUNCH_VERSION)

# Helper that ensures the new release's tag doesn't already exist on the origin
# remote.
define ENSURE_RELEASE_TAG_EXISTS =
	if ! git ls-remote --exit-code --tags $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) $(_RELEASE_TAG) 1>/dev/null; then \
		$(ECHO) "$(RED)Error: Tag '$(_RELEASE_TAG)' doesn't exist on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote.$(OFF)"; \
		exit 1; \
	fi
endef

# Helper that ensures the new release's tag doesn't already exist on the origin
# remote.
define ENSURE_RELEASE_TAG_DOES_NOT_EXIST =
	if git ls-remote --exit-code --tags $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) $(_RELEASE_TAG) 1>/dev/null; then \
		$(ECHO) "$(RED)Error: Tag '$(_RELEASE_TAG)' already exists on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote.$(OFF)"; \
		exit 1; \
	fi; \
	if git show-ref --quiet --tags $(_RELEASE_TAG); then \
		$(ECHO) "$(RED)Error: Tag '$(_RELEASE_TAG)' already exists locally.$(OFF)"; \
		exit 1; \
	fi
endef

# Name of the stable release branch (if the current version is appropriate).
_STABLE_BRANCH := $(shell python3 -c "exec(open('$(_PUNCH_VERSION_FILE)').read()); print('stable/{}.{}.x'.format(major, minor)) if patch == 0 else print('undefined')")

# Helper that ensures the stable branch name is valid.
define ENSURE_VALID_STABLE_BRANCH =
	if [[ "$(_STABLE_BRANCH)" == "undefined" ]]; then \
		$(ECHO) "$(RED)Error: Cannot create a stable release branch for version $(_PUNCH_VERSION).$(OFF)"; \
		exit 1; \
	fi
endef

# Helper that ensures the new stable branch doesn't already exist on the origin
# remote.
define ENSURE_STABLE_BRANCH_DOES_NOT_EXIST =
	if git ls-remote --exit-code --heads $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) $(_STABLE_BRANCH) 1>/dev/null; then \
		$(ECHO) "$(RED)Error: Branch '$(_STABLE_BRANCH)' already exists on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote.$(OFF)"; \
		exit 1; \
	fi; \
	if git show-ref --quiet --heads $(_STABLE_BRANCH); then \
		$(ECHO) "$(RED)Error: Branch '$(_STABLE_BRANCH)' already exists locally.$(OFF)"; \
		exit 1; \
	fi
endef

# Helper that ensures $(RELEASE_BRANCH) variable contains a valid release branch
# name.
define ENSURE_VALID_RELEASE_BRANCH_NAME =
	if [[ ! $(RELEASE_BRANCH) =~ ^(master|(stable/[0-9]+\.[0-9]+\.x$$)) ]]; then \
		$(ECHO) "$(RED)Error: Invalid release branch name: '$(RELEASE_BRANCH)'."; \
		exit 1; \
	fi
endef
