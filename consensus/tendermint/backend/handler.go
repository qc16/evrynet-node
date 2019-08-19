package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
)

//HandleMsg implement consensus.Handler interface
//TODO: implement this.
func (be *backend) HandleMsg(address common.Address, data p2p.Msg) (bool, error) {
	fmt.Printf("not implemented yet, but still beable to receive messsages: (%s)", data.String())
	return true, nil
}
