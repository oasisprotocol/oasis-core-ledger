package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// VersionInfo contains app version information.
type VersionInfo struct {
	AppMode uint8
	Major   uint8
	Minor   uint8
	Patch   uint8
}

func (c VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", c.Major, c.Minor, c.Patch)
}

// VersionRequiredError the command is not supported by this app.
type VersionRequiredError struct {
	Found    VersionInfo
	Required VersionInfo
}

func (e VersionRequiredError) Error() string {
	return fmt.Sprintf("App Version required %s - Version found: %s", e.Required, e.Found)
}

func newVersionRequiredError(req, ver VersionInfo) error {
	return &VersionRequiredError{
		Found:    ver,
		Required: req,
	}
}

func checkVersion(ver, req VersionInfo) error {
	if ver.Major != req.Major {
		if ver.Major > req.Major {
			return nil
		}
		return newVersionRequiredError(req, ver)
	}

	if ver.Minor != req.Minor {
		if ver.Minor > req.Minor {
			return nil
		}
		return newVersionRequiredError(req, ver)
	}

	if ver.Patch >= req.Patch {
		return nil
	}
	return newVersionRequiredError(req, ver)
}

func getBip44bytes(bip44Path []uint32, hardenCount int) ([]byte, error) {
	message := make([]byte, 20)
	if len(bip44Path) != 5 {
		return nil, fmt.Errorf("path should contain 5 elements")
	}
	for index, element := range bip44Path {
		pos := index * 4
		value := element
		if index < hardenCount {
			value = 0x80000000 | element
		}
		binary.LittleEndian.PutUint32(message[pos:], value)
	}
	return message, nil
}

func prepareChunks(bip44PathBytes, context, transaction []byte, chunkSize int) ([][]byte, error) {
	if len(context) > 255 {
		return nil, fmt.Errorf("maximum supported context size is 255 bytes")
	}

	contextSizeByte := []byte{byte(len(context))}
	body := append(contextSizeByte, context...)
	body = append(body, transaction...)

	packetCount := 1 + len(body)/chunkSize
	if len(body)%chunkSize > 0 {
		packetCount++
	}

	chunks := make([][]byte, 0, packetCount)
	chunks = append(chunks, bip44PathBytes) // First chunk is path.

	r := bytes.NewReader(body)
readLoop:
	for {
		toAppend := make([]byte, chunkSize)
		n, err := r.Read(toAppend)
		if n > 0 {
			// Note: n == 0 only when EOF.
			chunks = append(chunks, toAppend[:n])
		}
		switch err {
		case nil:
		case io.EOF:
			break readLoop
		default:
			// This can never happen, but handle it.
			return nil, err
		}
	}

	return chunks, nil
}
