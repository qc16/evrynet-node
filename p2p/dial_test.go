// Copyright 2015 The go-ethereum Authors
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

package p2p

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/internal/testlog"
	"github.com/Evrynetlabs/evrynet-node/log"
	"github.com/Evrynetlabs/evrynet-node/p2p/enode"
	"github.com/Evrynetlabs/evrynet-node/p2p/enr"
	"github.com/Evrynetlabs/evrynet-node/p2p/netutil"
)

func init() {
	spew.Config.Indent = "\t"
}

type dialtest struct {
	init   *dialstate // state before and after the test.
	rounds []round
}

type round struct {
	peers []*Peer // current peer set
	done  []task  // tasks that got done this round
	new   []task  // the result must match this one
}

func runDialTest(t *testing.T, test dialtest) {
	var (
		vtime   time.Time
		running int
	)
	pm := func(ps []*Peer) map[enode.ID]*Peer {
		m := make(map[enode.ID]*Peer)
		for _, p := range ps {
			m[p.ID()] = p
		}
		return m
	}
	for i, round := range test.rounds {
		for _, task := range round.done {
			running--
			if running < 0 {
				panic("running task counter underflow")
			}
			test.init.taskDone(task, vtime)
		}

		new := test.init.newTasks(running, pm(round.peers), nil, vtime)
		if !sametasks(new, round.new) {
			t.Errorf("ERROR round %d: \ngot %v\nwant %v\nstate: %v\nrunning: %v",
				i, spew.Sdump(new), spew.Sdump(round.new), spew.Sdump(test.init), spew.Sdump(running))
		}
		t.Logf("round %d new tasks: %s", i, strings.TrimSpace(spew.Sdump(new)))

		// Time advances by 16 seconds on every round.
		vtime = vtime.Add(16 * time.Second)
		running += len(new)
	}
}

type fakeTable []*enode.Node

func (t fakeTable) Self() *enode.Node                     { return new(enode.Node) }
func (t fakeTable) Close()                                {}
func (t fakeTable) LookupRandom() []*enode.Node           { return nil }
func (t fakeTable) Resolve(*enode.Node) *enode.Node       { return nil }
func (t fakeTable) ReadRandomNodes(buf []*enode.Node) int { return copy(buf, t) }
func (t fakeTable) LookupDiscoveredPeers() map[common.Address]*enode.Node {
	return map[common.Address]*enode.Node{}
}

// This test checks that dynamic dials are launched from discovery results.
func TestDialStateDynDial(t *testing.T) {
	privateKeys := generatePrivateKeys(8)
	config := &Config{
		Logger:     testlog.Logger(t, log.LvlTrace),
		PrivateKey: privateKeys[0],
	}
	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, fakeTable{}, 5, config),
		rounds: []round{
			// A discovery query is launched.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				new: []task{&discoverTask{}},
			},
			// Dynamic dials are launched when it completes.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				done: []task{
					&discoverTask{results: []*enode.Node{
						newNode(2, nil, privateKeys), // this one is already connected and not dialed.
						newNode(3, nil, privateKeys),
						newNode(4, nil, privateKeys),
						newNode(5, nil, privateKeys),
						newNode(6, nil, privateKeys), // these are not tried because max dyn dials is 5
						newNode(7, nil, privateKeys), // ...
					}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
				},
			},
			// Some of the dials complete but no new ones are launched yet because
			// the sum of active dial count and dynamic peer count is == maxDynDials.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(4, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
				},
			},
			// No new dial tasks are launched in the this round because
			// maxDynDials has been reached.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(4, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(5, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
				},
				new: []task{
					&waitExpireTask{Duration: 19 * time.Second},
				},
			},
			// In this round, the peer with id 2 drops off. The query
			// results from last discovery lookup are reused.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(4, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(5, nil, privateKeys)}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(6, nil, privateKeys)},
				},
			},
			// More peers (3,4) drop off and dial for ID 6 completes.
			// The last query result from the discovery lookup is reused
			// and a new one is spawned because more candidates are needed.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(5, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(6, nil, privateKeys)},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(7, nil, privateKeys)},
					&discoverTask{},
				},
			},
			// Peer 7 is connected, but there still aren't enough dynamic peers
			// (4 out of 5). However, a discovery is already running, so ensure
			// no new is started.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(5, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(7, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(7, nil, privateKeys)},
				},
			},
			// Finish the running node discovery with an empty set. A new lookup
			// should be immediately requested.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(0, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(5, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(7, nil, privateKeys)}},
				},
				done: []task{
					&discoverTask{},
				},
				new: []task{
					&discoverTask{},
				},
			},
		},
	})
}

