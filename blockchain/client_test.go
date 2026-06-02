package blockchain

import (
	"context"
	"crynux_relay_wallet/models"
	"database/sql"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
