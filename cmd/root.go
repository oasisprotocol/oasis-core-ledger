// Package cmd implements the oasis-core-ledger tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/oasisprotocol/oasis-core/go/common/logging"
	cmdCommon "github.com/oasisprotocol/oasis-core/go/oasis-node/cmd/common"

	"github.com/oasisprotocol/oasis-core-ledger/common"
)

const cfgLogLevel = "log.level"

var (
	rootCmd = &cobra.Command{
		Use:     "oasis-core-ledger",
		Short:   "Oasis Core Ledger Tool",
		Version: common.SoftwareVersion,
	}

	rootFlags = flag.NewFlagSet("", flag.ContinueOnError)
)

// RootCommand returns the root (top level) cobra.Command.
func RootCommand() *cobra.Command {
	return rootCmd
}

// Execute spawns the main entry point after handling the command line arguments.
func Execute() {
	var logLevel logging.Level
	if err := logLevel.Set(viper.GetString(cfgLogLevel)); err != nil {
		cmdCommon.EarlyLogAndExit(fmt.Errorf("root: failed to set log level: %w", err))
	}
	if err := logging.Initialize(os.Stdout, logging.FmtLogfmt, logLevel, nil); err != nil {
		cmdCommon.EarlyLogAndExit(fmt.Errorf("root: failed to initialize logging: %w", err))
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() { // nolint: gochecknoinits
	InitVersions(rootCmd)

	logLevel := logging.LevelInfo
	rootFlags.Var(&logLevel, cfgLogLevel, "log level")
	_ = viper.BindPFlags(rootFlags)
	rootCmd.PersistentFlags().AddFlagSet(rootFlags)

	// Register all of the sub-commands.
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showAddressCmd)
}