// Tests that bootnodes are dialed if no peers are connectd, but not otherwise.
func TestDialStateDynDialBootnode(t *testing.T) {
	privateKeys := generatePrivateKeys(9)
	config := &Config{
		PrivateKey: privateKeys[0],
		BootstrapNodes: []*enode.Node{
			newNode(1, nil, privateKeys),
			newNode(2, nil, privateKeys),
			newNode(3, nil, privateKeys),
		},
		Logger: testlog.Logger(t, log.LvlTrace),
	}
	table := fakeTable{
		newNode(4, nil, privateKeys),
		newNode(5, nil, privateKeys),
		newNode(6, nil, privateKeys),
		newNode(7, nil, privateKeys),
		newNode(8, nil, privateKeys),
	}
	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, table, 5, config),
		rounds: []round{
			// 2 dynamic dials attempted, bootnodes pending fallback interval
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
					&discoverTask{},
				},
			},
			// No dials succeed, bootnodes still pending fallback interval
			{
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
				},
			},
			// No dials succeed, bootnodes still pending fallback interval
			{},
			// No dials succeed, 2 dynamic dials attempted and 1 bootnode too as fallback interval was reached
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(1, nil, privateKeys)},
				},
			},
			// No dials succeed, 2nd bootnode is attempted
			{
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(1, nil, privateKeys)},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(2, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
				},
			},
			// No dials succeed, 3rd bootnode is attempted
			{
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(2, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
				},
			},
			// No dials succeed, 1st bootnode is attempted again, expired random nodes retried
			{
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
				},
				new: []task{},
			},
			// Random dial succeeds, no more bootnodes are attempted
			{
				new: []task{
					&waitExpireTask{3 * time.Second},
				},
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(4, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(1, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
				},
			},
		},
	})
}

func TestDialStateDynDialFromTable(t *testing.T) {
	privateKeys := generatePrivateKeys(13)
	// This table always returns the same random nodes
	// in the order given below.
	table := fakeTable{
		newNode(1, nil, privateKeys),
		newNode(2, nil, privateKeys),
		newNode(3, nil, privateKeys),
		newNode(4, nil, privateKeys),
		newNode(5, nil, privateKeys),
		newNode(6, nil, privateKeys),
		newNode(7, nil, privateKeys),
		newNode(8, nil, privateKeys),
	}

	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, table, 10, &Config{
			PrivateKey: privateKeys[0],
			Logger:     testlog.Logger(t, log.LvlTrace),
		}),
		rounds: []round{
			// 5 out of 8 of the nodes returned by ReadRandomNodes are dialed.
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(1, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(2, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
					&discoverTask{},
				},
			},
			// Dialing nodes 1,2 succeeds. Dials from the lookup are launched.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(1, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(2, nil, privateKeys)},
					&discoverTask{results: []*enode.Node{
						newNode(10, nil, privateKeys),
						newNode(11, nil, privateKeys),
						newNode(12, nil, privateKeys),
					}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(10, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(11, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(12, nil, privateKeys)},
					&discoverTask{},
				},
			},
			// Dialing nodes 3,4,5 fails. The dials from the lookup succeed.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(10, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(11, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(12, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: newNode(3, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(5, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(10, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(11, nil, privateKeys)},
					&dialTask{flags: dynDialedConn, dest: newNode(12, nil, privateKeys)},
				},
			},
			// Waiting for expiry. No waitExpireTask is launched because the
			// discovery query is still running.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(10, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(11, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(12, nil, privateKeys)}},
				},
			},
			// Nodes 3,4 are not tried again because only the first two
			// returned random nodes (nodes 1,2) are tried and they're
			// already connected.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(10, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(11, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(12, nil, privateKeys)}},
				},
			},
		},
	})
}

