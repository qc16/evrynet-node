package tests

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/backend"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/p2p/enode"
)

const (
	pkey2 = "AEC5EB6A80CC094363D206949C3ED475C2C5060A23049150310D4FD39F95AF99"
	pkey1 = "CAB57E606531AF83BFD023F55E1673713DA7E161D2612570A0ABAAA9507AACDF"
)

//TestStartingTendermint setup a test to with actual running components of a tendermint consensus
//The test is not finished yet but by running it, the procedure of a tendermint in implementation can be seens
//Current Expectation: if backend isc reated with nodePk1, it will be come proposer of the round and try to send propose message
// 					   if backend is created with nodePk2, it will wait for propose message and timeout
// 					   other logs are printed to indicate flow logic of core's consensus.
func TestStartingTendermint(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

	var (
		nodePk1    = mustGeneratePrivateKey(pkey1)
		nodePk2    = mustGeneratePrivateKey(pkey2)
		tbe1       = backend.New(tendermint.DefaultConfig, nodePk1)
		totalPeers = 2
		n1         = enode.MustParseV4("enode://" + hex.EncodeToString(crypto.FromECDSAPub(&nodePk1.PublicKey)[1:]) + "@33.4.2.1:30303")
		n2         = enode.MustParseV4("enode://" + hex.EncodeToString(crypto.FromECDSAPub(&nodePk2.PublicKey)[1:]) + "@33.4.2.1:30304")
	)
	assert.NoError(t, tbe1.Start(nil, nil))
	pm1, err := eth.NewTestProtocolManagerWithConsensus(tbe1)
	time.Sleep(2 * time.Second)
	assert.NoError(t, err)
	defer pm1.Stop()

	//Create 2 Pipe for read and write. These are full duplex
	io1, io2 := p2p.MsgPipe()
	//p1 will write to io2, p2 will receive from io1 and vice versal.
	p1 := pm1.NewPeer(63, p2p.NewPeerFromNode(n1, fmt.Sprintf("peer %d", 0), nil), io2)
	p2 := pm1.NewPeer(63, p2p.NewPeerFromNode(n2, fmt.Sprintf("peer %d", 1), nil), io1)
	assert.NoError(t, eth.RegisterNewPeer(pm1, p1))
	assert.NoError(t, eth.RegisterNewPeer(pm1, p2))

	//Making sure that the handlingMsg is done by calling pm.handleMsg
	var (
		errCh         = make(chan error, totalPeers)
		doneCh        = make(chan struct{}, totalPeers)
		receivedCount int
		expectedCount = 2
	)
	timeout := time.After(20 * time.Second)
	for _, p := range []*eth.Peer{p1, p2} {
		go func(p *eth.Peer) {
			for {
				if err := pm1.HandleMsg(p); err != nil {
					errCh <- err
				} else {
					doneCh <- struct{}{}
				}
			}
		}(p)
	}
outer:
	for {
		select {
		case err = <-errCh:
			fmt.Printf("handling error %v\n", err)
			break outer
		case <-doneCh:
			receivedCount++
			if receivedCount >= expectedCount {
				fmt.Printf("handling done ")
				break outer
			}

		case <-timeout:
			fmt.Printf("timdeout")

			t.Fail()
			break outer
		}

	}
	if err != nil {
		t.Errorf("error handling msg by peer: %v", err)
	}
}
