package main

import (
	"fmt"
	"os"

	"1inch-test-case/internal/commands"
	"1inch-test-case/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	cfg := &config.Config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		logger.Error("config process failed", zap.Error(err))
		os.Exit(1)
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", cfg.InfuraKey))
	if err != nil {
		logger.Error("dial infura failed", zap.Error(err))
		os.Exit(1)
	}

	u := commands.NewUni(logger.Named("uni cmd"), client)
	rootCmd.RunE = u.RunE

	err = rootCmd.Execute()
	if err != nil {
		logger.Error("execute failed", zap.Error(err))
		os.Exit(1)
	}
}
