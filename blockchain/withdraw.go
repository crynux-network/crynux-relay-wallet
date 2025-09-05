package blockchain

import (
	"context"
	"crynux_relay_wallet/blockchain/bindings"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"database/sql"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
)

func GetWithdrawalFeeAddress(ctx context.Context, network string) (common.Address, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return common.Address{}, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return client.WithdrawContractInstance.GetWithdrawalFeeAddress(opts)
}

func QueueWithdraw(ctx context.Context, db *gorm.DB, address common.Address, amount *big.Int, withdrawalFeeAmount *big.Int, network string) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[network]
	if !ok {
		return nil, ErrBlockchainNotFound
	}

	abi, err := bindings.WithdrawMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("withdraw", address, amount, withdrawalFeeAmount)
	if err != nil {
		return nil, err
	}
	dataStr := hexutil.Encode(data)

	totalAmount := big.NewInt(0).Add(amount, withdrawalFeeAmount)

	transaction := &models.BlockchainTransaction{
		Network:     network,
		Type:        "Withdraw::withdraw",
		Status:      models.TransactionStatusPending,
		FromAddress: blockchain.Account.Address,
		ToAddress:   blockchain.Contracts.Withdraw,
		Value:       totalAmount.String(),
		Data: sql.NullString{
			String: dataStr,
			Valid:  true,
		},
		MaxRetries: blockchain.MaxRetries,
	}

	if err := transaction.Save(ctx, db); err != nil {
		return nil, err
	}

	return transaction, nil
}