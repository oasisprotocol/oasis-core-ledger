package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	ledger_go "github.com/zondax/ledger-go"
)

const (
	testUseHardware = "OASIS_LEDGER_USE_HARDWARE"
	headerSize      = 5
	pathSize        = 5 * 4
)

var (
	_ ledger_go.LedgerDevice = (*MockOasisLedger)(nil)

	// Yes it would be better to derive these since we know the mnemonic
	// but dealing with BIP-44 and BIP-32 looked annoying.
	//
	// For now, just hard code the expected responses in a table, with
	// the caveat that it is not possible to test signing, and a bunch
	// of keys have placeholder addresses due to laziness.
	//
	// Ledger Test Mnemonic: equip will roof matter pink blind book anxiety banner elbow sun young
	testDeviceKeys = []*mockKeys{
		{
			decX("97e72e6e83ec39eb98d7e9189513aba662a08a210b9974b0f7197458483c7161"),
			"oasis1jlnjum5rasu7hxxhayvf2yat5e32pz3ppwvhfv8hr969sjpuw9sgn54g9",
		},
		{
			decX("54e98ea8afcf1321eddd2c91ee71f7f9237c38bd8c3242057be5c7ce3f46abbd"),
			"test key address 1",
		},
		{
			decX("7d10a11e1a4ef5adea33eb9f3332c6d221c12d461299de32d10e6cfffcd776d8"),
			"test key address 2",
		},
		{
			decX("00f3a005092933e8c2956d7ece62cbd39718678e35bf2a7370c344e9e755bc18"),
			"test key address 3",
		},
		{
			decX("3c713b1b2623c3a1c997b7b80c9dce4c49bf32c36dabb5cea6ce2cb6e89eb600"),
			"test key address 4",
		},
		{
			decX("636586ccbca4c1a5035552faccbce3b6ca59e6181ce17a3d84bcf6d9c5d120d1"),
			"test key address 5",
		},
		{
			decX("887fca7f936cad2733c6c8100c2ca8c612a37b9c7645b4a4b58445e5ceb6e862"),
			"test key address 6",
		},
		{
			decX("e2c22521953488a0135a4348dfd7544ff8ecfa1744fda1bef2f935476b909115"),
			"test key address 7",
		},
		{
			decX("5fec8d7031821c0a7ebbc18bdcaad826e1cf83323e172ce0a4f36a8e04792696"),
			"test key address 8",
		},
		{
			decX("72fde11509927324be809cdc815b258678ea74b2aa1d5e5490a960acd86c7a7e"),
			"test key address 9",
		},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		{
			// WARNING: The relevant app test uses `account = 5`, for this
			// test key.
			decX("d71c79ffd5a6d438de89c833e00222a2e80ed94e9929350ef7c1c97d1d13295d"),
			"oasis16uw8nl745m2r3h5feqe7qq3z5t5qak2wny5n2rhhc8yh68gn99wwcq4ef",
		},
	}
)

type mockKeys struct {
	publicKey []byte
	address   string
}

type MockOasisLedger struct {
	isClosed bool
}

func (dev *MockOasisLedger) Exchange(command []byte) ([]byte, error) {
	if dev.isClosed {
		return nil, os.ErrClosed
	}

	cmdLen := len(command)
	if cmdLen < headerSize {
		return nil, fmt.Errorf("oasis/ledger/mock: truncated command: %d", cmdLen)
	}

	// command[0] = CLA (ignored for now)
	// command[1] = instrution
	// command[2] = parameter 1
	// command[3] = parameter 2
	// command[4] = payload length
	switch command[1] {
	case insGetVersion:
		return dev.onGetVersion(command)
	case insGetAddrEd25519:
		return dev.onGetAddrEd25519(command)
	case insSignEd25519:
		return nil, fmt.Errorf("oasis/ledger/mock: sign not implemented yet")
	default:
		return nil, fmt.Errorf("oasis/ledger/mock: invalid command: %d", command[1])
	}
}

func (dev *MockOasisLedger) onGetVersion(cmd []byte) ([]byte, error) {
	return []byte{0x00, 0x00, 0x0d, 0x00, 0x00}, nil
}

func (dev *MockOasisLedger) onGetAddrEd25519(cmd []byte) ([]byte, error) {
	pathLen := int(cmd[4])
	if len(cmd) != headerSize+pathLen {
		return nil, fmt.Errorf("oasis/ledger/mock: truncated GetAddrEd25519: %d", len(cmd))
	}
	if pathLen != pathSize {
		return nil, fmt.Errorf("oasis/ledger/mock: truncated bip44 path: %d", pathLen)
	}

	path, err := parseBip44Path(cmd[headerSize:])
	if err != nil {
		return nil, err
	}

	addressIndex := int(path[4])
	if addressIndex >= len(testDeviceKeys) || testDeviceKeys[addressIndex] == nil {
		return nil, fmt.Errorf("oasis/ledger/mock: no key for address_index: %d", addressIndex)
	}

	key := testDeviceKeys[addressIndex]
	resp := append([]byte{}, key.publicKey...)
	resp = append(resp, []byte(key.address)...)

	return resp, nil
}

func (dev *MockOasisLedger) Close() error {
	if dev.isClosed {
		return os.ErrClosed
	}
	dev.isClosed = true
	return nil
}

func parseBip44Path(rawPath []byte) ([]uint32, error) {
	pathLen := len(rawPath)
	if pathLen != pathSize {
		return nil, fmt.Errorf("oasis/ledger/mock: truncated BIP44 path: %d", pathLen)
	}

	var path [5]uint32
	for i := range path {
		// Just unharden so the cheesy lookup table full of keys works.
		path[i] = binary.LittleEndian.Uint32(rawPath[i*4:]) & (^uint32(0x80000000))
	}

	return path[:], nil
}

func testFindLedgerOasisApp() (*LedgerOasis, error) {
	if testUsingHardware() {
		return FindLedgerOasisApp()
	}

	return newLedgerOasis(&MockOasisLedger{}, LedgerAppMode(0)), nil
}

func testUsingHardware() bool {
	return os.Getenv(testUseHardware) == "1"
}

func decX(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}

func checkTestKey(t *testing.T, pubKey []byte, address string, path []uint32) {
	index := int(path[4])
	if index >= len(testDeviceKeys) || testDeviceKeys[index] == nil {
		t.Fatalf("no known key for address index: %d", index)
	}

	key := testDeviceKeys[index]

	t.Logf("Public key %d: %x\n", index, pubKey)

	require := require.New(t)
	require.Len(pubKey, 32, "Public key should have expected length")
	require.Equal(key.publicKey, pubKey, "Public key should match %v", path)
	if address != "" {
		t.Logf("Bech32 addr %d: %s\n", index, address)

		require.Equal(key.address, address, "Address should match %v", path)
	}
}
