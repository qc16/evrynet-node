package types

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// IstanbulDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine
	IstanbulDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	IstanbulExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidIstanbulHeaderExtra = errors.New("invalid istanbul header extra-data")
)

type IstanbulExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// ExtractIstanbulExtra extracts all values of the IstanbulExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractIstanbulExtra(h *Header) (*IstanbulExtra, error) {
	if len(h.Extra) < IstanbulExtraVanity {
		return nil, ErrInvalidIstanbulHeaderExtra
	}

	var istanbulExtra *IstanbulExtra
	err := rlp.DecodeBytes(h.Extra[IstanbulExtraVanity:], &istanbulExtra)
	if err != nil {
		return nil, err
	}
	return istanbulExtra, nil
}
