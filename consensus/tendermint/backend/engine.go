package backend

import (
	"math/big"

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
	//TODO: need to validate address, block period and update block
	// clear previous data of proposal

	// post block into tendermint engine
	go sb.EventMux().Post(tendermint.RequestEvent{
		Proposal: block,
	})

	//TODO: faking logic to approve the block immediately
	go func() {
		select {
		case results <- block:
		default:
			log.Warn("Sealing result is not read by miner")
		}
	}()

	return nil
}

// Start implements consensus.Tendermint.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return tendermint.ErrStartedEngine
	}

	//TODO: clear previous data of proposal

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.Tendermint.Stop
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

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	panic("Author: implement me")
	//TODO: Research & Implement
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Warn("VerifyHeader: implement me")
	//TODO: Research & Implement
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("VerifyHeaders: implement me")
	//TODO: Research & Implement
}

func (sb *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("VerifyUncles: implement me")
	//TODO: Research & Implement
}

func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	panic("VerifySeal: implement me")
	//TODO: Research & Implement
}

func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	header.Difficulty = big.NewInt(1)
	log.Warn("Prepare: implement me")

	//TODO: Research & Implement
	return nil
}

func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	log.Warn("Finalize: implement me")
	//TODO: Research & Implement
}

func (sb *backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	log.Warn("FinalizeAndAssemble: implement me")
	//TODO: Research & Implement

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

func (sb *backend) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	//TODO: this logic is temporary.
	// I wanna make hash is different when SealHash() was called to bypass `func (w *worker) taskLoop()`

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	})
	hasher.Sum(hash[:0])
	return hash
}

func (sb *backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	log.Warn("CalcDifficulty: implement me")
	//TODO: Research & Implement
	return nil
}

func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	log.Warn("APIs: implement me")
	//TODO: Research & Implement
	return nil
}

func (sb *backend) Close() error {
	log.Warn("Close: implement me")
	//TODO: Research & Implement
	return nil
}
