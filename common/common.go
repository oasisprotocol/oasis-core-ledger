// Package common implements common things for Oasis Core Ledger.
package common

import (
	"runtime"
	"strings"
)

var (
	// SoftwareVersion represents the Oasis Core Ledger's version and should be
	// set by the linker.
	SoftwareVersion = "0.0-unset"

	// ToolchainVersion is the version of the Go compiler/standard library.
	ToolchainVersion = strings.TrimPrefix(runtime.Version(), "go")

	// Versions contains versions of relevant dependencies.
	Versions = struct {
		Toolchain string
	}{
		ToolchainVersion,
	}
)
