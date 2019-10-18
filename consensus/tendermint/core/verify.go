package core

import (
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
)

// Verify implements tendermint.Backend.Verify
func (c *Core) Verify(proposal tendermint.Proposal) error {
	var (
		block   = proposal.Block
		txs     = block.Transactions()
		txnHash = types.DeriveSha(txs)
	)

	// check block body
	if txnHash != block.Header().TxHash {
		return errMismatchTxhashes
	}

	// Verify transaction for CoreTxPool
	if c.backend.TxPool() != nil && c.backend.TxPool().CoreTxPool != nil {
		for _, t := range txs {
			if err := c.backend.TxPool().CoreTxPool.ValidateTx(t, false); err != nil {
				return err
			}
		}
	}

	// verify the header of proposed block
	err := c.backend.VerifyHeader(c.backend.Chain(), block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == errEmptyCommittedSeals {
		return nil
	}
	return err
}
