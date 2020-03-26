package tests_utils

import (
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/tests_utils"
	"github.com/Evrynetlabs/evrynet-node/utils"
)

func MakeBlockWithSeal(be tendermint.Backend, pHeader *types.Header) *types.Block {
	header := tests_utils.MakeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSealInvalid(be tendermint.Backend, pHeader *types.Header) *types.Block {
	header := tests_utils.MakeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	// generate a random private key not in validator set
	privateKey, _ := tests_utils.GeneratePrivateKey()
	//sign header with rand private key
	hashData := crypto.Keccak256(utils.SigHash(header).Bytes())
	invalidCommitSeal, _ := crypto.Sign(hashData, privateKey)
	tests_utils.AppendCommittedSeal(header, invalidCommitSeal)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSeal(be tendermint.Backend, pHeader *types.Header) *types.Block {
	header := tests_utils.MakeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	commitHash := utils.PrepareCommittedSeal(header.Hash())
	committedSeal, err := be.Sign(commitHash)
	if err != nil {
		panic(err)
	}
	tests_utils.AppendCommittedSeal(header, committedSeal)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

//AppendSeal sign the header with the engine's key and write the seal to the input header's extra data
func AppendSeal(header *types.Header, be tendermint.Backend) {
	// sign the hash
	seal, _ := be.Sign(utils.SigHash(header).Bytes())
	_ = utils.WriteSeal(header, seal)
}
