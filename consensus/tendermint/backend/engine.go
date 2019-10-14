package backend

import (
	"bytes"
	"errors"
	"math/big"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/evrynet-official/evrynet-client/rpc"
)

var (
	// TendermintBlockReward tempo fix the Block reward in wei for successfully mining a block
	// TODO: will modify after
	TendermintBlockReward = big.NewInt(5e+18)

	defaultDifficulty = big.NewInt(1)
	now               = time.Now
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
)

var (
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	errInvalidSignature = errors.New("invalid signature")
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	errInvalidExtraDataFormat = errors.New("invalid extra data format")
	// errInvalidMixDigest is returned if a block's mix digest is not Tendermint digest.
	errInvalidMixDigest = errors.New("invalid Tendermint mix digest")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	errInvalidCommittedSeals = errors.New("invalid committed seals")
	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	errEmptyCommittedSeals = errors.New("zero committed seals")
	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")
	// errCoinBaseInvalid is returned if the value of coin base is not equals proposer's address in header
	errCoinBaseInvalid = errors.New("invalid coin base address")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
	// errMalformedChannelData is returned if data return from blockFinalization does not conform to its struct definition
	errMalformedChannelData = errors.New("data received is not an event type")
)

func (sb *backend) addProposalSeal(h *types.Header) error {
	seal, err := sb.Sign(utils.SigHash(h).Bytes())
	if err != nil {
		return err
	}
	utils.WriteSeal(h, seal)
	return nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	blockNumber := header.Number.Uint64()

	// validate address of the validator
	// get snapshot
	snap, err := sb.snapshot(chain, blockNumber-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if err = sb.addProposalSeal(header); err != nil {
		return err
	}
	block = block.WithSeal(header)
	// checks the address must stored in snapshot
	if _, v := snap.ValSet.GetByAddress(sb.address); v == nil {
		return errUnauthorized
	}

	// update seal to header of the block
	parent := chain.GetHeader(header.ParentHash, blockNumber-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// wait for the timestamp of header, make sure this block does not come from the future
	headerTime := int64(block.Header().Time)
	delay := time.Unix(headerTime, 0).Sub(now())
	//TODO: revise delay
	select {
	case <-time.After(delay):
	case <-stop:
		results <- nil
		return
	}
	blockNumberStr := block.Number().String()

	if _, ok := sb.commitChs[blockNumberStr]; !ok {
		sb.commitChs[blockNumberStr] = make(chan *types.Block, 1)
	}
	log.Info("sealing...", "total number of channels", len(sb.commitChs))
	//block = sb.Prepare()
	//TODO: clear previous data of proposal
	// post block into tendermint engine
	go func(block *types.Block) {
		sb.EventMux().Post(tendermint.NewBlockEvent{
			Block: block,
		})
	}(block)
	// miner won't be able to interrupt a sealing task
	// a sealing task can only exist when core consensus agreed upon a block
	go func(blockNumberStr string) {
		ch := sb.commitChs[blockNumberStr]
		//TODO: DO we need timeout for consensus?
		for {
			select {
			case bl, ok := <-ch:
				if !ok {
					log.Info("committing... Channel closed, exit seal...", "number", blockNumberStr)
					return
				}
				if bl.Number().String() != blockNumberStr {
					log.Warn("committing.. Received a different block number than the sealing block number", "received", bl.Number().String(), "expected", blockNumberStr)
				}
				//this step is to stop other go routine wait for a block
				close(ch)
				delete(sb.commitChs, bl.Number().String())
				if bl == nil {
					log.Error("committing... Received nil ")
					return
				}

				//we only posted the block back to the miner if and only if the block is ours
				if bl.Coinbase() == sb.address {
					log.Info("committing... returned block to miner", "block_hash", bl.Hash(), "number", bl.Number())
					results <- bl
				} else {
					log.Info("committing... not this node's block, exit and let downloader sync the block from proposer...", "block_hash", block.Hash(), "number", block.Number())
				}
				return
				//case <-stop:
				//log.Warn("committing... refused to exit because the sealing task might be the finalize block. The seal only exit when core commit a block", "number", block.Number())
			}
		}
	}(blockNumberStr)
	return nil
}

// Start implements consensus.Tendermint.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return tendermint.ErrStartedEngine
	}

	//set chain reader
	sb.chain = chain
	sb.currentBlock = currentBlock

	if sb.commitChs != nil {
		for blockNo, chn := range sb.commitChs {
			close(chn)
			delete(sb.commitChs, blockNo)
		}
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

// Author retrieves the Evrynet address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	return blockProposer(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time > big.NewInt(now().Unix()).Uint64() {
		return consensus.ErrFutureBlock
	}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractTendermintExtra(header); err != nil {
		return errInvalidExtraDataFormat
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.TendermintDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	return sb.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// get block number from header of block
	blockNumber := header.Number.Uint64()
	if blockNumber == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, blockNumber-1)
	}
	if parent == nil || parent.Number.Uint64() != blockNumber-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time+sb.config.BlockPeriod > header.Time {
		//TODO: find out if tendermint is subject to error when Block Period is too fast
		//	return errInvalidTimestamp
		log.Warn("block time difference is too small", "different in ms", header.Time-sb.config.BlockPeriod)
	}

	// get snap shoot to prepare for the verify proposal and committed seal
	snap, err := sb.snapshot(chain, blockNumber-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	if err := sb.verifyProposalSeal(header, snap); err != nil {
		return err
	}

	return sb.verifyCommittedSeals(header, snap)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	errorHeaders := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := sb.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case errorHeaders <- err:
			}
		}
	}()
	return abort, errorHeaders
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (sb *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	blockNumber := header.Number.Uint64()
	if blockNumber == 0 {
		return errUnknownBlock
	}

	// get snap shoot to prepare for the verify proposal
	snap, err := sb.snapshot(chain, blockNumber-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	return sb.verifyProposalSeal(header, snap)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// set coinbase with the proposer's address
	header.Coinbase = sb.Address()
	// use the same difficulty and mixDigest for all blocks
	header.MixDigest = types.TendermintDigest
	// use the same difficulty for all blocks
	// TODO: thight might reflect 2F+1 value since our block have nothing included to indicate it
	header.Difficulty = sb.CalcDifficulty(chain, header.Time, nil)

	// get parent
	var (
		blockNumber = header.Number.Uint64()
		parent      = chain.GetHeader(header.ParentHash, blockNumber-1)
	)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// prepare extra data without validators
	extra, err := prepareExtra(header)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp from parent's timestamp and blockperiod
	var (
		parentTime  = new(big.Int).SetUint64(parent.Time)
		blockPeriod = new(big.Int).SetUint64(sb.config.BlockPeriod)
		headerTime  = new(big.Int).Add(parentTime, blockPeriod)
	)

	if headerTime.Int64() < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	} else {
		header.Time = headerTime.Uint64()
	}

	//TODO: modify valset data if epoch is reached.

	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	// Accumulate any block rewards and commit the final state root
	accumulateRewards(chain.Config(), state, header)

	// Since there is a change in stateDB, its trie must be update
	// In case block reached EIP158 hash, the state will attempt to delete empty object as EIP158 sepcification
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
}

// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block rewards and commit the final state root
	accumulateRewards(chain.Config(), state, header)

	// No block rewards, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (sb *backend) SealHash(header *types.Header) (hash common.Hash) {
	return utils.SigHash(header)
}

// CalcDifficulty tempo return default difficulty
func (sb *backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return defaultDifficulty
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

// snapshot retrieves the authorization snapshot at a given point in time.
func (sb *backend) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	// Loop and try to find a valid snapshot that contain the block hash we need, otherwise a list of headers and a
	// most recent snapshot then apply the headers onto that snapshot to get the snapshot we need
	for snap == nil {
		// If an in-memory snapshot was found, use that
		//TODO: get from cached if the snapshot is existed

		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			s, err := loadSnapshot(sb.config.Epoch, sb.db, hash)
			if err != nil {
				log.Warn("cannot load snapshot from db", "error", err)
			} else {
				log.Debug("Loaded voting snapshot form disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at block zero, make a snapshot
		if number == 0 {
			log.Debug("creating snapshot at block 0")
			genesis := chain.GetHeaderByNumber(0)
			if err := sb.VerifyHeader(chain, genesis, false); err != nil {
				return nil, err
			}
			extra, err := types.ExtractTendermintExtra(genesis)
			if err != nil {
				return nil, err
			}
			snap = newSnapshot(sb.config.Epoch, 0, genesis.Hash(), validator.NewSet(extra.Validators, sb.config.ProposerPolicy, int64(0)))
			if err := snap.store(sb.db); err != nil {
				return nil, err
			}
			log.Trace("Stored genesis voting snapshot to disk")
			break
		}
		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}
		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}
	//revert the headers's array index , i.e, block n..1 become 1..n
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}
	// apply the list of headers found on top of it
	snap, err := snap.apply(headers)
	if err != nil {
		return nil, err
	}
	//TODO: add to cached for snapshot

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(sb.db); err != nil {
			return nil, err
		}
		log.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// verifyProposalSeal checks proposal seal is signed by validator
func (sb *backend) verifyProposalSeal(header *types.Header, snap *Snapshot) error {
	// resolve the authorization key and check against signers
	signer, err := blockProposer(header)
	if err != nil {
		log.Error("proposal seal is invalid", "error", err)
		return err
	}
	// compare with coin base that contain the address of proposer.
	if signer != header.Coinbase {
		return errCoinBaseInvalid
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.ValSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(header *types.Header, snap *Snapshot) error {
	extra, err := types.ExtractTendermintExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return errEmptyCommittedSeals
	}

	vals := snap.ValSet.Copy()
	// Check whether the committed seals are generated by parent's validators
	validSeal := 0
	proposalSeal := tendermintCore.PrepareCommittedSeal(header.Hash())
	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := utils.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			log.Error("not a valid address", "err", err)
			return errInvalidSignature
		}
		// Every validator can have only one seal. If more than one seals are signed by a
		// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
		if vals.RemoveValidator(addr) {
			validSeal++
		} else {
			return errInvalidCommittedSeals
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= 2*snap.ValSet.F() {
		return errInvalidCommittedSeals
	}

	return nil
}

// blockProposer extracts the Ethereum account address from a signed header.
func blockProposer(header *types.Header) (common.Address, error) {
	//TODO: check if existed in the cached

	// Retrieve the signature from the header extra-data
	extra, err := types.ExtractTendermintExtra(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := utils.GetSignatureAddress(utils.SigHash(header).Bytes(), extra.Seal)
	if err != nil {
		return addr, err
	}
	//TODO: will be caching address
	return addr, nil
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

// AccumulateRewards credits the coinbase of the given block with the proposing
// reward.
func accumulateRewards(config *params.ChainConfig, state *state.StateDB, header *types.Header) {
	// Accumulate the rewards for the proposer
	reward := new(big.Int).Set(TendermintBlockReward)

	state.AddBalance(header.Coinbase, reward)
}
