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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var coinContext = "oasis-core/consensus: tx for chain test-chain-id"

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

	app.api.Logging = true

	version, err := app.GetVersion()
	require.Nil(t, err, "Detected error")
	fmt.Println(version)

	assert.Equal(t, uint8(0x0), version.AppMode, "TESTING MODE ENABLED!!")
	assert.Equal(t, uint8(0x0), version.Major, "Wrong Major version")
	assert.Equal(t, uint8(0x5), version.Minor, "Wrong Minor version")
	assert.Equal(t, uint8(0x0), version.Patch, "Wrong Patch version")
}

func Test_UserGetPublicKey(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	app.api.Logging = true

	path := []uint32{44, 118, 5, 0, 21}

	pubKey, err := app.GetPublicKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	assert.Equal(t, 32, len(pubKey),
		"Public key has wrong length: %x, expected length: %x\n", pubKey, 32)
	fmt.Printf("PUBLIC KEY: %x\n", pubKey)

	assert.Equal(t,
		"3fbe6fd4cba729c23400510bbc49f90ec76d44204e7e921795e764e6a51fe387",
		hex.EncodeToString(pubKey),
		"Unexpected pubkey")
}

func Test_GetAddressPubKeyEd25519_Zero(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	app.api.Logging = true

	path := []uint32{44, 118, 0, 0, 0}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("PUBLIC KEY : %x\n", pubKey)
	fmt.Printf("BECH32 ADDR: %s\n", addr)

	assert.Equal(t, 32, len(pubKey), "Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

	assert.Equal(t, "2dc2ba9143c81c44830da47b0962ed026ba3e2e3d7a24d4a6ff1a418e9d154be", hex.EncodeToString(pubKey), "Unexpected pubkey")
	assert.Equal(t, "oasis19hpt4y2reqwyfqcd53asjchdqf468chr673y6jn07xjp36w32jlscf0me", addr, "Unexpected addr")
}

func Test_GetAddressPubKeyEd25519(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	app.api.Logging = true

	path := []uint32{44, 118, 5, 0, 21}

	pubKey, addr, err := app.GetAddressPubKeyEd25519(path)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("PUBLIC KEY : %x\n", pubKey)
	fmt.Printf("BECH32 ADDR: %s\n", addr)

	assert.Equal(t, 32, len(pubKey), "Public key has wrong length: %x, expected length: %x\n", pubKey, 32)

	assert.Equal(t, "3fbe6fd4cba729c23400510bbc49f90ec76d44204e7e921795e764e6a51fe387", hex.EncodeToString(pubKey), "Unexpected pubkey")
	assert.Equal(t, "oasis187lxl4xt5u5uydqq2y9mcj0epmrk63pqfelfy9u4uajwdfgluwrk0e5vx", addr, "Unexpected addr")
}

func Test_UserPK_HDPaths(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	app.api.Logging = true

	path := []uint32{44, 118, 0, 0, 0}

	expected := []string{
		"2dc2ba9143c81c44830da47b0962ed026ba3e2e3d7a24d4a6ff1a418e9d154be",
		"660e4d231f86a24549cae9e0e57879379fb33306374d0afbfee17075df569792",
		"2d5dfd35f26079c8b45cdb028ddd0d1d3c9a49a6e67935def053a0019528a357",
		"be01ac0e3b210c0488bb3e2ade554805f9fb92a2b6f2b3124a9d3055bd899432",
		"b31c272d30133ddd3750efaf631b7cf7f69724b9f9e345f1881017a288e19081",
		"658c6217e868b05bc221cd9fc7ec1720b156a3e7def0181d9bf772d4c988c625",
		"0efa4ee32aa2b0570fee44a1c9890fbf49c841f190d11be4bbe78a29ec21ee0d",
		"addc2e6cfe74b887d36a184a7c530846cc506101f7948afea1803df49ff844cb",
		"08d2c564382a9b2d9ed481bac0f23e98a24b6f685fb46f01b1388aa7735b22c7",
		"540c266ce814964a652c9dccd23377b8f9adc157b6b3c8269838fcd2223c875d",
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
			"Public key 44'/118'/0'/0/%d does not match\n", i)
	}
}

func getDummyTx() []byte {
	base64tx := "pGNmZWWiY2dhcwBmYW1vdW50QGRib2R5omd4ZmVyX3RvWCBkNhaFWEyIEubmS3EVtRLTanD3U+vDV5fke4Obyq" +
		"83CWt4ZmVyX3Rva2Vuc0Blbm9uY2UAZm1ldGhvZHBzdGFraW5nLlRyYW5zZmVy"
	tx, _ := base64.StdEncoding.DecodeString(base64tx)
	return tx;
}

func Test_Sign(t *testing.T) {
	app, err := FindLedgerOasisApp()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer app.Close()

	app.api.Logging = true

	path := []uint32{44, 118, 0, 0, 5}

	message := getDummyTx()
	signature, err := app.SignEd25519(path, message)
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

	app.api.Logging = true

	path := []uint32{44, 118, 0, 0, 5}

	message := getDummyTx()
	garbage := []byte{65}
	message = append(garbage, message...)

	_, err = app.SignEd25519(path, message)
	assert.Error(t, err)
	errMessage := err.Error()
	assert.Equal(t, errMessage, "Unexpected data type")

	message = getDummyTx()
	garbage = []byte{65}
	message = append(message, garbage...)

	_, err = app.SignEd25519(path, message)
	assert.Error(t, err)
	errMessage = err.Error()
	assert.Equal(t, errMessage, "Unexpected data type")

}
