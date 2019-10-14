// Copyright 2017 The go-ethereum Authors
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

package core

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/log"
)

// func (c *core) handleRequest(request *istanbul.Request) error {
// 	logger := c.logger.New("state", c.state, "seq", c.current.sequence)

// 	if err := c.checkRequestMsg(request); err != nil {
// 		if err == errInvalidMessage {
// 			logger.Warn("invalid request")
// 			return err
// 		}
// 		logger.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.Hash())
// 		return err
// 	}

// 	logger.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

// 	c.current.pendingRequest = request
// 	if c.state == StateAcceptRequest {
// 		c.sendPreprepare(request)
// 	}
// 	return nil
// }

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(msg message, blockNumber float32) error {
	state := c.CurrentState()

	if blockNumber < float32(state.BlockNumber().Int64()) {
		return errOldMessage
	}
	if blockNumber > float32(state.BlockNumber().Int64()) {
		return errFutureMessage
	}
	return nil
}

func (c *core) storeRequestMsg(msg message, blockNumber *big.Int) {

	log.Info("Store future request", "address", msg.Address, "blockNumber", blockNumber.Int64())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(msg, float32(blockNumber.Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(message)
		if !ok {
			log.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r, prio)
		if err != nil {
			if err == errFutureMessage {
				log.Info("Stop processing request", "number", prio)
				c.pendingRequests.Push(m, prio)
				break
			}
			log.Info("Skip the pending request", "number", prio)
			continue
		}
		log.Info("Post pending request", "number", prio)

		payload, err := c.FinalizeMsg(&r)
		if err != nil {
			log.Warn("Cannot finalize msg")
			continue
		}
		if err := c.backend.Broadcast(c.valSet, payload); err != nil {
			log.Warn("Failed to Broadcast mgs", "payload", payload)
			continue
		}
	}
}
