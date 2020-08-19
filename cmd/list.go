package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/oasisprotocol/oasis-core-ledger/internal"
)

var listCmd = &cobra.Command{
	Use:   "list_devices",
	Short: "list available devices",
	Run:   doList,
}

func doList(cmd *cobra.Command, args []string) {
	for _, appInfo := range internal.ListApps(internal.ListingDerivationPath) {
		fmt.Printf("- Wallet ID: %s\n", appInfo.WalletID)
		fmt.Printf("  App version: %s\n", appInfo.Version)
	}
}
