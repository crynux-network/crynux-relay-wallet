package blockchain

import (
	"context"
	"crynux_relay_wallet/models"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TransactionSender sends pending transactions from database to blockchain
type TransactionSender struct {
	db            *gorm.DB
	processingTxs sync.Map
	txQueue       chan *models.BlockchainTransaction
	stopChan      chan struct{}
	isRunning     bool
	batchSize     int
	pollInterval  time.Duration
}

// NewTransactionSender creates a new transaction sender instance
func NewTransactionSender(db *gorm.DB) *TransactionSender {
	return &TransactionSender{
		db:           db,
		txQueue:      make(chan *models.BlockchainTransaction, 100),
		stopChan:     make(chan struct{}),
		isRunning:    false,
		batchSize:    50,
		pollInterval: 5 * time.Second,
	}
}

// Start starts the transaction sender goroutine
func (ts *TransactionSender) Start(ctx context.Context) {
	if ts.isRunning {
		return
	}

	ts.isRunning = true
	go ts.run(ctx)
	log.Info("Transaction sender started")
}

// Stop stops the transaction sender goroutine
func (ts *TransactionSender) Stop() {
	if !ts.isRunning {
		return
	}

	close(ts.stopChan)
	ts.isRunning = false
	log.Info("Transaction sender stopped")
}

// run is the main loop for sending transactions
func (ts *TransactionSender) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go ts.getPendingTransactions(ctx)
	go ts.processPendingTransactions(ctx)

	select {
	case <-ts.stopChan:
		close(ts.txQueue)
		return
	case <-ctx.Done():
		close(ts.txQueue)
		return
	}
}

// getPendingTransactions gets pending transactions and adds them to the queue
func (ts *TransactionSender) getPendingTransactions(ctx context.Context) {
	ticker := time.NewTicker(ts.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get pending transactions that need to be sent
			transactions, err := func(ctx context.Context) ([]models.BlockchainTransaction, error) {
				var allTransactions []models.BlockchainTransaction
				offset := 0
				for {
					transactions, err := models.GetPendingTransactions(ctx, ts.db, offset, ts.batchSize)
					if err != nil {
						return nil, err
					}

					if len(transactions) == 0 {
						break
					}
					allTransactions = append(allTransactions, transactions...)
					offset += len(transactions)
				}
				return allTransactions, nil
			}(ctx)
			if err != nil {
				log.Errorf("Error getting pending transactions: %v", err)
				continue
			}
			if len(transactions) == 0 {
				continue
			}
			var cnt int
			for _, transaction := range transactions {
				_, loaded := ts.processingTxs.LoadOrStore(transaction.ID, struct{}{})
				if !loaded {
					select {
					case <-ctx.Done():
						return
					case ts.txQueue <- &transaction:
						cnt++
					}
				}
			}
			log.Infof("Processing %d pending transactions for sending", cnt)
		}
	}
}

// processPendingTransactions processes transactions from the queue
func (ts *TransactionSender) processPendingTransactions(ctx context.Context) {
	for transaction := range ts.txQueue {
		if err := ts.sendTransaction(ctx, transaction); err != nil {
			log.Errorf("Failed to send transaction %d: %v", transaction.ID, err)
		}
	}
}

// sendTransaction sends a single transaction to the blockchain
func (ts *TransactionSender) sendTransaction(ctx context.Context, transaction *models.BlockchainTransaction) error {
	defer func() {
		ts.processingTxs.Delete(transaction.ID)
	}()

	if transaction.Status != models.TransactionStatusPending {
		return nil
	}

	if transaction.TxHash.Valid {
		return nil
	}

	if transaction.NextRetryAt.Valid && transaction.NextRetryAt.Time.After(time.Now()) {
		return nil
	}


	// Send transaction based on type
	txHash, err := ts.sendRawTransaction(ctx, transaction)
	if err != nil {
		return err
	}

	// Update transaction status to sent
	if err := transaction.MarkSent(ctx, ts.db, txHash); err != nil {
		return err
	}

	log.Infof("Transaction %d sent successfully with hash: %s", transaction.ID, txHash)
	return nil
}

// sendRawTransaction sends a raw transaction to the blockchain
func (ts *TransactionSender) sendRawTransaction(ctx context.Context, transaction *models.BlockchainTransaction) (string, error) {
	client, err := GetBlockchainClient(transaction.Network)
	if err != nil {
		return "", err
	}

	if transaction.FromAddress != client.Address {
		return "", fmt.Errorf("from address is not the same as the client address")
	}

	auth, err := client.GetAuth(ctx)
	if err != nil {
		return "", err
	}

	client.NonceMu.Lock()
	defer client.NonceMu.Unlock()

	nonce, err := client.GetNonce(ctx)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = client.GasPrice
	auth.GasLimit = client.GasLimit

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	auth.Context = callCtx

	value, _ := new(big.Int).SetString(transaction.Value, 10)
	toAddress := common.HexToAddress(transaction.ToAddress)
	var data []byte
	if transaction.Data.Valid {
		data, err = hexutil.Decode(transaction.Data.String)
		if err != nil {
			return "", err
		}
	}

	baseTx := &types.LegacyTx{
		To:       &toAddress,
		Nonce:    nonce,
		GasPrice: client.GasPrice,
		Gas:      client.GasLimit,
		Value:    value,
		Data:     data,
	}

	rawTx := types.NewTx(baseTx)

	signedTx, err := auth.Signer(auth.From, rawTx)
	if err != nil {
		return "", err
	}

	err = client.RpcClient.SendTransaction(callCtx, signedTx)
	if err != nil {
		err = client.processSendingTxError(err)
		return "", err
	}

	client.IncrementNonce()

	return signedTx.Hash().Hex(), nil
}
