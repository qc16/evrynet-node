package utils

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/sha3"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/rlp"
)

var (
	ErrInvalidSealLength = errors.New("seal is expected to be multiplication of 65")
)

const (
	msgCommit uint64 = iota
)

// sigHash returns the hash
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
func SigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.New256()

	// Clean seal is required for calculating proposer seal.
	_ = rlp.Encode(hasher, types.TendermintFilteredHeader(header, false))
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

// WriteValSet writes the extra-data field of the given header with the given val-sets address.
func WriteValSet(h *types.Header, validators []byte) error {
	tendermintExtra, err := types.ExtractTendermintExtra(h)
	if err != nil {
		return err
	}

	tendermintExtra.ValidatorAdds = validators
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

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}

// GetCheckpointNumber returns check-point base on epoch duration and block-number
func GetCheckpointNumber(epochDuration uint64, blockNumber uint64) uint64 {
	if blockNumber == 0 || blockNumber < epochDuration {
		return 0
	}
	return epochDuration * (blockNumber / epochDuration)
}
