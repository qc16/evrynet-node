package tests_utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func MustGeneratePrivateKey(key string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		panic(err)
	}
	return privateKey
}

// ------------------------------------

func MakeNodeKey() *ecdsa.PrivateKey {
	key, _ := GeneratePrivateKey()
	return key
}

func MustCreateStateDB(t *testing.T) *state.StateDB {
	var (
		statedb, err = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	)
	if err != nil {
		t.Fatalf("failed to create stateDB, error %s", err)

	}
	return statedb
}


func MakeGenesisHeader(validators []common.Address) *types.Header {
	var header = &types.Header{
		Number:     big.NewInt(int64(0)),
		ParentHash: common.HexToHash("0x01"),
		UncleHash:  types.CalcUncleHash(nil),
		Root:       common.HexToHash("0x0"),
		Difficulty: big.NewInt(1),
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := prepareExtra(header)

	var buf bytes.Buffer
	buf.Write(extra[:types.TendermintExtraVanity])
	tdm := &types.TendermintExtra{
		Validators:    validators,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	payload, _ := rlp.EncodeToBytes(&tdm)

	header.Extra = append(buf.Bytes(), payload...)
	return header
}

func MakeBlockWithoutSeal(pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	return types.NewBlockWithHeader(header)
}

func MakeBlockWithSeal(be tendermint.Backend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSealInvalid(be tendermint.Backend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	invalidCommitSeal := make([]byte, types.TendermintExtraSeal)
	_, err := rand.Read(invalidCommitSeal)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, invalidCommitSeal)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSeal(be tendermint.Backend, pHeader *types.Header, validators []common.Address) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	commitHash := utils.PrepareCommittedSeal(header.Hash())
	committedSeal, err := be.Sign(commitHash)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, committedSeal)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

//AppendSeal sign the header with the engine's key and write the seal to the input header's extra data
func AppendSeal(header *types.Header, be tendermint.Backend) {
	// sign the hash
	seal, _ := be.Sign(utils.SigHash(header).Bytes())
	utils.WriteSeal(header, seal)
}

//appendCommittedSeal
func appendCommittedSeal(header *types.Header, committedSeal []byte) {
	//TODO: make this logic as the same as AppendSeal, which involve signing commit before writeCommittedSeal
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	utils.WriteCommittedSeals(header, committedSeals)
}

//makeHeaderFromParent return a new block With valid information from its parents.
func makeHeaderFromParent(parent *types.Block) *types.Header {
	header := &types.Header{
		Coinbase:   GetAddress(),
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Difficulty: big.NewInt(1),
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := prepareExtra(header)
	header.Extra = extra
	return header
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough TendermintExtraVanity bytes.
	if len(header.Extra) < types.TendermintExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.TendermintExtraVanity])

	tdm := &types.TendermintExtra{}
	payload, err := rlp.EncodeToBytes(&tdm)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}
