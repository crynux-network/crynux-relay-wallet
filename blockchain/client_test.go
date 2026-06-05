package blockchain

import (
	"context"
	"crynux_relay_wallet/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestGetTransactionTransferReadsRawCustomTypeFields(t *testing.T) {
	txHash := common.HexToHash("0x1234")
	from := common.HexToAddress("0x1111111111111111111111111111111111111111")
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
			ID     json.RawMessage   `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode rpc request: %v", err)
		}
		if request.Method != "eth_getTransactionByHash" {
			t.Fatalf("unexpected rpc method %s", request.Method)
		}

		response := map[string]any{
			"jsonrpc": "2.0",
			"id":      request.ID,
			"result": map[string]any{
				"type":  "0x7e",
				"hash":  txHash.Hex(),
				"from":  from.Hex(),
				"to":    to.Hex(),
				"value": "0x2a",
				"input": "0x",
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("encode rpc response: %v", err)
		}
	}))
	defer server.Close()

	client := &BlockchainClient{RpcEndpoint: server.URL}
	transfer, err := client.GetTransactionTransfer(context.Background(), txHash)
	if err != nil {
		t.Fatalf("GetTransactionTransfer failed: %v", err)
	}
	if transfer.Hash != txHash {
		t.Fatalf("expected hash %s, got %s", txHash.Hex(), transfer.Hash.Hex())
	}
	if transfer.From != from {
		t.Fatalf("expected from %s, got %s", from.Hex(), transfer.From.Hex())
	}
	if transfer.To == nil || *transfer.To != to {
		t.Fatalf("expected to %s, got %v", to.Hex(), transfer.To)
	}
	if transfer.Value.Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("expected value 42, got %s", transfer.Value.String())
	}
	if len(transfer.Input) != 0 {
		t.Fatalf("expected empty input, got %x", transfer.Input)
	}
}

func TestBuildCallMsgFromTransactionUsesDatabaseFields(t *testing.T) {
	from := common.HexToAddress("0x1111111111111111111111111111111111111111")
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")

	msg, err := BuildCallMsgFromTransaction(&models.BlockchainTransaction{
		FromAddress: from.Hex(),
		ToAddress:   to.Hex(),
		Value:       "42",
		Data:        sql.NullString{String: "0x1234", Valid: true},
	})
	if err != nil {
		t.Fatalf("BuildCallMsgFromTransaction failed: %v", err)
	}
	if msg.From != from {
		t.Fatalf("expected from %s, got %s", from.Hex(), msg.From.Hex())
	}
	if msg.To == nil || *msg.To != to {
		t.Fatalf("expected to %s, got %v", to.Hex(), msg.To)
	}
	if msg.Value.Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("expected value 42, got %s", msg.Value.String())
	}
	if string(msg.Data) != string([]byte{0x12, 0x34}) {
		t.Fatalf("expected data 0x1234, got %x", msg.Data)
	}
}

func TestParseERC20TransferAmount(t *testing.T) {
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")
	amount := big.NewInt(123)
	data, err := erc20ABI.Pack("transfer", to, amount)
	if err != nil {
		t.Fatalf("pack transfer failed: %v", err)
	}

	got, err := parseERC20TransferAmount(&models.BlockchainTransaction{
		Data: sql.NullString{String: hexutil.Encode(data), Valid: true},
	})
	if err != nil {
		t.Fatalf("parseERC20TransferAmount failed: %v", err)
	}
	if got.Cmp(amount) != 0 {
		t.Fatalf("expected amount %s, got %s", amount.String(), got.String())
	}
}

func TestIsHotWalletBalanceError(t *testing.T) {
	err := fmt.Errorf("wrapped error: %w", ErrHotWalletERC20BalanceInsufficient)
	if !isHotWalletBalanceError(err) {
		t.Fatalf("expected hot wallet balance error")
	}
	if isHotWalletBalanceError(fmt.Errorf("other error")) {
		t.Fatalf("expected non-balance error")
	}
}

func TestAddGasLimitBuffer(t *testing.T) {
	gasLimit, err := addGasLimitBuffer(21000, 20)
	if err != nil {
		t.Fatalf("addGasLimitBuffer failed: %v", err)
	}
	if gasLimit != 25200 {
		t.Fatalf("expected gas limit 25200, got %d", gasLimit)
	}
}

func TestBuildDynamicFeeTransactionUsesType2FeeCaps(t *testing.T) {
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")
	chainID := big.NewInt(8453)
	gasFeeCap := big.NewInt(100)
	gasTipCap := big.NewInt(2)

	tx, err := buildDynamicFeeTransaction(&models.BlockchainTransaction{
		ToAddress: to.Hex(),
		Value:     "42",
		Data:      sql.NullString{String: "0x1234", Valid: true},
	}, 7, 25200, gasFeeCap, gasTipCap, chainID)
	if err != nil {
		t.Fatalf("buildDynamicFeeTransaction failed: %v", err)
	}
	if tx.Type() != types.DynamicFeeTxType {
		t.Fatalf("expected dynamic fee transaction type, got %d", tx.Type())
	}
	if tx.ChainId().Cmp(chainID) != 0 {
		t.Fatalf("expected chain ID %s, got %s", chainID.String(), tx.ChainId().String())
	}
	if tx.Nonce() != 7 {
		t.Fatalf("expected nonce 7, got %d", tx.Nonce())
	}
	if tx.Gas() != 25200 {
		t.Fatalf("expected gas 25200, got %d", tx.Gas())
	}
	if tx.GasFeeCap().Cmp(gasFeeCap) != 0 {
		t.Fatalf("expected fee cap %s, got %s", gasFeeCap.String(), tx.GasFeeCap().String())
	}
	if tx.GasTipCap().Cmp(gasTipCap) != 0 {
		t.Fatalf("expected tip cap %s, got %s", gasTipCap.String(), tx.GasTipCap().String())
	}
	if tx.To() == nil || *tx.To() != to {
		t.Fatalf("expected to %s, got %v", to.Hex(), tx.To())
	}
	if tx.Value().Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("expected value 42, got %s", tx.Value().String())
	}
	if string(tx.Data()) != string([]byte{0x12, 0x34}) {
		t.Fatalf("expected data 0x1234, got %x", tx.Data())
	}
}

func TestValidateDynamicFeeCapsRejectsExceededCaps(t *testing.T) {
	err := validateDynamicFeeCaps(big.NewInt(101), big.NewInt(2), big.NewInt(100), big.NewInt(2))
	if err == nil {
		t.Fatalf("expected max fee cap error")
	}

	err = validateDynamicFeeCaps(big.NewInt(100), big.NewInt(3), big.NewInt(100), big.NewInt(2))
	if err == nil {
		t.Fatalf("expected max priority fee cap error")
	}

	err = validateDynamicFeeCaps(big.NewInt(101), big.NewInt(3), nil, nil)
	if err != nil {
		t.Fatalf("expected nil caps to allow dynamic fee values, got %v", err)
	}
}
