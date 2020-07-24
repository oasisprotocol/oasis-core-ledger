package cmd

import (
	"github.com/spf13/cobra"

	"github.com/oasisprotocol/oasis-core-ledger/internal"
)

var listCmd = &cobra.Command{
	Use:   "list_devices",
	Short: "list available devices by address",
	Run:   doList,
}

func doList(cmd *cobra.Command, args []string) {
	internal.ListOasisDevices(internal.ListingDerivationPath)
}
