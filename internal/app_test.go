package internal

import (
	"crypto/ed25519"
	"crypto/sha512"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindLedger(t *testing.T) {
	if !testUsingHardware() {
		t.Skipf("Hardware not configured for tests")
	}

	require := require.New(t)

	app, err := FindApp()
	require.NoError(err, "FindLedgerOasisApp")
	require.NotNil(app, "Must find a ledger device and initialize the interface")

	defer app.Close()
}

func TestUserGetVersion(t *testing.T) {
	require, assert := require.New(t), assert.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	version, err := app.GetVersion()
	require.NoError(err, "GetVersion")

	t.Logf("Version: %s", version)

	assert.Equal(uint8(0x0), version.AppMode, "TESTING MODE ENABLED!!")
	assert.Equal(uint8(0x0), version.Major, "Wrong Major version")
	assert.Equal(uint8(0xd), version.Minor, "Wrong Minor version")
}

func TestUserGetPublicKey(t *testing.T) {
	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, err := app.GetPublicKeyEd25519(path)
	require.NoError(err, "GetPublicKeyEd25519")

	checkTestKey(t, pubKey, "", path)
}

func TestGetAddressPubKeyEd25519_Zero(t *testing.T) {
	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 0}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	require.NoError(err, "GetAddressPublicKeyEd25519")

	checkTestKey(t, pubKey, addr, path)
}

func TestGetAddressPubKeyEd25519(t *testing.T) {
	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	require.NoError(err, "GetAddressPublicKeyEd25519")

	checkTestKey(t, pubKey, addr, path)
}

func TestShowAddressPubKeyEd25519(t *testing.T) {
	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, addr, err := app.ShowAddressPubKeyEd25519(path)
	require.NoError(err, "ShowAddressPublicKeyEd25519")

	checkTestKey(t, pubKey, addr, path)
}

func TestUserPKHDPaths(t *testing.T) {
	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 0}

	for i := uint32(0); i < 10; i++ {
		path[4] = i

		pubKey, err := app.GetPublicKeyEd25519(path)
		require.NoError(err, "GetPublicKeyEd25519")

		checkTestKey(t, pubKey, "", path)
	}
}

func TestSign(t *testing.T) {
	if !testUsingHardware() {
		// Can't sign with expected public key due to the incomplete
		// mock implementation.
		t.Skipf("Hardware not configured for tests")
	}

	require := require.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 5}

	message := getDummyTx()
	context := []byte(coinContext)

	t.Logf("Context: %x", context)
	t.Logf("Message: %x", message)

	signature, err := app.SignEd25519(path, context, message)
	require.NoError(err, "SignEd25519")

	// Verify Signature
	pubKey, err := app.GetPublicKeyEd25519(path)
	require.NoError(err, "GetPublicKeyEd25519")

	message = append(context, message...)
	hash := sha512.Sum512(message)

	verified := ed25519.Verify(pubKey, hash[:], signature)
	require.True(verified, "ed25519.Verify")
}

func TestSignFails(t *testing.T) {
	if !testUsingHardware() {
		// The mock implementation does not return the expected
		// error messages.
		t.Skipf("Hardware not configured for tests")
	}

	require, assert := require.New(t), assert.New(t)

	app, err := testFindLedgerOasisApp()
	require.NoError(err, "FindLedgerOasisApp")
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 5}

	message := getDummyTx()
	garbage := []byte{65}
	message = append(garbage, message...)

	_, err = app.SignEd25519(path, []byte(coinContext), message)
	assert.Error(err, "Signing unexpected data types should fail")
	assert.Equal(t, "Unexpected data type", err.Error())

	message = getDummyTx()
	garbage = []byte{65}
	message = append(message, garbage...)

	_, err = app.SignEd25519(path, []byte(coinContext), message)
	assert.Error(err, "Signing truncated CBOR payloads should fail")
	assert.Equal("Unexpected CBOR EOF", err.Error())
}
