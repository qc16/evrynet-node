// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package types

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// TendermintDigest TODO: Add digest for our consensus, Digest represents a hash of ""
	// to identify whether the block is from Tendermint consensus engine
	TendermintDigest = common.HexToHash("0x0")

	// TendermintExtraVanity Fixed number of extra-data bytes reserved for validator vanity
	TendermintExtraVanity = 32
	// TendermintExtraSeal Fixed number of extra-data bytes reserved for validator seal
	TendermintExtraSeal = 65

	// ErrInvalidTendermintHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidTendermintHeaderExtra = errors.New("invalid tendermint header extra-data")
)

// TendermintExtra extra data for Tendermint consensus
type TendermintExtra struct {
	LastCommitHash []byte // commit from validators from the last block

	// hashes from the app output from the prev block
	ValidatorsHash     []byte // validators for the current block
	NextValidatorsHash []byte // validators for the next block

	EvidenceHash []byte // evidence of malicious validators included in the block

	Seal          []byte   // Proposer seal 65 bytes
	CommittedSeal [][]byte // Committed seal, 65 * len(Validators) bytes
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *TendermintExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.LastCommitHash,
		ist.ValidatorsHash,
		ist.NextValidatorsHash,
		ist.EvidenceHash,
		ist.Seal,
		ist.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder, and load the tendermint fields from a RLP stream.
func (ist *TendermintExtra) DecodeRLP(s *rlp.Stream) error {
	var tendermintExtra struct {
		LastCommitHash     []byte
		ValidatorsHash     []byte
		NextValidatorsHash []byte
		EvidenceHash       []byte
		Seal               []byte
		CommittedSeal      [][]byte
	}
	if err := s.Decode(&tendermintExtra); err != nil {
		return err
	}
	ist.LastCommitHash, ist.ValidatorsHash, ist.NextValidatorsHash = tendermintExtra.LastCommitHash, tendermintExtra.ValidatorsHash, tendermintExtra.NextValidatorsHash
	ist.EvidenceHash = tendermintExtra.EvidenceHash
	ist.Seal, ist.CommittedSeal = tendermintExtra.Seal, tendermintExtra.CommittedSeal
	return nil
}

// ExtractTendermintExtra extracts all values of the TendermintExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractTendermintExtra(h *Header) (*TendermintExtra, error) {
	if len(h.Extra) < TendermintExtraVanity {
		return nil, ErrInvalidTendermintHeaderExtra
	}

	var tendermintExtra *TendermintExtra
	err := rlp.DecodeBytes(h.Extra[TendermintExtraVanity:], &tendermintExtra)
	if err != nil {
		return nil, err
	}
	return tendermintExtra, nil
}

// TendermintFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Tendermint hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func TendermintFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	tendermintExtra, err := ExtractTendermintExtra(newHeader)
	// Returns nil if ValidatorHash is missing, since a Header is not valid unless there is
	// a ValidatorsHash (corresponding to the validator set).
	if err != nil || len(tendermintExtra.ValidatorsHash) == 0 {
		return nil
	}

	if !keepSeal {
		tendermintExtra.Seal = []byte{}
	}
	tendermintExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&tendermintExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:TendermintExtraVanity], payload...)

	return newHeader
}
