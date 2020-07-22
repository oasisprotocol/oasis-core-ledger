package ledger_oasis_go

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var coinContext = "oasis-core/consensus: tx for chain 7b02d647e8997bacebce96723f6904029ec78b67c261c4bdddb5e47de1ab31fa"

func getDummyTx() []byte {
	base64tx := "pGNmZWWiY2dhcxkD6GZhbW91bnRCB9BkYm9keaJneGZlcl90b1UA4ywoibwEEhHt7fqvlNL9hmmLsH9reGZlcl90b2tlbnNFJ5TKJABlbm9uY2UHZm1ldGhvZHBzdGFraW5nLlRyYW5zZmVy"
	tx, _ := base64.StdEncoding.DecodeString(base64tx)
	println(hex.EncodeToString(tx))
	return tx
}

func Test_PrintVersion(t *testing.T) {
	reqVersion := VersionInfo{0, 1, 2, 3}
	s := fmt.Sprintf("%v", reqVersion)
	assert.Equal(t, "1.2.3", s)
}

func Test_PathGeneration0(t *testing.T) {
	bip44Path := []uint32{44, 100, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 0)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("Path: %x\n", pathBytes)

	assert.Equal(
		t,
		20,
		len(pathBytes),
		"PathBytes has wrong length: %x, expected length: %x\n", pathBytes, 40)

	assert.Equal(
		t,
		"2c00000064000000000000000000000000000000",
		fmt.Sprintf("%x", pathBytes),
		"Unexpected PathBytes\n")
}

func Test_PathGeneration2(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 2)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("Path: %x\n", pathBytes)

	assert.Equal(
		t,
		20,
		len(pathBytes),
		"PathBytes has wrong length: %x, expected length: %x\n", pathBytes, 40)

	assert.Equal(
		t,
		"2c0000807b000080000000000000000000000000",
		fmt.Sprintf("%x", pathBytes),
		"Unexpected PathBytes\n")
}

func Test_PathGeneration3(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 3)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	fmt.Printf("Path: %x\n", pathBytes)

	assert.Equal(
		t,
		20,
		len(pathBytes),
		"PathBytes has wrong length: %x, expected length: %x\n", pathBytes, 40)

	assert.Equal(
		t,
		"2c0000807b000080000000800000000000000000",
		fmt.Sprintf("%x", pathBytes),
		"Unexpected PathBytes\n")
}

func Test_ChunkGeneration(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 0)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	message := getDummyTx()

	chunks, err := prepareChunks(pathBytes, []byte(coinContext), message, userMessageChunkSize)

	assert.Equal(
		t,
		chunks[0],
		pathBytes,
		"First chunk should be pathBytes\n")
}

func Test_ChunkGeneration2(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 0)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	context := []byte{0xa1, 0xa2, 0xa3, 0xa4}
	message := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	println(hex.EncodeToString(context))
	println(hex.EncodeToString(message))

	chunks, err := prepareChunks(pathBytes, context, message, 4)

	for index, c := range chunks {
		println(index, hex.EncodeToString(c))
	}

	assert.Equal(
		t,
		5,
		len(chunks),
		"incorrect number of chunks\n")
}

func Test_ChunkGeneration_invalidContextLength(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 0)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	message := getDummyTx()

	coinContext := strings.Repeat("A", 256)

	_, errChunk := prepareChunks(pathBytes, []byte(coinContext), message, userMessageChunkSize)

	fmt.Printf("Error: %s\n", errChunk)

	assert.Error(t, errChunk)
}

func Test_ChunkGeneration_contextLengthIsZero(t *testing.T) {
	bip44Path := []uint32{44, 123, 0, 0, 0}

	pathBytes, err := GetBip44bytes(bip44Path, 0)
	if err != nil {
		t.Fatalf("Detected error, err: %s\n", err.Error())
	}

	message := getDummyTx()

	coinContext := ""

	chunks, _ := prepareChunks(pathBytes, []byte(coinContext), message, userMessageChunkSize)

	assert.Equal(
		t,
		byte(0),
		chunks[1][0],
		"First byte should be 0 because context is empty\n")
}
