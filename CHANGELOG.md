# Change Log

All notables changes to this project are documented in this file.

The format is inspired by [Keep a Changelog].

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/

<!-- markdownlint-disable no-duplicate-heading -->

<!-- NOTE: towncrier will not alter content above the TOWNCRIER line below. -->

<!-- TOWNCRIER -->

## 1.0.0 (2020-08-20)

### Process

- Add Change Log and the Change Log fragments process for assembling it
  ([#41](https://github.com/oasisprotocol/oasis-core-ledger/issues/41))

  This follows the same Change Log fragments process as is used by [Oasis Core].

  For more details, see [Change Log fragments].

  [Oasis Core]: https://github.com/oasisprotocol/oasis-core
  [Change Log fragments]: .changelog/README.md

- Define project's versioning
  ([#42](https://github.com/oasisprotocol/oasis-core-ledger/issues/42))

  Adopt a [Semantic Versioning 2.0.0].

  For more details, see [Versioning].

  [Semantic Versioning 2.0.0]: https://semver.org/spec/v2.0.0.html
  [Versioning]: README.md#versioning

- Define project's release process
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

  For more details, see [Release Process]

  [Release Process]: docs/release-process.md

### Removals and Breaking Changes

- Rename project to Oasis Core Ledger and transfer to Oasis Protocol Foundation
  ([#14](https://github.com/oasisprotocol/oasis-core-ledger/issues/14))

  The new home is at <https://github.com/oasisprotocol/oasis-core-ledger>.

- Rename `ledger_oasis_go` package to `internal` and move it to its directory
  ([#14](https://github.com/oasisprotocol/oasis-core-ledger/issues/14))

- Change identification of devices to use wallet IDs instead of App Addresses
  ([#27](https://github.com/oasisprotocol/oasis-core-ledger/issues/27))

  The new wallet IDs are six-characters hex strings deterministically derived
  from the mnemonics the Ledger devices were initialized with.

### Features

- ledger-signer: Add the Ledger signer plugin
  ([#16](https://github.com/oasisprotocol/oasis-core-ledger/issues/16))

- cmd: Add `oasis-core-ledger` executable with the `list_devices` CLI command
  ([#16](https://github.com/oasisprotocol/oasis-core-ledger/issues/16))

- common/wallet: Initial implementation of the wallet ID
  ([#46](https://github.com/oasisprotocol/oasis-core-ledger/issues/46))

  Wallet ID is computed as a truncated hash of a public key for a specific BIP32
  path.

  This means that two wallet IDs will be the same if and only if both Ledger
  devices were initialized with the same mnemonic.

- cmd: Improve listing of available devices
  ([#46](https://github.com/oasisprotocol/oasis-core-ledger/issues/46))

### Documentation Improvements

- Add [Oasis Core Ledger documentation] to GitBook
  ([#47](https://github.com/oasisprotocol/oasis-core-ledger/issues/47))

  [Oasis Core Ledger documentation]: https://docs.oasis.dev/oasis-core-ledger/

- docs: Add [README] and [Usage docs]
  ([#47](https://github.com/oasisprotocol/oasis-core-ledger/issues/47))

  [README]: docs/README.md
  [Usage docs]: docs/README.md#usage

### Internal Changes

- Replace Circle CI with [*ci-tests* GitHub Actions workflow]
  ([#12](https://github.com/oasisprotocol/oasis-core-ledger/issues/12))

  <!-- markdownlint-disable line-length -->
  [*ci-tests* GitHub Actions workflow]:
    https://github.com/oasisprotocol/oasis-core-ledger/actions?query=workflow:ci-tests
  <!-- markdownlint-enable line-length -->
  <!--
  gitlint-ignore: body-max-line-length
  -->

- Add new `Makefile`
  ([#14](https://github.com/oasisprotocol/oasis-core-ledger/issues/14))

- github: Add [*ci-lint* GitHub Actions workflow]
  ([#14](https://github.com/oasisprotocol/oasis-core-ledger/issues/14))

  <!-- markdownlint-disable line-length -->
  [*ci-lint* GitHub Actions workflow]:
    https://github.com/oasisprotocol/oasis-core-ledger/actions?query=workflow:ci-lint
  <!-- markdownlint-enable line-length -->
  <!--
  gitlint-ignore: body-max-line-length
  -->

- Add configuration for new linters: gitlint, markdownlint, golangci-lint
  ([#14](https://github.com/oasisprotocol/oasis-core-ledger/issues/14))

- internal: Clean up app tests
  ([#24](https://github.com/oasisprotocol/oasis-core-ledger/issues/24))

  Match Oasis Core style and exercise some of the code without a Ledger device.

- Add linting for Change Log fragments
  ([#41](https://github.com/oasisprotocol/oasis-core-ledger/issues/41),
   [#43](https://github.com/oasisprotocol/oasis-core-ledger/issues/43))

  Add `lint-changelog` Make target and *Lint Change Log fragments* step to the
  *ci-lint* GitHub Actions workflow.

- Use [Punch] tool for tracking and bumping project's version
  ([#42](https://github.com/oasisprotocol/oasis-core-ledger/issues/42))

  [Punch]: https://github.com/lgiordani/punch

- Make: Add `changelog` target for assembling the Change Log
  ([#42](https://github.com/oasisprotocol/oasis-core-ledger/issues/42))

- Make: Add `fetch-git` target for fetching changes from the canonical git repo
  ([#42](https://github.com/oasisprotocol/oasis-core-ledger/issues/42))

- Make: Reorganize how project's version is determined from git
  ([#42](https://github.com/oasisprotocol/oasis-core-ledger/issues/42))

- Use a proper tag for the github.com/oasisprotocol/oasis-core/go Go module
  ([#44](https://github.com/oasisprotocol/oasis-core-ledger/issues/44))

- Make: Add `release-build` target for building and publishing a release
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

- github: Add [*release* GitHub Actions workflow]
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

  <!-- markdownlint-disable line-length -->
  [*release* GitHub Actions workflow]:
    https://github.com/oasisprotocol/oasis-core-ledger/actions?query=workflow:release
  <!-- markdownlint-enable line-length -->
  <!--
  gitlint-ignore: body-max-line-length
  -->

- Make: Adjust handling of Go's linker flags
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

  Export `GOLDFLAGS_VERSION` so it can be consumed by other tools.

- Use [GoReleaser] tool for building and publishing releases
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

  [GoReleaser]: https://goreleaser.com/

- Make: Add `release-tag` and `release-stable-branch` targets
  ([#45](https://github.com/oasisprotocol/oasis-core-ledger/issues/45))

  The `release-tag` target can be used to tag the next release and
  `release-stable-branch` for creating and pushing a stable branch for the
  current release.

- internal: Rename functions for listing and connecting to Oasis Ledger Apps
  ([#46](https://github.com/oasisprotocol/oasis-core-ledger/issues/46))

  Renames:

  - `ListOasisDevices()` -> `ListApps()`
  - `ConnectLedgerOasisApp()` -> `ConnectApp()`
  - `FindLedgerOasisApp()` -> `FindApp()`

- internal: Add `AppInfo` type
  ([#46](https://github.com/oasisprotocol/oasis-core-ledger/issues/46))

  Refactor `ListApps()` to return a list of `AppInfo` pointers and leave the
  presentation of application information to the callers.
