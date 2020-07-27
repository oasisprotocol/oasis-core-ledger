package internal

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var coinContext = "oasis-core/consensus: tx for chain 7b02d647e8997bacebce96723f6904029ec78b67c261c4bdddb5e47de1ab31fa"

func getDummyTx() []byte {
	base64tx := "pGNmZWWiY2dhcxkD6GZhbW91bnRCB9BkYm9keaJneGZlcl90b1UA4ywoibwEEhHt7fqvlNL9hmmLsH9reGZlcl90b2tlbnNFJ5TKJABlbm9uY2UHZm1ldGhvZHBzdGFraW5nLlRyYW5zZmVy" //nolint: lll
	tx, _ := base64.StdEncoding.DecodeString(base64tx)
	return tx
}

func TestPrintVersion(t *testing.T) {
	require := require.New(t)

	reqVersion := VersionInfo{0, 1, 2, 3}
	s := fmt.Sprintf("%v", reqVersion)
	require.Equal("1.2.3", s)
}

func TestPathGeneration(t *testing.T) {
	const expectedLength = 20

	type pathGenerationTest struct {
		name        string
		path        []uint32
		hardenCount int
		expected    string
	}

	for _, testCase := range []pathGenerationTest{
		{
			name:        "TC0",
			path:        []uint32{44, 100, 0, 0, 0},
			hardenCount: 0,
			expected:    "2c00000064000000000000000000000000000000",
		},
		{
			name:        "TC2",
			path:        []uint32{44, 123, 0, 0, 0},
			hardenCount: 2,
			expected:    "2c0000807b000080000000000000000000000000",
		},
		{
			name:        "TC3",
			path:        []uint32{44, 123, 0, 0, 0},
			hardenCount: 3,
			expected:    "2c0000807b000080000000800000000000000000",
		},
	} {
		tc := testCase // Shut up scopelint.
		t.Run(testCase.name, func(tt *testing.T) {
			require := require.New(tt)

			pathBytes, err := getBip44bytes(tc.path, tc.hardenCount)
			require.NoError(err, "GetBip44bytes")

			tt.Logf("Path: %x", pathBytes)

			require.Len(pathBytes, expectedLength, "pathBytes has the correct length")
			require.Equal(tc.expected, fmt.Sprintf("%x", pathBytes), "Generated path is correct")
		})
	}
}

func validateChunks(
	t *testing.T,
	require *require.Assertions,
	pathBytes, context, message []byte,
	chunks [][]byte,
	chunkSize int,
) {
	expected := append(pathBytes, byte(len(context)))
	expected = append(expected, context...)
	expected = append(expected, message...)

	nrChunks := len(chunks)

	reassembled := make([]byte, 0, len(expected))
	for i, v := range chunks {
		if i == 0 {
			require.Equal(pathBytes, v, "First chunk should be pathBytes")
		} else {
			if i < nrChunks-1 {
				require.Len(v, chunkSize, "Non-final/non-path chunks should be fully packed")
			} else {
				require.NotZero(len(v), "Final chunk should have non-zero length")
				require.LessOrEqualf(len(v), chunkSize, "Chunk[%d]: invalid length")
			}
		}

		reassembled = append(reassembled, v...)

		t.Logf("Chunk[%d]: %x", i, v)
	}

	require.Equal(expected, reassembled, "Reconstructed message should match")
}

func TestChunkGeneration(t *testing.T) {
	require := require.New(t)

	bip44Path := []uint32{44, 123, 0, 0, 0}
	pathBytes, err := getBip44bytes(bip44Path, 0)
	require.NoError(err, "GetBip44bytes")

	context := []byte(coinContext)
	message := getDummyTx()

	chunks, err := prepareChunks(pathBytes, context, message, userMessageChunkSize)
	require.NoError(err, "prepareChunks")

	validateChunks(t, require, pathBytes, context, message, chunks, userMessageChunkSize)
}

func TestChunkGeneration2(t *testing.T) {
	require := require.New(t)

	bip44Path := []uint32{44, 123, 0, 0, 0}
	pathBytes, err := getBip44bytes(bip44Path, 0)
	require.NoError(err, "GetBip44bytes")

	context := []byte{0xa1, 0xa2, 0xa3, 0xa4}
	message := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	t.Logf("Context: %x", context)
	t.Logf("Message: %x", message)

	const chunkSize = 4
	chunks, err := prepareChunks(pathBytes, context, message, chunkSize)
	require.NoError(err, "prepareChunks")

	validateChunks(t, require, pathBytes, context, message, chunks, chunkSize)
	require.Len(chunks, 5, "incorrect number of chunks")
}

func TestChunkGenerationInvalidContextLength(t *testing.T) {
	require := require.New(t)

	bip44Path := []uint32{44, 123, 0, 0, 0}
	pathBytes, err := getBip44bytes(bip44Path, 0)
	require.NoError(err, "GetBip44bytes")

	message := getDummyTx()
	context := []byte(strings.Repeat("A", 256))

	_, err = prepareChunks(pathBytes, context, message, userMessageChunkSize)
	require.Error(err, "oversized context should fail")

	t.Logf("Error: %v", err)
}

func TestChunkGenerationContextLengthIsZero(t *testing.T) {
	require := require.New(t)

	bip44Path := []uint32{44, 123, 0, 0, 0}
	pathBytes, err := getBip44bytes(bip44Path, 0)
	require.NoError(err, "GetBip44bytes")

	message := getDummyTx()
	context := []byte{}

	chunks, err := prepareChunks(pathBytes, context, message, userMessageChunkSize)
	require.NoError(err, "prepareChunks")

	validateChunks(t, require, pathBytes, context, message, chunks, userMessageChunkSize)
	require.Zero(chunks[1][0], "First non-path byte should be 0 because context is empty")
}