func newNode(id int, ip net.IP, priKeys []*ecdsa.PrivateKey) *enode.Node {
	var r enr.Record
	if ip != nil {
		r.Set(enr.IP(ip))
	}
	if err := enode.SignV4(&r, priKeys[id]); err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	node, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return nil
	}
	return node
}

// This test checks that candidates that do not match the netrestrict list are not dialed.
func TestDialStateNetRestrict(t *testing.T) {
	privateKeys := generatePrivateKeys(9)
	// This table always returns the same random nodes
	// in the order given below.
	table := fakeTable{
		newNode(1, net.ParseIP("127.0.0.1"), privateKeys),
		newNode(2, net.ParseIP("127.0.0.2"), privateKeys),
		newNode(3, net.ParseIP("127.0.0.3"), privateKeys),
		newNode(4, net.ParseIP("127.0.0.4"), privateKeys),
		newNode(5, net.ParseIP("127.0.2.5"), privateKeys),
		newNode(6, net.ParseIP("127.0.2.6"), privateKeys),
		newNode(7, net.ParseIP("127.0.2.7"), privateKeys),
		newNode(8, net.ParseIP("127.0.2.8"), privateKeys),
	}
	restrict := new(netutil.Netlist)
	restrict.Add("127.0.2.0/24")

	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, table, 10, &Config{
			NetRestrict: restrict,
			PrivateKey:  privateKeys[0],
		}),
		rounds: []round{
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: table[4]},
					&discoverTask{},
				},
			},
		},
	})
}

// This test checks that static dials are launched.
func TestDialStateStaticDial(t *testing.T) {
	privateKeys := generatePrivateKeys(6)
	config := &Config{
		PrivateKey: privateKeys[0],
		StaticNodes: []*enode.Node{
			newNode(1, nil, privateKeys),
			newNode(2, nil, privateKeys),
			newNode(3, nil, privateKeys),
			newNode(4, nil, privateKeys),
			newNode(5, nil, privateKeys),
		},
		Logger: testlog.Logger(t, log.LvlTrace),
	}
	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, fakeTable{}, 0, config),
		rounds: []round{
			// Static dials are launched for the nodes that
			// aren't yet connected.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(3, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(5, nil, privateKeys)},
				},
			},
			// No new tasks are launched in this round because all static
			// nodes are either connected or still being dialed.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(3, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(3, nil, privateKeys)},
				},
			},
			// No new dial tasks are launched because all static
			// nodes are now connected.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(4, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(5, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(4, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(5, nil, privateKeys)},
				},
				new: []task{
					&waitExpireTask{Duration: 19 * time.Second},
				},
			},
			// Wait a round for dial history to expire, no new tasks should spawn.
			{
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: dynDialedConn, node: newNode(2, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(4, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(5, nil, privateKeys)}},
				},
			},
			// If a static node is dropped, it should be immediately redialed,
			// irrespective whether it was originally static or dynamic.
			{
				done: []task{
					&waitExpireTask{Duration: 19 * time.Second},
				},
				peers: []*Peer{
					{rw: &conn{flags: dynDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(3, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(5, nil, privateKeys)}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(2, nil, privateKeys)},
				},
			},
		},
	})
}

