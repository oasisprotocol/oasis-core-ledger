package ledger_oasis_go

import (
	"encoding/binary"
	"fmt"
)

const (
	userMessageChunkSize = 250

	LogModuleName = "ledger-signer"
)

// VersionInfo contains app version information
type VersionInfo struct {
	AppMode uint8
	Major   uint8
	Minor   uint8
	Patch   uint8
}

func (c VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", c.Major, c.Minor, c.Patch)
}

// VersionRequiredError the command is not supported by this app
type VersionRequiredError struct {
	Found    VersionInfo
	Required VersionInfo
}

func (e VersionRequiredError) Error() string {
	return fmt.Sprintf("App Version required %s - Version found: %s", e.Required, e.Found)
}

func NewVersionRequiredError(req, ver VersionInfo) error {
	return &VersionRequiredError{
		Found:    ver,
		Required: req,
	}
}

// CheckVersion compares the current version with the required version
func CheckVersion(ver, req VersionInfo) error {
	if ver.Major != req.Major {
		if ver.Major > req.Major {
			return nil
		}
		return NewVersionRequiredError(req, ver)
	}

	if ver.Minor != req.Minor {
		if ver.Minor > req.Minor {
			return nil
		}
		return NewVersionRequiredError(req, ver)
	}

	if ver.Patch >= req.Patch {
		return nil
	}
	return NewVersionRequiredError(req, ver)
}

func GetBip44bytes(bip44Path []uint32, hardenCount int) ([]byte, error) {
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

	packetIndex := 0

	contextSizeByte := []byte{byte(len(context))}
	body := append(contextSizeByte, context...)
	body = append(body, transaction...)

	packetCount := 1 + len(body)/chunkSize
	if len(body)%chunkSize > 0 {
		packetCount += 1
	}

	chunks := make([][]byte, packetCount)

	// First chunk is path
	chunks[0] = bip44PathBytes
	packetIndex++

	chunkIndex := 0
	for chunkIndex < packetCount-1 {
		start := chunkIndex * chunkSize
		end := (chunkIndex + 1) * chunkSize
		if end > len(body) {
			end = len(body)
		}

		chunks[1+chunkIndex] = body[start:end]
		chunkIndex++
	}

	return chunks, nil
}
