// +build ledgerhw

/*******************************************************************************
*   (c) 2019 ZondaX GmbH
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
********************************************************************************/

package ledger_oasis_go

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ledger Test Mnemonic: equip will roof matter pink blind book anxiety banner elbow sun young

func Test_FindLedger(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, app)
	defer app.Close()
}

func Test_UserGetVersion(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	version, err := app.GetVersion()
	require.Nil(t, err, "Detected error")
	fmt.Println(version)

	assert.Equal(t, uint8(0x0), version.AppMode, "TESTING MODE ENABLED!!")
	assert.Equal(t, uint8(0x0), version.Major, "Wrong Major version")
	assert.Equal(t, uint8(0xd), version.Minor, "Wrong Minor version")
}

func Test_UserGetPublicKey(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, err := app.GetPublicKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	assert.Equal(t, 32, len(pubKey),
		"Public key has wrong length: %x, expected length: %x\n", pubKey, 32)
	fmt.Printf("PUBLIC KEY: %x\n", pubKey)

	assert.Equal(t,
		"d71c79ffd5a6d438de89c833e00222a2e80ed94e9929350ef7c1c97d1d13295d",
		hex.EncodeToString(pubKey),
		"Unexpected pubkey")
}

func Test_GetAddressPubKeyEd25519_Zero(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 0}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("PUBLIC KEY : %x\n", pubKey)
	fmt.Printf("BECH32 ADDR: %s\n", addr)

	assert.Equal(t, 32, len(pubKey), "Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

	assert.Equal(t, "97e72e6e83ec39eb98d7e9189513aba662a08a210b9974b0f7197458483c7161", hex.EncodeToString(pubKey), "Unexpected pubkey")
	assert.Equal(t, "oasis1jlnjum5rasu7hxxhayvf2yat5e32pz3ppwvhfv8hr969sjpuw9sgn54g9", addr, "Unexpected addr")
}

func Test_GetAddressPubKeyEd25519(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("PUBLIC KEY : %x\n", pubKey)
	fmt.Printf("BECH32 ADDR: %s\n", addr)

	assert.Equal(t, 32, len(pubKey), "Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

	assert.Equal(t, "d71c79ffd5a6d438de89c833e00222a2e80ed94e9929350ef7c1c97d1d13295d", hex.EncodeToString(pubKey), "Unexpected pubkey")
	assert.Equal(t, "oasis16uw8nl745m2r3h5feqe7qq3z5t5qak2wny5n2rhhc8yh68gn99wwcq4ef", addr, "Unexpected addr")
}

func Test_ShowAddressPubKeyEd25519(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 5, 0, 21}

	pubKey, addr, err := app.ShowAddressPubKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("PUBLIC KEY : %x\n", pubKey)
	fmt.Printf("BECH32 ADDR: %s\n", addr)

	assert.Equal(t, 32, len(pubKey), "Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

	assert.Equal(t, "d71c79ffd5a6d438de89c833e00222a2e80ed94e9929350ef7c1c97d1d13295d", hex.EncodeToString(pubKey), "Unexpected pubkey")
	assert.Equal(t, "oasis16uw8nl745m2r3h5feqe7qq3z5t5qak2wny5n2rhhc8yh68gn99wwcq4ef", addr, "Unexpected addr")
}

func Test_UserPK_HDPaths(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 0}

	expected := []string{
		"97e72e6e83ec39eb98d7e9189513aba662a08a210b9974b0f7197458483c7161",
		"54e98ea8afcf1321eddd2c91ee71f7f9237c38bd8c3242057be5c7ce3f46abbd",
		"7d10a11e1a4ef5adea33eb9f3332c6d221c12d461299de32d10e6cfffcd776d8",
		"00f3a005092933e8c2956d7ece62cbd39718678e35bf2a7370c344e9e755bc18",
		"3c713b1b2623c3a1c997b7b80c9dce4c49bf32c36dabb5cea6ce2cb6e89eb600",
		"636586ccbca4c1a5035552faccbce3b6ca59e6181ce17a3d84bcf6d9c5d120d1",
		"887fca7f936cad2733c6c8100c2ca8c612a37b9c7645b4a4b58445e5ceb6e862",
		"e2c22521953488a0135a4348dfd7544ff8ecfa1744fda1bef2f935476b909115",
		"5fec8d7031821c0a7ebbc18bdcaad826e1cf83323e172ce0a4f36a8e04792696",
		"72fde11509927324be809cdc815b258678ea74b2aa1d5e5490a960acd86c7a7e",
	}

	for i := uint32(0); i < 10; i++ {
		path[4] = i

		pubKey, err := app.GetPublicKeyEd25519(path)
		if err != nil {
			t.Fatalf("Detected error, err: %s\n", err.Error())
		}

		assert.Equal(
			t,
			32,
			len(pubKey),
			"Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

		assert.Equal(
			t,
			expected[i],
			hex.EncodeToString(pubKey),
			"Public key 44'/474'/0'/0/%d does not match\n", i)
	}
}

func Test_Sign(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 5}

	message := getDummyTx()

	println(coinContext)
	println(hex.EncodeToString(message))

	signature, err := app.SignEd25519(path, []byte(coinContext), message)
	if err != nil {
		t.Fatalf("[Sign] Error: %s\n", err.Error())
	}

	// Verify Signature
	pubKey, err := app.GetPublicKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	if err != nil {
		t.Fatalf("[GetPK] Error: " + err.Error())
		return
	}

	message = append([]byte(coinContext), message...)
	hash := sha512.Sum512(message)

	verified := ed25519.Verify(pubKey, hash[:], signature)
	if !verified {
		t.Fatalf("[VerifySig] Error verifying signature")
		return
	}
}

func Test_Sign_Fails(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	path := []uint32{44, 474, 0, 0, 5}

	message := getDummyTx()
	garbage := []byte{65}
	message = append(garbage, message...)

	_, err = app.SignEd25519(path, []byte(coinContext), message)
	assert.Error(t, err)
	errMessage := err.Error()
	assert.Equal(t, errMessage, "Unexpected data type")

	message = getDummyTx()
	garbage = []byte{65}
	message = append(message, garbage...)

	_, err = app.SignEd25519(path, []byte(coinContext), message)
	assert.Error(t, err)
	errMessage = err.Error()
	assert.Equal(t, errMessage, "Unexpected CBOR EOF")
}
