package config

import (
	"github.com/namsral/flag"
)

type Config struct {
	StakingContractAddress string `yaml:"stakingContractAddress"`
	BinanceSmartChainUrl   string `yaml:"binanceSmartChainUrl"`
	ValidatorAddress       string `yaml:"validatorAddress"`
	OwnerPrKey             string `yaml:"ownerPrKey"`
	Method                 string `yaml:"method"`
}

func NewConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.BinanceSmartChainUrl, "rpc-url", "https://rpc.bombchain.com", "RPC URL of the chain") // https://bombchain-testnet.ankr.com/bas_full_rpc_1  or   https://rpc.bombchain.com
	flag.StringVar(&config.ValidatorAddress, "validator-address", "0xe5024719F3000066E1D4283A590C3446C27a6E32", "Public Key of the Validator address")
	flag.StringVar(&config.OwnerPrKey, "owner-private-key", "0000000000000000000000000000000000000000000000000000000000000000", "Private Key of the Owner address")
	flag.StringVar(&config.StakingContractAddress, "staking-contract", "0x0000000000000000000000000000000000001000", "Address of the Staking system smart contract")
	flag.StringVar(&config.Method, "method", "unjail", "Contract method to call, currently only unjail is available")
	flag.Parse()
	return config
}
