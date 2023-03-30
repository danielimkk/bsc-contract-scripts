package main

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/danielimkk/bsc-contract-scripts/app/config"
	"github.com/danielimkk/bsc-contract-scripts/app/contract"
)

func main() {
	cfg := config.NewConfig()

	contractService := contract.NewContractService(cfg)

	contractService.Start()

	oldAddress := common.HexToAddress(cfg.ValidatorAddress)
	oldPrivateKey := cfg.OwnerPrKey

	tx, err := contractService.UnjailValidator(context.Background(), oldAddress, oldPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Validator unjailed succesfully: ", tx.Hash().Hex())

}
