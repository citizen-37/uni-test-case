package pool

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

var ErrInvalidTokens = errors.New("invalid tokens")

type (
	Pool struct {
		logger   *zap.Logger
		contract Contract
	}

	Request struct {
		From   common.Address
		To     common.Address
		Amount *big.Int
	}
)

func NewSwapper(logger *zap.Logger, contract Contract) *Pool {
	return &Pool{
		logger:   logger,
		contract: contract,
	}
}

func (s *Pool) GetAmountOut(ctx context.Context, req Request) (*big.Int, error) {
	s.logger.Debug("got request", zap.Any("request", req))

	token0, token1, err := s.loadTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("load tokens failed: %w", err)
	}

	s.logger.Debug("tokens loaded", zap.Any("token0", token0), zap.Any("token1", token1))

	if !s.validateTokens(req, []common.Address{token0, token1}) {
		return nil, ErrInvalidTokens
	}

	res0, res1, err := s.getReserves(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reserves failed: %w", err)
	}

	s.logger.Debug("got reserves", zap.Stringer("res0", res0), zap.Stringer("res1", res1))

	resFrom, resTo := s.sortReserves(res0, res1, token0, req.From)

	s.logger.Debug("sorted reserves", zap.Stringer("res_from", resFrom), zap.Stringer("res_to", resTo))

	return s.calculateAmountOut(req.Amount, resFrom, resTo), nil
}

func (s *Pool) loadTokens(ctx context.Context) (common.Address, common.Address, error) {
	token0, err := s.contract.Token0(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [20]byte{}, [20]byte{}, fmt.Errorf("load token 0 failed: %w", err)
	}

	token1, err := s.contract.Token1(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [20]byte{}, [20]byte{}, fmt.Errorf("load token 1 failed: %w", err)
	}

	return token0, token1, nil
}

func (s *Pool) validateTokens(req Request, list []common.Address) bool {
	if !contains(req.From, list) {
		return false
	}

	if !contains(req.To, list) {
		return false
	}

	return true
}

func (s *Pool) getReserves(ctx context.Context) (*big.Int, *big.Int, error) {
	res, err := s.contract.GetReserves(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("load reserves failed: %w", err)
	}

	return res.Reserve0, res.Reserve1, nil
}

func (s *Pool) sortReserves(res0, res1 *big.Int, token0, from common.Address) (*big.Int, *big.Int) {
	if token0 == from {
		return res0, res1
	}

	return res1, res0
}

func (s *Pool) calculateAmountOut(amountIn, reserveIn, reserveOut *big.Int) *big.Int {
	numerator, denominator, amountInWithFee, result := &big.Int{}, &big.Int{}, &big.Int{}, &big.Int{}

	amountInWithFee.Mul(amountIn, big.NewInt(997))
	numerator.Mul(amountInWithFee, reserveOut)

	denominator.Mul(reserveIn, big.NewInt(1000))
	denominator.Add(denominator, amountInWithFee)

	result.Div(numerator, denominator)

	return result
}

func contains(address common.Address, list []common.Address) bool {
	for _, a := range list {
		if address == a {
			return true
		}
	}

	return false
}