// This test checks that past dials are not retried for some time.
func TestDialStateCache(t *testing.T) {
	privateKeys := generatePrivateKeys(4)
	config := &Config{
		PrivateKey: privateKeys[0],
		StaticNodes: []*enode.Node{
			newNode(1, nil, privateKeys),
			newNode(2, nil, privateKeys),
			newNode(3, nil, privateKeys),
		},
		Logger: testlog.Logger(t, log.LvlTrace),
	}
	runDialTest(t, dialtest{
		init: newDialState(enode.ID{}, fakeTable{}, 0, config),
		rounds: []round{
			// Static dials are launched for the nodes that
			// aren't yet connected.
			{
				peers: nil,
				new: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(1, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(2, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(3, nil, privateKeys)},
				},
			},
			// No new tasks are launched in this round because all static
			// nodes are either connected or still being dialed.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(1, nil, privateKeys)},
					&dialTask{flags: staticDialedConn, dest: newNode(2, nil, privateKeys)},
				},
			},
			// A salvage task is launched to wait for node 3's history
			// entry to expire.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(3, nil, privateKeys)},
				},
				new: []task{
					&waitExpireTask{Duration: 19 * time.Second},
				},
			},
			// Still waiting for node 3's entry to expire in the cache.
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(2, nil, privateKeys)}},
				},
			},
			{
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(2, nil, privateKeys)}},
				},
			},
			// The cache entry for node 3 has expired and is retried.
			{
				done: []task{
					&waitExpireTask{Duration: 19 * time.Second},
				},
				peers: []*Peer{
					{rw: &conn{flags: staticDialedConn, node: newNode(1, nil, privateKeys)}},
					{rw: &conn{flags: staticDialedConn, node: newNode(2, nil, privateKeys)}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: newNode(3, nil, privateKeys)},
				},
			},
		},
	})
}

func TestDialResolve(t *testing.T) {
	privateKeys := generatePrivateKeys(2)
	config := &Config{
		PrivateKey: privateKeys[0],
		Logger:     testlog.Logger(t, log.LvlTrace),
		Dialer:     TCPDialer{&net.Dialer{Deadline: time.Now().Add(-5 * time.Minute)}},
	}
	resolved := newNode(1, net.IP{127, 0, 55, 234}, privateKeys)
	table := &resolveMock{answer: resolved}
	state := newDialState(enode.ID{}, table, 0, config)

	// Check that the task is generated with an incomplete ID.
	dest := newNode(1, nil, privateKeys)
	state.addStatic(dest)
	tasks := state.newTasks(0, nil, nil, time.Time{})
	if !reflect.DeepEqual(tasks, []task{&dialTask{flags: staticDialedConn, dest: dest}}) {
		t.Fatalf("expected dial task, got %#v", tasks)
	}

	// Now run the task, it should resolve the ID once.
	srv := &Server{ntab: table, log: config.Logger, Config: *config}
	tasks[0].Do(srv)
	if !reflect.DeepEqual(table.resolveCalls, []*enode.Node{dest}) {
		t.Fatalf("wrong resolve calls, got %v", table.resolveCalls)
	}

	// Report it as done to the dialer, which should update the static node record.
	state.taskDone(tasks[0], time.Now())
	if state.static[dest.ID()].dest != resolved {
		t.Fatalf("state.dest not updated")
	}
}

// compares task lists but doesn't care about the order.
func sametasks(a, b []task) bool {
	if len(a) != len(b) {
		return false
	}
next:
	for _, ta := range a {
		for _, tb := range b {
			if reflect.DeepEqual(ta, tb) {
				continue next
			}
		}
		return false
	}
	return true
}

func uintID(i uint32) enode.ID {
	var id enode.ID
	binary.BigEndian.PutUint32(id[:], i)
	return id
}

// implements discoverTable for TestDialResolve
type resolveMock struct {
	resolveCalls []*enode.Node
	answer       *enode.Node
}

func (t *resolveMock) Resolve(n *enode.Node) *enode.Node {
	t.resolveCalls = append(t.resolveCalls, n)
	return t.answer
}

func (t *resolveMock) Self() *enode.Node                     { return new(enode.Node) }
func (t *resolveMock) Close()                                {}
func (t *resolveMock) LookupRandom() []*enode.Node           { return nil }
func (t *resolveMock) ReadRandomNodes(buf []*enode.Node) int { return 0 }
func (t *resolveMock) LookupDiscoveredPeers() map[common.Address]*enode.Node {
	return map[common.Address]*enode.Node{}
}

func generatePrivateKeys(num int) []*ecdsa.PrivateKey {
	var privateKeys []*ecdsa.PrivateKey
	for i := 0; i < num; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic("Failed to generate private key. Error: " + err.Error())
		}
		privateKeys = append(privateKeys, privateKey)
	}
	return privateKeys
}
