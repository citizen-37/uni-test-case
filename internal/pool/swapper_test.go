package pool_test

import (
	"context"
	"math/big"
	"testing"

	"1inch-test-case/internal/pool"
	mocks "1inch-test-case/mocks/internal_/pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func TestSwapper(t *testing.T) {
	suite.Run(t, &SwapperSuite{})
}

type SwapperSuite struct {
	suite.Suite
	contract *mocks.Contract
}

func (s *SwapperSuite) SetupTest() {
	token0 := addr("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	token1 := addr("0xdac17f958d2ee523a2206206994597c13d831ec7")

	contract := mocks.NewContract(s.T())
	contract.On("Token0", mock.Anything).Return(token0, nil).Maybe()
	contract.On("Token1", mock.Anything).Return(token1, nil).Maybe()

	s.contract = contract
}

func (s *SwapperSuite) TestGetAmountOut_all_good() {
	res0 := intString("17923792376848098321135")
	res1 := intString("30004336940045")
	expected := intString("1668880080")

	s.contract.On("GetReserves", mock.Anything).Return(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}{
		Reserve0:           res0,
		Reserve1:           res1,
		BlockTimestampLast: 0,
	}, nil)

	sw := pool.NewSwapper(zap.NewNop(), s.contract)

	result, err := sw.GetAmountOut(context.TODO(), pool.Request{
		From:   addr("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		To:     addr("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		Amount: intString("1000000000000000000"),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, result)
}

func (s *SwapperSuite) TestGetAmountOut_all_good_reversed() {
	res0 := intString("17923792376848098321135")
	res1 := intString("30004336940045")
	expected := intString("17923252983346018099788")

	s.contract.On("GetReserves", mock.Anything).Return(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}{
		Reserve0:           res0,
		Reserve1:           res1,
		BlockTimestampLast: 0,
	}, nil)

	sw := pool.NewSwapper(zap.NewNop(), s.contract)

	result, err := sw.GetAmountOut(context.TODO(), pool.Request{
		// addresses swapped
		From:   addr("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		To:     addr("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		Amount: intString("1000000000000000000"),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, result)
}

func (s *SwapperSuite) TestGetAmountOut_invalid_tokens() {
	sw := pool.NewSwapper(zap.NewNop(), s.contract)

	_, err := sw.GetAmountOut(context.TODO(), pool.Request{
		From:   addr("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		To:     addr("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		Amount: intString("1000000000000000000"),
	})
	require.ErrorIs(s.T(), err, pool.ErrInvalidTokens)
}

func intString(src string) *big.Int {
	result := &big.Int{}
	result.SetString(src, 10)
	return result
}

func addr(src string) common.Address {
	return common.HexToAddress(src)
}
