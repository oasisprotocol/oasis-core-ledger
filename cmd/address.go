package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/common/logging"

	"github.com/oasisprotocol/oasis-core-ledger/common/wallet"
	"github.com/oasisprotocol/oasis-core-ledger/internal"
)

const (
	// cfgWalletID configures wallet ID.
	cfgWalletID = "wallet_id"

	// cfgIndex configures the wallet's account index (0-based).
	cfgIndex = "index"

	// cfgSkipDevice configures whether showing staking account address on
	// device's screen should be skipped or not.
	cfgSkipDevice = "skip-device"

	burnTxHex = "a463666565a26367617319099c66616d6f756e744064626f6479a166616d6f756e74417b656e6f6e636507666d6574686f646c7374616b696e672e4275726e"
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
	var walletID *wallet.ID
	hexWalletID := viper.GetString(cfgWalletID)
	if hexWalletID != "" {
		walletID = new(wallet.ID)
		if err := walletID.UnmarshalHex(hexWalletID); err != nil {
			logger.Error("failed to parse wallet ID",
				"err", err,
			)
			os.Exit(1)
		}
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

	pubKey, _, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		logger.Error("failed to get account address",
			"wallet_id", walletID,
			"index", index,
			"err", err,
		)
		os.Exit(1)
	}

	// Get the hard-coded burn transaction's CBOR.
	burnTxRaw, err := hex.DecodeString(burnTxHex)
	if err != nil {
		panic(err)
	}

	// Come up with an ad-hoc signature context and allow it.
	var signingContext signature.Context = "oasis-core/consensus: tx for chain 355e1d1e341e2bdf61df883fa58d907f2be3afc86ab3092a261a9af73b1a3569"
	signature.UnsafeAllowUnregisteredContexts()

	// Sign the transaction with the Ledger app.
	signingContextRaw := []byte(signingContext)
	sig, err := app.SignEd25519(path, signingContextRaw, burnTxRaw)
	if err != nil {
		logger.Error("failed to sign the first burn tx",
			"err", err,
		)
		os.Exit(1)
	}

	// Verify the signature with Go.
	data, err := signature.PrepareSignerMessage(signingContext, burnTxRaw)
	if err != nil {
		logger.Error("failed to prepare signer message for verification",
			"err", err,
		)
		os.Exit(1)
	}

	verified := ed25519.Verify(pubKey, data, sig)
	if verified {
		fmt.Println("Signature is OK")
	} else {
		fmt.Println("Signature is INVALID")
	}

	sig2, err := app.SignEd25519(path, signingContextRaw, burnTxRaw)
	if err != nil {
		logger.Error("failed to sign the second burn tx",
			"err", err,
		)
		os.Exit(1)
	}

	verified = ed25519.Verify(pubKey, data, sig2)
	if verified {
		fmt.Println("Signature is OK")
	} else {
		fmt.Println("Signature is INVALID")
	}

	// if !viper.GetBool(cfgSkipDevice) {
	// 	fmt.Fprintln(os.Stderr, "Ensure account address shown on device's screen matches the outputted address.")
	// 	_, _, err = app.ShowAddressPubKeyEd25519(path)
	// 	if err != nil {
	// 		logger.Error("failed to show account address",
	// 			"wallet_id", walletID,
	// 			"index", index,
	// 			"err", err,
	// 		)
	// 		os.Exit(1)
	// 	}
	// }
}

func init() { //nolint:gochecknoinits
	showAddressFlags.String(cfgWalletID, "", "wallet ID (can be omitted if only a single device is connected)")
	showAddressFlags.Uint32(cfgIndex, 0, "wallet's account index (0-based) (default 0)")
	showAddressFlags.Bool(cfgSkipDevice, false, "skip showing account address on device")
	_ = viper.BindPFlags(showAddressFlags)

	showAddressCmd.Flags().AddFlagSet(showAddressFlags)
}
