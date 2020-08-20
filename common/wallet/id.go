package wallet

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"errors"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/hash"
)

const (
	// IDSize is the size of a wallet ID in bytes.
	// NOTE: The length of wallet ID's string is IDSize * 2.
	IDSize = 3
)

var (
	// ErrMalformedID is the error returned when a wallet ID is malformed.
	ErrMalformedID = errors.New("wallet: malformed ID")

	_ encoding.BinaryMarshaler   = ID{}
	_ encoding.BinaryUnmarshaler = (*ID)(nil)
)

// ID is a wallet ID computed as a truncated hash of a public key for a specific
// BIP32 path.
type ID [IDSize]byte

// MarshalBinary encodes a wallet ID into binary form.
func (id ID) MarshalBinary() (data []byte, err error) {
	data = append([]byte{}, id[:]...)
	return
}

// UnmarshalBinary decodes a binary marshaled wallet ID.
func (id *ID) UnmarshalBinary(data []byte) error {
	if len(data) != IDSize {
		return ErrMalformedID
	}
	copy(id[:], data)
	return nil
}

// UnmarshalHex deserializes a hexadecimal text string into wallet ID.
func (id *ID) UnmarshalHex(text string) error {
	b, err := hex.DecodeString(text)
	if err != nil {
		return err
	}
	return id.UnmarshalBinary(b)
}

// Equal compares vs another wallet ID for equality.
func (id ID) Equal(cmp ID) bool {
	return bytes.Equal(id[:], cmp[:])
}

// String returns the string representation of a wallet ID.
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// IsValid checks whether a wallet ID is well-formed.
func (id ID) IsValid() bool {
	return len(id) == IDSize
}

// NewID creates a new wallet ID the given data.
func NewID(data []byte) (id ID) {
	h := hash.NewFromBytes(data)
	truncatedHash, err := h.Truncate(IDSize)
	if err != nil {
		panic(err)
	}
	_ = id.UnmarshalBinary(truncatedHash)
	return
}
