include common.mk

# Check if Go's linkers flags are set in common.mk and add them as extra flags.
ifneq ($(GOLDFLAGS),)
	GO_EXTRA_FLAGS += -ldflags $(GOLDFLAGS)
endif

# Set all target as the default target.
all: build build-plugin

# Build.
build:
	@$(ECHO) "$(MAGENTA)*** Building oasis-core-ledger...$(OFF)"
	@$(GO) build $(GOFLAGS) $(GO_EXTRA_FLAGS) .

# Build plugin.
build-plugin:
	@$(ECHO) "$(MAGENTA)*** Building ledger signer plugin ...$(OFF)"
	@$(GO) build $(GOFLAGS) $(GO_EXTRA_FLAGS) -o ./ledger-signer/ledger-signer ./ledger-signer

# Format code.
fmt:
	@$(ECHO) "$(CYAN)*** Running Go formatters...$(OFF)"
	@gofumpt -s -w .
	@gofumports -w -local github.com/oasisprotocol/oasis-core-ledger .

# Lint code, commits and documentation.
lint-targets := lint-go lint-docs lint-changelog lint-git lint-go-mod-tidy

lint-go:
	@$(ECHO) "$(CYAN)*** Running Go linters...$(OFF)"
	@env -u GOPATH golangci-lint run

lint-git:
	@$(ECHO) "$(CYAN)*** Runnint gitlint...$(OFF)"
	@$(CHECK_GITLINT)

lint-docs:
	@$(ECHO) "$(CYAN)*** Runnint markdownlint-cli...$(OFF)"
	@npx markdownlint-cli '**/*.md' --ignore .changelog/

lint-changelog:
	@$(CHECK_CHANGELOG_FRAGMENTS)

lint-go-mod-tidy:
	@$(ECHO) "$(CYAN)*** Checking go mod tidy...$(OFF)"
	@$(ENSURE_GIT_CLEAN)
	@$(CHECK_GO_MOD_TIDY)

lint: $(lint-targets)

# Test.
test-targets := test-unit

test-unit:
	@$(ECHO) "$(CYAN)*** Running unit tests...$(OFF)"
	@$(GO) test -v -race ./...

test: $(test-targets)

# Clean.
clean:
	@$(ECHO) "$(CYAN)*** Cleaning up ...$(OFF)"
	@$(GO) clean -x
	rm -f ./ledger-signer/ledger-signer
	rm -rf ./__pycache__/

# Fetch all the latest changes (including tags) from the canonical upstream git
# repository.
fetch-git:
	@$(ECHO) "Fetching the latest changes (including tags) from $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote..."
	@git fetch $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) --tags

# Private target for bumping project's version using the Punch tool.
# NOTE: It should not be invoked directly.
_version-bump: fetch-git
	@$(ENSURE_GIT_VERSION_LATEST_TAG_EQUALS_PUNCH_VERSION)
	@$(PUNCH_BUMP_VERSION)
	@git add $(_PUNCH_VERSION_FILE)

# Private target for assembling the Change Log.
# NOTE: It should not be invoked directly.
_changelog:
	@$(ECHO) "$(CYAN)*** Generating Change Log for version $(_PUNCH_VERSION)...$(OFF)"
	@$(BUILD_CHANGELOG)
	@$(ECHO) "Next, review the staged changes, commit them and make a pull request."
	@$(WARN_BREAKING_CHANGES)

# Assemble Change Log.
# NOTE: We need to call Make recursively since _version-bump target updates
# Punch's version and hence we need Make to re-evaluate the _PUNCH_VERSION
# variable.
changelog: _version-bump
	@$(MAKE) --no-print-directory _changelog

# Tag the next release.
release-tag: fetch-git
	@$(ECHO) "Checking if we can tag version $(_PUNCH_VERSION) as the next release..."
	@$(ENSURE_VALID_RELEASE_BRANCH_NAME)
	@$(ENSURE_RELEASE_TAG_DOES_NOT_EXIST)
	@$(ENSURE_NO_CHANGELOG_FRAGMENTS)
	@$(ENSURE_NEXT_RELEASE_IN_CHANGELOG)
	@$(ECHO) "All checks have passed. Proceeding with tagging the $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH)'s HEAD with tag '$(_RELEASE_TAG)'."
	@$(CONFIRM_ACTION)
	@$(ECHO) "If this appears to be stuck, you might need to touch your security key for GPG sign operation."
	@git tag --sign --message="Version $(_PUNCH_VERSION)" $(_RELEASE_TAG) $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE)/$(RELEASE_BRANCH)
	@git push $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) $(_RELEASE_TAG)
	@$(ECHO) "$(CYAN)*** Tag '$(_RELEASE_TAG)' has been successfully pushed to $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote.$(OFF)"

# Create and push a stable branch for the current release.
release-stable-branch: fetch-git
	@$(ECHO) "Checking if we can create a stable release branch for version $(_PUNCH_VERSION)...$(OFF)"
	@$(ENSURE_VALID_STABLE_BRANCH)
	@$(ENSURE_RELEASE_TAG_EXISTS)
	@$(ENSURE_STABLE_BRANCH_DOES_NOT_EXIST)
	@$(ECHO) "All checks have passed. Proceeding with creating the '$(_STABLE_BRANCH)' branch on $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote."
	@$(CONFIRM_ACTION)
	@git branch $(_STABLE_BRANCH) $(_RELEASE_TAG)
	@git push $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) $(_STABLE_BRANCH)
	@$(ECHO) "$(CYAN)*** Branch '$(_STABLE_BRANCH)' has been sucessfully pushed to $(OASIS_CORE_LEDGER_GIT_ORIGIN_REMOTE) remote.$(OFF)"

# List of targets that are not actual files.
.PHONY: \
	all build build-plugin \
	fmt \
	$(lint-targets) lint \
	$(test-targets) test \
	clean \
	fetch-git \
	_version-bump _changelog \
	changelog \
	release-tag release-stable-branch
