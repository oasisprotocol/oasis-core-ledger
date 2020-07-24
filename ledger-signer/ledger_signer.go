// Package main implements the ledger backed oasis-node signer plugin.
package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"

	"github.com/oasisprotocol/oasis-core-ledger/internal"
)

var (
	_ signature.SignerFactoryCtor = newPluginFactory
	_ signature.SignerFactory     = (*ledgerFactory)(nil)
	_ signature.Signer            = (*ledgerSigner)(nil)

	// signerPathCoinType is set to 474, the number associated with Oasis ROSE.
	signerPathCoinType uint32 = 474
	// signerPathAccount is the account index used to sign transactions.
	signerPathAccount uint32 = 0
	// SignerPathChange indicates an external chain.
	signerPathChange uint32 = 0

	// signerEntityDerivationRootPath is the BIP-0032 path prefix used for generating
	// an Entity signer.
	signerEntityDerivationRootPath = []uint32{
		internal.PathPurposeBIP44,
		signerPathCoinType,
		signerPathAccount,
		signerPathChange,
	}
	// signerConsensusDerivationRootPath is the derivation path prefix used for
	// generating a consensus signer.
	signerConsensusDerivationRootPath = []uint32{
		internal.PathPurposeConsensus,
		signerPathCoinType,
		internal.PathSubPurposeConsensus,
		signerPathAccount,
	}

	roleDerivationRootPaths = map[signature.SignerRole][]uint32{
		signature.SignerEntity:    signerEntityDerivationRootPath,
		signature.SignerConsensus: signerConsensusDerivationRootPath,
	}
)

// GetPluginCtor is the plugin's entry point.
func GetPluginCtor() signature.SignerFactoryCtor {
	return newPluginFactory
}

type factoryConfig struct {
	address string
	index   uint32
}

func newFactoryConfig(cfgStr string) (*factoryConfig, error) {
	kvStrs := strings.Split(cfgStr, ";")

	var (
		cfg                      factoryConfig
		foundAddress, foundIndex bool
	)
	for _, v := range kvStrs {
		spl := strings.Split(v, "=")
		if len(spl) != 2 {
			return nil, fmt.Errorf("malformed k/v pair: '%s'", v)
		}

		key := strings.ToLower(spl[0])
		switch key {
		case "address":
			if foundAddress {
				return nil, fmt.Errorf("address already configured")
			}
			cfg.address = spl[1]
			foundAddress = true
		case "index":
			if foundIndex {
				return nil, fmt.Errorf("index already configured")
			}
			idx, err := strconv.ParseUint(spl[1], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("malformed index: %w", err)
			}
			cfg.index = uint32(idx)
			foundIndex = true
		default:
			return nil, fmt.Errorf("unknown configuration option: '%v'", spl[0])
		}
	}

	if !foundAddress {
		return nil, fmt.Errorf("address not configured")
	}
	if !foundIndex {
		return nil, fmt.Errorf("index not configured")
	}

	return &cfg, nil
}

func newPluginFactory(config interface{}, roles ...signature.SignerRole) (signature.SignerFactory, error) {
	cfgStr, ok := config.(string)
	if !ok {
		return nil, fmt.Errorf("ledger: invalid configuration type")
	}
	cfg, err := newFactoryConfig(cfgStr)
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to parse configuration: %w", err)
	}

	return &ledgerFactory{
		roles:   roles,
		address: cfg.address,
		index:   cfg.index,
	}, nil
}

type ledgerFactory struct {
	roles   []signature.SignerRole
	address string
	index   uint32
}

func (fac *ledgerFactory) EnsureRole(role signature.SignerRole) error {
	for _, v := range fac.roles {
		if v == role {
			return nil
		}
	}
	return signature.ErrRoleMismatch
}

func (fac *ledgerFactory) Generate(role signature.SignerRole, _rng io.Reader) (signature.Signer, error) {
	// Generate has the same functionality as Load, since all keys are
	// generated on the Ledger device.
	return fac.Load(role)
}

func (fac *ledgerFactory) Load(role signature.SignerRole) (signature.Signer, error) {
	pathPrefix, ok := roleDerivationRootPaths[role]
	if !ok {
		return nil, fmt.Errorf("ledger: role %d is not supported when using the Ledger backed signer", role)
	}
	path := append(pathPrefix, fac.index)
	device, err := internal.ConnectLedgerOasisApp(fac.address, path)
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to connect to device: %w", err)
	}

	return &ledgerSigner{device, path, nil}, nil
}

type ledgerSigner struct {
	device    *internal.LedgerOasis
	path      []uint32
	publicKey *signature.PublicKey
}

func (s *ledgerSigner) Public() signature.PublicKey {
	if s.publicKey != nil {
		return *s.publicKey
	}

	var pubKey signature.PublicKey
	retrieved, err := s.device.GetPublicKeyEd25519(s.path)
	if err != nil {
		panic(fmt.Errorf("ledger: failed to retrieve public key from device: %w", err))
	}
	if err = pubKey.UnmarshalBinary(retrieved); err != nil {
		panic(fmt.Errorf("ledger: device returned malfored public key: %w", err))
	}
	s.publicKey = &pubKey

	return pubKey
}

// ContextSign generates a signature with the private key over the context and
// message.
func (s *ledgerSigner) ContextSign(context signature.Context, message []byte) ([]byte, error) {
	preparedContext, err := signature.PrepareSignerContext(context)
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to prepare signing context: %w", err)
	}

	signature, err := s.device.SignEd25519(s.path, preparedContext, message)
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to sign message: %w", err)
	}

	return signature, nil
}

// String returns the address of the account on the Ledger device.
func (s *ledgerSigner) String() string {
	return fmt.Sprintf("[ledger signer: %s]", s.Public())
}

// Reset tears down the Signer.
func (s *ledgerSigner) Reset() {
	s.device.Close()
}
