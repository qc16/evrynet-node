package core

import (
	"io"
	"math/big"
	"strconv"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/rlp"
)

//Proposal represent a propose message to be sent in the case of the node is a proposer
//for its Round.
type Proposal struct {
	Block    *types.Block
	Round    int64
	POLRound int64
}

func (p *Proposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		p.Block,
		strconv.FormatInt(p.Round, 10),
		strconv.FormatInt(p.POLRound, 10),
	})
}

func (p *Proposal) DecodeRLP(s *rlp.Stream) error {
	var ps struct {
		Block   *types.Block
		RStr    string
		POLRStr string
	}
	if err := s.Decode(&ps); err != nil {
		return err
	}
	round, err := strconv.ParseInt(ps.RStr, 10, 64)
	if err != nil {
		return err
	}
	polcr, err := strconv.ParseInt(ps.POLRStr, 10, 64)
	if err != nil {
		return err
	}
	p.Block = ps.Block
	p.Round = round
	p.POLRound = polcr
	return nil
}

type Vote struct {
	BlockHash   *common.Hash
	BlockNumber *big.Int
	Round       int64
	Seal        []byte
}

func (v *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		v.BlockHash,
		v.BlockNumber,
		strconv.FormatInt(v.Round, 10),
		v.Seal,
	})
}

func (v *Vote) DecodeRLP(s *rlp.Stream) error {
	var vs struct {
		BlockHash   *common.Hash
		BlockNumber *big.Int
		RStr        string
		Seal        []byte
	}
	if err := s.Decode(&vs); err != nil {
		return err
	}
	round, err := strconv.ParseInt(vs.RStr, 10, 64)
	if err != nil {
		return err
	}
	v.BlockHash = vs.BlockHash
	v.BlockNumber = vs.BlockNumber
	v.Round = round
	v.Seal = vs.Seal
	return nil
}

type CatchUpMsg struct {
	BlockNumber *big.Int
	Round       int64
	Step        RoundStepType
}

func (msg *CatchUpMsg) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		msg.BlockNumber,
		uint64(msg.Round),
		msg.Step,
	})
}

func (msg *CatchUpMsg) DecodeRLP(s *rlp.Stream) error {
	var vs struct {
		BlockNumber *big.Int
		Round       uint64
		Step        RoundStepType
	}
	if err := s.Decode(&vs); err != nil {
		return err
	}
	msg.BlockNumber = vs.BlockNumber
	msg.Round = int64(vs.Round)
	msg.Step = vs.Step
	return nil
}
