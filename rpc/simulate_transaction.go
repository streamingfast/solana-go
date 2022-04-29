package rpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type SimulateTransactionResponse struct {
	Err  *TransactionError
	Logs []string
}

func (c *Client) SimulateTransaction(transaction *solana.Transaction) (*SimulateTransactionResponse, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewEncoder(buf).Encode(transaction); err != nil {
		return nil, fmt.Errorf("send transaction: encode transaction: %w", err)
	}
	trxData := buf.Bytes()

	obj := map[string]interface{}{
		"encoding": "base64",
	}

	b64Data := base64.StdEncoding.EncodeToString(trxData)
	params := []interface{}{
		b64Data,
		obj,
	}

	var out *SimulateTransactionResponse
	if err := c.DoRequest(&out, "simulateTransaction", params...); err != nil {
		return nil, fmt.Errorf("send transaction: rpc send: %w", err)
	}

	return out, nil

}
