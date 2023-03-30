package contract

import (
	"context"
	"math/big"

	"github.com/danielimkk/bsc-contract-scripts/app/config"
	contract "github.com/danielimkk/bsc-contract-scripts/app/contract/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	bsc        *ethclient.Client
	config     *config.Config
	contract   *contract.Staking
	BscChainID *big.Int
}

func NewContractService(
	config *config.Config,
) *Service {
	return &Service{
		config: config,
	}
}

func calcGasPrice(ctx context.Context, ethClient *ethclient.Client) (*big.Int, error) {
	// Check if network is Mainnet
	return gasPrice(ctx, ethClient)
}

func gasPrice(ctx context.Context, ethClient *ethclient.Client) (*big.Int, error) {
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	log.Info("gas price: ", gasPrice)

	return gasPrice, nil
}

func SignTx(tx *types.Transaction, signer types.Signer, pk string) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return signedTx, nil
}

func (s *Service) Start() {
	err := s.connectEthereum()
	if err != nil {
		log.Fatalf("unable to start contract service: %s", err)
	}
	contractAddress := common.HexToAddress(s.config.StakingContractAddress)
	instance, err := contract.NewStaking(contractAddress, s.bsc)
	if err != nil {
		log.Fatal(err)
	}
	s.contract = instance
}

func (s *Service) dialEthClientOrFatal(url string) (*ethclient.Client, *big.Int) {
	dial, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Unable to dial eth1 client for url (%s): %s", url, err)
	}
	chainID, err := dial.ChainID(context.Background())
	if err != nil {
		log.Fatalf("unable to get chain id: %s", err)
	}
	log.Infof("Chain id for url (%s) is %s", url, chainID.String())
	return dial, chainID
}

func (s *Service) connectEthereum() error {
	/* connect to bsc */
	s.bsc, s.BscChainID = s.dialEthClientOrFatal(s.config.BinanceSmartChainUrl)
	return nil
}

func (s *Service) UnjailValidator(ctx context.Context, valAddress common.Address, pk string) (*types.Transaction, error) {

	fromAddress := valAddress

	opts, err := s.getEthTxOptionsWithPk(ctx, fromAddress, pk)
	if err != nil {
		log.Error("Failed to get tx opts. ", err)
		return nil, err
	}
	log.Info("Opts: ", opts.GasFeeCap)

	txHash, err := s.contract.ReleaseValidatorFromJail(opts, valAddress)
	if err != nil {
		log.Error("Failed to send tx. ", err)
		return nil, err
	}

	return txHash, nil
}

func (s *Service) getEthTxOptionsWithPk(ctx context.Context, from common.Address, pk string) (*bind.TransactOpts, error) {
	ethGasPrice, err := calcGasPrice(ctx, s.bsc)
	if err != nil {
		log.Error("Failed to calculate gas price. ", err)
		return nil, err
	}

	balance, err := s.bsc.BalanceAt(ctx, from, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Balance of validator: ", balance)

	return &bind.TransactOpts{
		Context: ctx,
		From:    from,
		Nonce:   nil,
		Signer: func(from common.Address, tx *types.Transaction) (*types.Transaction, error) {
			log.Info("tx gas required: ", tx.Gas())
			log.Info("chain id: ", s.BscChainID)
			signer := types.NewEIP155Signer(s.BscChainID)
			return SignTx(tx, signer, pk)
		},
		GasLimit: 0,
		GasPrice: ethGasPrice,
		Value:    nil,
	}, nil
}
