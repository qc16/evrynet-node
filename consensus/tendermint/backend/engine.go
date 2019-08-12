package backend

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// post block into tendermint engine
	go sb.EventMux().Post(tendermint.RequestEvent{
		Proposal: block,
	})
	return nil
}

// Start implements consensus.Tendermint.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return tendermint.ErrStartedEngine
	}

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.Istanbul.Stop
func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return tendermint.ErrStoppedEngine
	}
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

func (sb *backend) Author(header *types.Header) (common.Address, error) {
	panic("Author: implement me")
}

func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Warn("VerifyHeader: implement me")
	return nil
}

func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("VerifyHeaders: implement me")
}

func (sb *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("VerifyUncles: implement me")
}

func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	panic("VerifySeal: implement me")
}

func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Warn("Prepare: implement me")
	return nil
}

func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	log.Warn("Finalize: implement me")
}

func (sb *backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	log.Warn("FinalizeAndAssemble: implement me")

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

func (sb *backend) SealHash(header *types.Header) (hash common.Hash) {
	log.Warn("SealHash: implement me")

	//TODO: this logic is temporary.
	// I wanna make hash is different when SealHash() was called to bypass `func (w *worker) taskLoop()`
	hasher := sha3.NewLegacyKeccak256()
	rd := rand.Uint64()
	rlp.Encode(hasher, rd)
	hasher.Sum(hash[:0])
	return hash
}

func (sb *backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	log.Warn("CalcDifficulty: implement me")
	return nil
}

func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	log.Warn("APIs: implement me")
	return nil
}

func (sb *backend) Close() error {
	log.Warn("Close: implement me")
	return nil
}
