package utils

import (
	"errors"

	"github.com/evrynet-official/evrynet-client/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/rlp"
)

var (
	ErrInvalidSealLength = errors.New("seal is expected to be multiplication of 65")
)

// sigHash returns the hash
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
// if keepSeal = true, the SigHash will return the hash with proposalSeal, otherwise the pure hash
func SigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.New256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.TendermintFilteredHeader(header, false))
	hasher.Sum(hash[:0])

	return hash
}

// WriteSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func WriteSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.TendermintExtraSeal != 0 {
		return ErrInvalidSealLength
	}

	tendermintExtra, err := types.ExtractTendermintExtra(h)
	if err != nil {
		return err
	}

	tendermintExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&tendermintExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.TendermintExtraVanity], payload...)
	return nil
}

// WriteCommittedSeals writes the extra-data field of a block header with given committed seals.
func WriteCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return ErrInvalidSealLength
	}

	for _, seal := range committedSeals {
		if len(seal) != types.TendermintExtraSeal {
			return ErrInvalidSealLength
		}
	}

	tendermintExtra, err := types.ExtractTendermintExtra(h)
	if err != nil {
		return err
	}

	tendermintExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(tendermintExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&tendermintExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.TendermintExtraVanity], payload...)
	return nil
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
