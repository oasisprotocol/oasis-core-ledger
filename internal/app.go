package internal

import (
	"encoding/hex"
	"fmt"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/common/logging"
	ledger_go "github.com/zondax/ledger-go"
)

const (
	// PathPurposeConsensus is set to 43, matching ledger's Validator app
	PathPurposeConsensus uint32 = 43

	errMsgInvalidParameters = "[APDU_CODE_BAD_KEY_HANDLE] The parameters in the data field are incorrect"
	errMsgInvalidated       = "[APDU_CODE_DATA_INVALID] Referenced data reversibly blocked (invalidated)"
	errMsgRejected          = "[APDU_CODE_COMMAND_NOT_ALLOWED] Sign request rejected"

	userMessageChunkSize = 250

	claConsumer  = 0x05
	claValidator = 0xF5

	insGetVersion     = 0
	insGetAddrEd25519 = 1
	insSignEd25519    = 2

	payloadChunkInit = 0
	payloadChunkAdd  = 1
	payloadChunkLast = 2
)

var (
	// ErrSignRequestRejected is the error returned when the user
	// explicitly rejects a signature request.
	ErrSignRequestRejected = fmt.Errorf("ledger/oasis: transaction rejected on Ledger device")

	logger = logging.GetLogger("oasis/ledger")

	minimumRequiredVersion = VersionInfo{0, 0, 3, 0}
)

type LedgerAppMode int

const (
	ValidatorMode LedgerAppMode = 1 + iota
	ConsumerMode
	UnknownMode
)

// LedgerOasis represents a connection to the Ledger app.
type LedgerOasis struct {
	device  ledger_go.LedgerDevice
	version VersionInfo
}

func newLedgerOasis(device ledger_go.LedgerDevice, mode LedgerAppMode) *LedgerOasis {
	return &LedgerOasis{
		device: device,
		version: VersionInfo{
			AppMode: uint8(mode),
		},
	}
}

func getModeForRole(role signature.SignerRole) LedgerAppMode { // nolint: deadcode,unused
	switch role {
	case signature.SignerConsensus:
		return ValidatorMode
	default:
		return ConsumerMode
	}
}

func getModeForPath(path []uint32) LedgerAppMode {
	switch path[0] {
	case PathPurposeConsensus:
		return ValidatorMode
	default:
		return ConsumerMode
	}
}

// Displays existing Ledger Oasis apps by address.
func ListOasisDevices(path []uint32) {
	ledgerAdmin := ledger_go.NewLedgerAdmin()

	mode := getModeForPath(path)

	for i := 0; i < ledgerAdmin.CountDevices(); i++ {
		ledgerDevice, err := ledgerAdmin.Connect(i)
		if err != nil {
			continue
		}
		defer ledgerDevice.Close()

		app := newLedgerOasis(ledgerDevice, mode)
		defer app.Close()

		appVersion, err := app.GetVersion()
		if err != nil {
			continue
		}

		_, address, err := app.GetAddressPubKeyEd25519(path)
		if err != nil {
			continue
		}

		fmt.Printf("============ Device found\n")
		fmt.Printf("Oasis App Version : %x\n", appVersion)
		fmt.Printf("Oasis App Address : %s\n", address)
	}
}

// ConnectLedgerOasisApp connects to Oasis app based on address.
func ConnectLedgerOasisApp(seekingAddress string, path []uint32) (*LedgerOasis, error) {
	ledgerAdmin := ledger_go.NewLedgerAdmin()

	mode := getModeForPath(path)

	for i := 0; i < ledgerAdmin.CountDevices(); i++ {
		ledgerDevice, err := ledgerAdmin.Connect(i)
		if err != nil {
			continue
		}

		app := newLedgerOasis(ledgerDevice, mode)

		_, address, err := app.GetAddressPubKeyEd25519(path)
		if err != nil {
			defer app.Close()
			continue
		}
		if seekingAddress == "" || address == seekingAddress {
			return app, nil
		}
	}
	return nil, fmt.Errorf("ledger/oasis: no app with specified address found")
}

// FindLedgerOasisApp finds the Oasis app running in a Ledger device.
func FindLedgerOasisApp() (*LedgerOasis, error) {
	ledgerAdmin := ledger_go.NewLedgerAdmin()

	for i := 0; i < ledgerAdmin.CountDevices(); i++ {
		ledgerDevice, err := ledgerAdmin.Connect(i)
		if err != nil {
			continue
		}

		app := newLedgerOasis(ledgerDevice, LedgerAppMode(0))

		appVersion, err := app.GetVersion()
		if err != nil {
			app.Close()
			continue
		}

		if err = app.CheckVersion(*appVersion); err != nil {
			app.Close()
			continue
		}

		return app, nil
	}

	return nil, fmt.Errorf("ledger/oasis: no app found")
}

// Close closes a connection with the Oasis user app.
func (ledger *LedgerOasis) Close() error {
	return ledger.device.Close()
}

// CheckVersion returns nil if the App version is supported by this library.
func (ledger *LedgerOasis) CheckVersion(ver VersionInfo) error {
	return checkVersion(ver, minimumRequiredVersion)
}

