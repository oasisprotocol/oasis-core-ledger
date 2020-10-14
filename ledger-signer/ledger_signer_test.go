package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oasisprotocol/oasis-core-ledger/common/wallet"
)

func TestNewFactoryConfig(t *testing.T) {
	require := require.New(t)

	tWalletID := wallet.NewID([]byte("1640 Riverside Drive"))

	for _, t := range []struct { //nolint:maligned
		index    uint32
		walletID *wallet.ID
		cfgStr   string
		valid    bool
		errorMsg string
	}{
		// Valid configurations.

		// Wallet ID and account index are given.
		{17, &tWalletID, fmt.Sprintf("wallet_id:%s,index:17", tWalletID), true, ""},
		// Only account index is given.
		{17, nil, "index:17", true, ""},
		// Only wallet ID is given.
		{0, &tWalletID, fmt.Sprintf("wallet_id:%s", tWalletID), true, ""},
		// Empty configuration string is also valid.
		{0, nil, "", true, ""},

		// Invalid configurations.

		// Wallet ID is listed twice.
		{0, nil, fmt.Sprintf("wallet_id:%s,wallet_id:%s", tWalletID, tWalletID), false, "wallet ID already configured"},
		// Index is listed twice.
		{5, nil, fmt.Sprintf("wallet_id:%s,index:5,index:6", tWalletID), false, "index already configured"},
		// Empty configuration keys.
		{0, nil, ",", false, "malformed k/v pair: ''"},
	} {
		cfg, err := newPluginConfig(t.cfgStr)
		if !t.valid {
			require.EqualError(err, t.errorMsg, "newPluginConfig should fail to parse invalid config")
		} else {
			require.NoError(err, "newPluginConfig should parse valid config")
			require.Equal(t.index, cfg.index, "parsed index should be equal")
			require.Equal(t.walletID, cfg.walletID, "parsed wallet ID should be equal")
		}
	}
}
