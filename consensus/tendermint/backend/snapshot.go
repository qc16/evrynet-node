package backend

import (
	"bytes"
	"encoding/json"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethdb"
)

const (
	dbKeySnapshotPrefix = "tendermint-snapshot"
)

// Vote represents a single vote that an authorized validator made to modify the
// list of authorizations.
type Vote struct {
	Validator         common.Address `json:"validator"`         // Authorized validator that cast this vote
	BlockNumber       uint64         `json:"blockNumber"`       // Block number the vote was cast in (expire old votes)
	ModifiedValidator common.Address `json:"modifiedValidator"` // Account being voted on to change its authorization
	Authorize         bool           `json:"authorize"`         // Whether to authorize or de-authorize the voted account
}

// Tally is a simple vote tally to keep the current score of votes. Votes that
// go against the proposal aren't counted since it's equivalent to not voting.
type Tally struct {
	Authorize bool `json:"authorize"` // Whether the vote it about authorizing or kicking someone
	Votes     int  `json:"votes"`     // Number of votes until now wanting to pass the proposal
}

// Snapshot is the state of the authorization voting at a given point in time.
// It doesn't have anything to do with voting of creating new Block, only voting for changes in validator sets
type Snapshot struct {
	Epoch  uint64 // The number of blocks after which to checkpoint and reset the pending votes
	Number uint64 // Block number where the snapshot was created

	Votes []*Vote                  // List of votes cast in chronological order
	Tally map[common.Address]Tally // Current vote tally to avoid recalculating

	Hash   common.Hash             // Block hash where the snapshot was created
	ValSet tendermint.ValidatorSet // Set of authorized validators at this moment
}

type snapshotJSON struct {
	Epoch  uint64      `json:"epoch"`
	Number uint64      `json:"number"`
	Hash   common.Hash `json:"hash"`

	Votes []*Vote                  `json:"votes"`
	Tally map[common.Address]Tally `json:"tally"`

	// for validator set
	Validators []common.Address          `json:"validators"`
	Policy     tendermint.ProposerPolicy `json:"policy"`
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:  s.Epoch,
		Number: s.Number,
		Hash:   s.Hash,
		ValSet: s.ValSet.Copy(),
		Votes:  make([]*Vote, len(s.Votes)),
		Tally:  make(map[common.Address]Tally),
	}

	for address, tally := range s.Tally {
		cpy.Tally[address] = tally
	}
	copy(cpy.Votes, s.Votes)

	return cpy
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(headers []*types.Header) (*Snapshot, error) {
	countHeader := len(headers)
	// Allow passing in no headers for cleaner code
	if countHeader == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < countHeader-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errInvalidVotingChain
	}
	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	for _, header := range headers {
		blockNumber := header.Number.Uint64()
		if blockNumber%s.Epoch == 0 {
			// Remove any votes on checkpoint blocks
			snap.Votes = nil
			snap.Tally = make(map[common.Address]Tally)
		}

		// Resolve the authorization key and check against validators
		validator, err := blockProposer(header)
		if err != nil {
			return nil, err
		}
		if _, v := snap.ValSet.GetByAddress(validator); v == nil {
			return nil, errUnauthorized
		}

		// get address of the candidate to validate
		modifiedValidator, err := getModifiedValidator(*header)
		if err != nil {
			return nil, errInvalidCandidate
		}

		if snap.Votes != nil {
			// Header authorized, discard any previous votes from the validator
			for i, vote := range snap.Votes {
				if vote.Validator == validator && vote.ModifiedValidator == modifiedValidator {
					// Un-cast the vote from the cached tally
					snap.uncast(vote.ModifiedValidator, vote.Authorize)

					// Un-cast the vote from the chronological list
					snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
					break // only one vote allowed
				}
			}
		}

		// Tally up the new vote from the validator
		var authorize bool
		switch {
		case bytes.Compare(header.Nonce[:], nonceAuthVote) == 0:
			authorize = true
		case bytes.Compare(header.Nonce[:], nonceDropVote) == 0:
			authorize = false
		default:
			return nil, errInvalidVote
		}

		// cast the validator's vote
		// add to the tally and the vote collectors of snapshot
		if snap.cast(modifiedValidator, authorize) {
			snap.Votes = append(snap.Votes, &Vote{
				Validator:         validator,
				BlockNumber:       blockNumber,
				ModifiedValidator: modifiedValidator,
				Authorize:         authorize,
			})
		}

		// If the vote passed, update the list of validators
		// if the number of votes > 50%
		if tally := snap.Tally[modifiedValidator]; tally.Votes > snap.ValSet.Size()/2 {
			if tally.Authorize {
				// if the authorize is ok add modified validator to valset collectors
				snap.ValSet.AddValidator(modifiedValidator)
			} else {
				snap.ValSet.RemoveValidator(modifiedValidator)

				// Discard any previous votes the de-authorized validator cast
				for i := 0; i < len(snap.Votes); i++ {
					if snap.Votes[i].Validator == modifiedValidator {
						// Un-cast the vote from the cached tally
						snap.uncast(snap.Votes[i].ModifiedValidator, snap.Votes[i].Authorize)

						// Un-cast the vote from the chronological list
						snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)

						i--
					}
				}
			}
			// Discard any previous votes around the just changed account
			for i := 0; i < len(snap.Votes); i++ {
				if snap.Votes[i].ModifiedValidator == modifiedValidator {
					snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
					i--
				}
			}
			// remove tally for new/old validator
			delete(snap.Tally, modifiedValidator)
		}
	}

	//Recalcualte valset
	blockHeightWhenApplyHeaders := snap.ValSet.Height() + int64(countHeader)
	newValSet := validator.NewSet(snap.validators(), snap.ValSet.Policy(), blockHeightWhenApplyHeaders)
	snap.ValSet = newValSet

	snap.Number = snap.Number + uint64(countHeader)
	snap.Hash = headers[countHeader-1].Hash()

	return snap, nil
}

