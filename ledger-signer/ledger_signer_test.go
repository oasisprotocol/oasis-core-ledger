package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFactoryConfig(t *testing.T) {
	require := require.New(t)

	expectedConfig := &factoryConfig{
		address: "1640 Riverside Drive",
		index:   17,
	}

	cfgStr := fmt.Sprintf("address=%s;index=%d", expectedConfig.address, expectedConfig.index)
	cfg, err := newFactoryConfig(cfgStr)
	require.NoError(err, "newFactoryConfig")
	require.Equal(expectedConfig, cfg, "config should parse")
}
