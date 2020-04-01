package tests_utils

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/Evrynetlabs/evrynet-node/event"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/utils"
	"github.com/Evrynetlabs/evrynet-node/core/rawdb"
	"github.com/Evrynetlabs/evrynet-node/core/state"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/rlp"
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
	extra, _ := PrepareExtra(header)

	var buf bytes.Buffer
	buf.Write(extra[:types.TendermintExtraVanity])
	valSetData, _ := rlp.EncodeToBytes(validators)
	tdm := &types.TendermintExtra{
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
		ValidatorAdds: valSetData,
	}
	payload, _ := rlp.EncodeToBytes(&tdm)

	header.Extra = append(buf.Bytes(), payload...)
	return header
}

func MakeBlockWithoutSeal(pHeader *types.Header) *types.Block {
	header := MakeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	return types.NewBlockWithHeader(header)
}

//AppendCommittedSeal
func AppendCommittedSeal(header *types.Header, committedSeal []byte) {
	//TODO: make this logic as the same as AppendSeal, which involve signing commit before writeCommittedSeal
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	_ = utils.WriteCommittedSeals(header, committedSeals)
}

//makeHeaderFromParent return a new block With valid information from its parents.
func MakeHeaderFromParent(parent *types.Block) *types.Header {
	header := &types.Header{
		Coinbase:   GetAddress(),
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   parent.GasLimit(),
		GasUsed:    0,
		Difficulty: big.NewInt(1),
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := PrepareExtra(header)
	header.Extra = extra
	return header
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// PrepareExtra returns a extra-data of the given header and validators
func PrepareExtra(header *types.Header) ([]byte, error) {
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

func GenerateMockChainReader() (*MockChainReader, error) {
	var (
		nodePrivateKey = MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{nodeAddr}
		genesisHeader  = MakeGenesisHeader(validators)
		stateDB, err   = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	)
	if err != nil {
		return nil, errors.New("failed to create stateDB, error: " + err.Error())
	}

	return &MockChainReader{
		GenesisHeader: genesisHeader,
		MockBlockChain: &MockBlockChain{
			Statedb:          stateDB,
			GasLimit:         1000000000,
			ChainHeadFeed:    new(event.Feed),
			MockCurrentBlock: types.NewBlockWithHeader(genesisHeader),
		},
	}, nil
}
