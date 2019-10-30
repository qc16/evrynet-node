package backend

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func TestSaveAndLoad(t *testing.T) {
	var (
		hash   = common.HexToHash("1234567890")
		valSet = validator.NewSet([]common.Address{
			common.HexToAddress("1234567894"),
			common.HexToAddress("1234567895"),
		}, tendermint.RoundRobin, int64(0))
		snap = newSnapshot(5, 10, hash, valSet)
	)

	db := ethdb.NewMemDatabase()
	assert.NoError(t, snap.store(db))

	snap1, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	assert.NoError(t, err)
	assert.NotNil(t, snap1)
	assert.Equal(t, snap.Epoch, snap1.Epoch)
	assert.Equal(t, snap.Hash, snap1.Hash)
	if !reflect.DeepEqual(snap.ValSet, snap.ValSet) {
		t.Errorf("validator set mismatch: have %v, want %v", snap1.ValSet, snap.ValSet)
	}
}

func TestApplyHeaders(t *testing.T) {
	var (
		hash     = common.HexToHash("1234567890")
		addr1    = common.HexToAddress("1")
		addr2    = common.HexToAddress("2")
		addr3    = common.HexToAddress("3")
		addr4    = common.HexToAddress("4")
		fakeAddr = common.HexToAddress("99")
		addrs    = []common.Address{
			addr1, addr2, addr3, addr4,
		}
		valSet = validator.NewSet(addrs, tendermint.RoundRobin, int64(1))

		//init snapshot with epoch is 5, number is 0
		snap    = newSnapshot(5, 1, hash, valSet)
		headers []*types.Header
	)

	if val := snap.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr1) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr1)
	}

	// got error unauthorized with fake proposal
	headers = createHeaderArr(2, 1, []common.Address{fakeAddr})
	_, err := snap.apply(headers)
	assert.Equal(t, err, errUnauthorized)

	//Apply new 2 headers, new proposer should be addr2
	headers = createHeaderArr(2, 2, addrs)
	snap2, err := snap.apply(headers)
	assert.NoError(t, err)

	if val := snap2.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr3) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr3)
	}

	//Apply new 5 headers, new proposer should be addr2
	headers = createHeaderArr(2, 5, addrs)
	snap3, err := snap.apply(headers)
	assert.NoError(t, err)

	if val := snap3.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr2) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr2)
	}

	// test save and load after that apply headers
	db := ethdb.NewMemDatabase()
	assert.NoError(t, snap.store(db))
	snap4, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	if val := snap4.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr1) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr1)
	}

	//Apply new 9 headers, new proposer should be addr2
	headers = createHeaderArr(2, 9, addrs)
	snap5, err := snap4.apply(headers)
	assert.NoError(t, err)

	if val := snap5.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr2) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr2)
	}

	//vote kick out addr1
	headers = createHeaderArr(2, 4, addrs)
	applyVotes(headers, addr1, false)
	snap6, err := snap4.apply(headers)
	assert.NoError(t, err)
	assert.Equal(t, 3, snap6.ValSet.Size())

	//vote kick in fakeAddr
	headers = createHeaderArr(2, 4, addrs)
	applyVotes(headers, fakeAddr, true)
	snap7, err := snap4.apply(headers)
	assert.NoError(t, err)
	assert.Equal(t, 5, snap7.ValSet.Size())
}

func applyVotes(headers []*types.Header, addr common.Address, vote bool) {
	for _, header := range headers {
		extra, _ := prepareExtraWithModified(header, addr)
		header.Extra = extra
		if vote {
			copy(header.Nonce[:], tendermint.NonceAuthVote)
		} else {
			copy(header.Nonce[:], tendermint.NonceDropVote)
		}
	}
}

func createHeaderArr(startNumber int, countHeader int, addrs []common.Address) []*types.Header {
	headers := make([]*types.Header, countHeader)
	for i := startNumber; i < startNumber+countHeader; i++ {
		index := i % len(addrs)
		header := &types.Header{
			Number:   big.NewInt(int64(i)),
			Coinbase: addrs[index],
		}
		tdm := &types.TendermintExtra{}
		payload, _ := rlp.EncodeToBytes(&tdm)
		tendermintExtraVanity := bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity)
		header.Extra = append(tendermintExtraVanity, payload...)

		headers[i-startNumber] = header
	}
	return headers
}

func prepareExtraWithModified(header *types.Header, modifiedAddress common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough TendermintExtraVanity bytes.
	if len(header.Extra) < types.TendermintExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.TendermintExtraVanity])

	tdm := &types.TendermintExtra{
		ModifiedValidator: modifiedAddress,
	}
	payload, err := rlp.EncodeToBytes(&tdm)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}