// GetVersion returns the current version of the Oasis user app.
func (ledger *LedgerOasis) GetVersion() (*VersionInfo, error) {
	message := []byte{ledger.getCLA(), insGetVersion, 0, 0, 0}
	response, err := ledger.device.Exchange(message)

	logger.Debug("GetVersion",
		"err", err,
		"message", hex.EncodeToString(message),
		"response", hex.EncodeToString(response),
	)

	if err != nil {
		return nil, fmt.Errorf("ledger/oasis: failed GetVersion request: %w", err)
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("ledger/oasis: truncated GetVersion response")
	}

	// WTF this tramples over the AppMode used to connect to the device.
	ledger.version = VersionInfo{
		AppMode: response[0],
		Major:   response[1],
		Minor:   response[2],
		Patch:   response[3],
	}

	return &ledger.version, nil
}

// SignEd25519 signs a transaction using Oasis user app
//
// NOTE: This command requires user confirmation on the device.
func (ledger *LedgerOasis) SignEd25519(bip44Path []uint32, context, transaction []byte) ([]byte, error) {
	return ledger.sign(bip44Path, context, transaction)
}

// GetPublicKeyEd25519 retrieves the public key for the corresponding BIP44
// derivation path.
//
// NOTE: This command DOES NOT require user confirmation on the device.
func (ledger *LedgerOasis) GetPublicKeyEd25519(bip44Path []uint32) ([]byte, error) {
	pubkey, _, err := ledger.retrieveAddressPubKeyEd25519(bip44Path, false)
	return pubkey, err
}

// GetAddressPubKeyEd25519 returns the pubkey and address (Bech32-encoded).
//
// NOTE: This command DOES NOT require user confirmation on the device.
func (ledger *LedgerOasis) GetAddressPubKeyEd25519(bip44Path []uint32) (pubkey []byte, addr string, err error) {
	return ledger.retrieveAddressPubKeyEd25519(bip44Path, false)
}

// ShowAddressPubKeyEd25519 returns the pubkey (compressed) and address (Bech32-encoded).
//
// NOTE: This command requires user confirmation on the device.
func (ledger *LedgerOasis) ShowAddressPubKeyEd25519(bip44Path []uint32) (pubkey []byte, addr string, err error) {
	return ledger.retrieveAddressPubKeyEd25519(bip44Path, true)
}

func (ledger *LedgerOasis) getCLA() byte {
	switch LedgerAppMode(ledger.version.AppMode) {
	case ValidatorMode:
		return claValidator
	default:
		return claConsumer
	}
}

func (ledger *LedgerOasis) sign(bip44Path []uint32, context, transaction []byte) ([]byte, error) {
	pathBytes, err := getBip44bytes(bip44Path, 5)
	if err != nil {
		return nil, fmt.Errorf("ledger/oasis: failed to get BIP44 bytes: %w", err)
	}

	chunks, err := prepareChunks(pathBytes, context, transaction, userMessageChunkSize)
	if err != nil {
		return nil, fmt.Errorf("ledger/oasis: failed to prepare chunks: %w", err)
	}

	var finalResponse []byte
	for idx, chunk := range chunks {
		payloadLen := byte(len(chunk))

		var payloadDesc byte
		switch idx {
		case 0:
			payloadDesc = payloadChunkInit
		case len(chunks) - 1:
			payloadDesc = payloadChunkLast
		default:
			payloadDesc = payloadChunkAdd
		}

		message := []byte{ledger.getCLA(), insSignEd25519, payloadDesc, 0, payloadLen}
		message = append(message, chunk...)

		response, err := ledger.device.Exchange(message)

		logger.Debug("Sign",
			"err", err,
			"message", hex.EncodeToString(message),
			"response", hex.EncodeToString(response),
		)

		if err != nil {
			switch err.Error() {
			case errMsgInvalidParameters, errMsgInvalidated:
				return nil, fmt.Errorf("ledger/oasis: failed to sign: %s", string(response))
			case errMsgRejected:
				return nil, ErrSignRequestRejected
			}
			return nil, fmt.Errorf("ledger/oasis: failed to sign: %w", err)
		}

		finalResponse = response
	}

	return finalResponse, nil
}

// retrieveAddressPubKeyEd25519 returns the pubkey and address (Bech32-encoded).
func (ledger *LedgerOasis) retrieveAddressPubKeyEd25519(
	bip44Path []uint32,
	requireConfirmation bool,
) (pubkey []byte, addr string, err error) {
	pathBytes, err := getBip44bytes(bip44Path, 5)
	if err != nil {
		return nil, "", fmt.Errorf("ledger/oasis: failed to get BIP44 bytes: %w", err)
	}

	p1 := byte(0)
	if requireConfirmation {
		p1 = byte(1)
	}

	// Prepare message
	header := []byte{ledger.getCLA(), insGetAddrEd25519, p1, 0, 0}
	message := append(header, pathBytes...)
	message[4] = byte(len(message) - len(header)) // update length

	response, err := ledger.device.Exchange(message)

	logger.Debug("GetAddrEd25519",
		"err", err,
		"message", hex.EncodeToString(message),
		"response", hex.EncodeToString(response),
	)

	if err != nil {
		return nil, "", fmt.Errorf("ledger/oasis: failed to request public key: %w", err)
	}
	if len(response) < 39 {
		return nil, "", fmt.Errorf("ledger/oasis: truncated GetAddrEd25519 response")
	}

	pubkey = response[0:32]
	addr = string(response[32:])

	return pubkey, addr, nil
}