// getModifiedValidator get modified validator in the extra data for cals votes
func getModifiedValidator(header types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	extra, err := types.ExtractTendermintExtra(&header)
	if err != nil {
		return common.Address{}, err
	}
	return extra.ModifiedValidator, nil
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db ethdb.Database) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return db.Put(append([]byte(dbKeySnapshotPrefix), s.Hash[:]...), blob)
}

// checkVote return whether it's a valid vote
func (s *Snapshot) checkVote(address common.Address, authorize bool) bool {
	_, validator := s.ValSet.GetByAddress(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

// cast adds a new vote into the tally.
func (s *Snapshot) cast(address common.Address, authorize bool) bool {
	// Ensure the vote is meaningful
	if !s.checkVote(address, authorize) {
		return false
	}
	// Cast the vote into an existing or new tally
	if old, ok := s.Tally[address]; ok {
		old.Votes++
		s.Tally[address] = old
	} else {
		s.Tally[address] = Tally{Authorize: authorize, Votes: 1}
	}
	return true
}

// uncast removes a previously cast vote from the tally.
func (s *Snapshot) uncast(address common.Address, authorize bool) bool {
	// If there's no tally, it's a dangling vote, just drop
	tally, ok := s.Tally[address]
	if !ok {
		return false
	}
	// Ensure we only revert counted votes
	if tally.Authorize != authorize {
		return false
	}
	// Otherwise revert the vote
	if tally.Votes > 1 {
		tally.Votes--
		s.Tally[address] = tally
	} else {
		delete(s.Tally, address)
	}
	return true
}

// newSnapshot create a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use if for
// the genesis block.
func newSnapshot(epoch uint64, number uint64, hash common.Hash, valSet tendermint.ValidatorSet) *Snapshot {
	snap := &Snapshot{
		Epoch:  epoch,
		Number: number,
		Hash:   hash,
		ValSet: valSet,
		Tally:  make(map[common.Address]Tally),
	}
	return snap
}

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(epoch uint64, db ethdb.Database, hash common.Hash) (*Snapshot, error) {
	blob, err := db.Get(append([]byte(dbKeySnapshotPrefix), hash[:]...))
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal(blob, snap); err != nil {
		return nil, err
	}
	snap.Epoch = epoch

	return snap, nil
}

func (s *Snapshot) toJSONStruct() *snapshotJSON {
	return &snapshotJSON{
		Epoch:      s.Epoch,
		Number:     s.Number,
		Hash:       s.Hash,
		Validators: s.validators(),
		Policy:     s.ValSet.Policy(),
		Votes:      s.Votes,
		Tally:      s.Tally,
	}
}

// UnmarshalJSON from a json byte array
func (s *Snapshot) UnmarshalJSON(b []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	s.Epoch = j.Epoch
	s.Number = j.Number
	s.Hash = j.Hash
	s.Votes = j.Votes
	s.Tally = j.Tally
	s.ValSet = validator.NewSet(j.Validators, j.Policy, int64(j.Number))
	return nil
}

// MarshalJSON to a json byte array
func (s *Snapshot) MarshalJSON() ([]byte, error) {
	j := s.toJSONStruct()
	return json.Marshal(j)
}
