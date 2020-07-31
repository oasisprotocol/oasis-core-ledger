include common.mk

# Check if Go's linkers flags are set in common.mk and add them as extra flags.
ifneq ($(GOLDFLAGS),)
	GO_EXTRA_FLAGS += -ldflags $(GOLDFLAGS)
endif

# Set all target as the default target.
all: build build-plugin

# Build.
build:
	@$(ECHO) "$(MAGENTA)*** Building Go code...$(OFF)"
	@$(GO) build $(GOFLAGS) .

# Build plugin.
build-plugin:
	@$(ECHO) "$(MAGENTA)*** Building ledger signer plugin code...$(OFF)"
	@$(GO) build $(GOFLAGS) $(GO_EXTRA_FLAGS) -o ./ledger-signer/ledger-signer ./ledger-signer

# Format code.
fmt:
	@$(ECHO) "$(CYAN)*** Running Go formatters...$(OFF)"
	@gofumpt -s -w .
	@gofumports -w -local github.com/oasisprotocol/oasis-core-ledger .

# Lint code, commits and documentation.
lint-targets := lint-go lint-docs lint-git lint-go-mod-tidy

lint-go:
	@$(ECHO) "$(CYAN)*** Running Go linters...$(OFF)"
	@env -u GOPATH golangci-lint run

lint-git:
	@$(ECHO) "$(CYAN)*** Runnint gitlint...$(OFF)"
	@$(CHECK_GITLINT)

lint-docs:
	@$(ECHO) "$(CYAN)*** Runnint markdownlint-cli...$(OFF)"
	@npx markdownlint-cli '**/*.md'

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
	@rm -f ./ledger-signer/ledger-signer

# List of targets that are not actual files.
.PHONY: \
	all build build-plugin \
	fmt \
	$(lint-targets) lint \
	$(test-targets) test \
	clean
