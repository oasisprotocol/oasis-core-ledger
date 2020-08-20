package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oasisprotocol/oasis-core-ledger/common/wallet"
)

func TestNewFactoryConfig(t *testing.T) {
	require := require.New(t)

	expectedConfig := &pluginConfig{
		walletID: wallet.NewID([]byte("1640 Riverside Drive")),
		index:    17,
	}

	cfgStr := fmt.Sprintf("wallet_id:%s,index:%d", expectedConfig.walletID, expectedConfig.index)
	cfg, err := newPluginConfig(cfgStr)
	require.NoError(err, "newPluginConfig")
	require.Equal(expectedConfig, cfg, "config should parse")
}
