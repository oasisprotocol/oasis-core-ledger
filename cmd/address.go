package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/oasisprotocol/oasis-core/go/common/logging"

	"github.com/oasisprotocol/oasis-core-ledger/common/wallet"
	"github.com/oasisprotocol/oasis-core-ledger/internal"
)

const (
	// cfgWalletID configures wallet ID.
	cfgWalletID = "wallet_id"

	// cfgIndex configures the wallet's account index (0-based).
	cfgIndex = "index"

	// cfgDontShowOnDevice configures whether staking account address should not
	// be shown on device's screen.
	cfgDontShowOnDevice = "no-show"
)

var (
	showAddressFlags = flag.NewFlagSet("", flag.ContinueOnError)

	showAddressCmd = &cobra.Command{
		Use:   "show_address",
		Short: "show staking account address",
		Run:   doShowAddress,
	}

	logger = logging.GetLogger("cmd")
)

func doShowAddress(cmd *cobra.Command, args []string) {
	var walletID wallet.ID
	if err := walletID.UnmarshalHex(viper.GetString(cfgWalletID)); err != nil {
		logger.Error("failed to parse wallet ID",
			"err", err,
		)
		os.Exit(1)
	}

	index := viper.GetUint32(cfgIndex)
	path := internal.GetPath(index)

	app, err := internal.ConnectApp(walletID, internal.ListingDerivationPath)
	if err != nil {
		logger.Error("failed to connect to ledger device",
			"wallet_id", walletID,
			"err", err,
		)
		os.Exit(1)
	}

	_, address, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		logger.Error("failed to get account address",
			"wallet_id", walletID,
			"index", index,
			"err", err,
		)
		os.Exit(1)
	}

	fmt.Println(address)

	if !viper.GetBool(cfgDontShowOnDevice) {
		fmt.Fprintln(os.Stderr, "Ensure account address shown on device's screen matches the outputted address.")
		_, _, err = app.ShowAddressPubKeyEd25519(path)
		if err != nil {
			logger.Error("failed to show account address",
				"wallet_id", walletID,
				"index", index,
				"err", err,
			)
			os.Exit(1)
		}
	}
}

func init() { //nolint:gochecknoinits
	showAddressFlags.String(cfgWalletID, "", "wallet ID")
	showAddressFlags.Uint32(cfgIndex, 0, "wallet's account index (0-based) (default 0)")
	showAddressFlags.Bool(cfgDontShowOnDevice, false, "don't show account address on device")
	_ = viper.BindPFlags(showAddressFlags)

	showAddressCmd.Flags().AddFlagSet(showAddressFlags)
}
