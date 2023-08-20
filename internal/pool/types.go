package pool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type (
	Contract interface {
		Token0(opts *bind.CallOpts) (common.Address, error)
		Token1(opts *bind.CallOpts) (common.Address, error)
		GetReserves(opts *bind.CallOpts) (struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}, error)
	}
)
