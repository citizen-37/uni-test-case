package commands

import (
	"context"
	"fmt"
	"math/big"

	"1inch-test-case/internal/contracts"
	"1inch-test-case/internal/pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Uni struct {
	logger *zap.Logger
	client *ethclient.Client
}

func NewUni(logger *zap.Logger, client *ethclient.Client) *Uni {
	return &Uni{
		logger: logger,
		client: client,
	}
}

func (u *Uni) RunE(cmd *cobra.Command, _ []string) error {
	ctx := context.TODO()

	cont, err := contracts.NewPair(poolAddress(cmd), u.client)
	if err != nil {
		return fmt.Errorf("load contract failed: %w", err)
	}

	result, err := pool.NewSwapper(u.logger.Named("swapper"), cont).GetAmountOut(ctx, pool.Request{
		From:   fromAddress(cmd),
		To:     toAddress(cmd),
		Amount: amountIn(cmd),
	})
	if err != nil {
		return fmt.Errorf("get amount out failed: %w", err)
	}

	u.logger.Info("result", zap.Stringer("result", result))
	return nil
}

func poolAddress(cmd *cobra.Command) common.Address {
	return common.HexToAddress(cmd.Flag("pool").Value.String())
}

func fromAddress(cmd *cobra.Command) common.Address {
	return common.HexToAddress(cmd.Flag("from").Value.String())
}

func toAddress(cmd *cobra.Command) common.Address {
	return common.HexToAddress(cmd.Flag("to").Value.String())
}

func amountIn(cmd *cobra.Command) *big.Int {
	bigAmount := &big.Int{}
	bigAmount.SetString(cmd.Flag("amount").Value.String(), 10)

	return bigAmount
}
