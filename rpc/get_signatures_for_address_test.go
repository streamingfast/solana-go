package rpc

import (
	"encoding/json"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_GetSignaturesForAddress(t *testing.T) {
	limit := uint64(2)
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		address     solana.PublicKey
		opts        *GetSignaturesForAddressOpts
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":[{"blockTime":1649888395,"confirmationStatus":"finalized","err":null,"memo":null,"signature":"sbmnRxSb3uW1nmGD4mXrj6MJKtfErcxPV2AGKy9GzziDB37M6W9p87qxWobHTZmvRfJ2dYomWzLE6XCWYWASLa6","slot":129612604},{"blockTime":1649888394,"confirmationStatus":"finalized","err":null,"memo":null,"signature":"3andPyixBAjKRaM7rQwpsXZf68ug4Lr19EhedN69F6KX9218Us9AydC3Z5VCA452X3jUhVuh1fK9VQhr8D4RSv8p","slot":129612602}],"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getSignaturesForAddress", "params": []interface{}{
						"5PzHeoZPEbsW8GuQZ8ct9feFhSPdwBJY4Qb8CBaDLzN7",
						map[string]interface{}{
							"limit": float64(2),
						},
					}}, server.RequestBody(t))
				}
			},
			address: solana.MustPublicKeyFromBase58("5PzHeoZPEbsW8GuQZ8ct9feFhSPdwBJY4Qb8CBaDLzN7"),
			opts: &GetSignaturesForAddressOpts{
				Limit: &limit,
			},
			expectOut: GetSignaturesForAddressResult([]*TransactionSignature{
				{
					BlockTime:          1649888395,
					ConfirmationStatus: "finalized",
					Err:                nil,
					Memo:               nil,
					Signature:          "sbmnRxSb3uW1nmGD4mXrj6MJKtfErcxPV2AGKy9GzziDB37M6W9p87qxWobHTZmvRfJ2dYomWzLE6XCWYWASLa6",
					Slot:               129612604,
				},
				{
					BlockTime:          1649888394,
					ConfirmationStatus: "finalized",
					Err:                nil,
					Memo:               nil,
					Signature:          "3andPyixBAjKRaM7rQwpsXZf68ug4Lr19EhedN69F6KX9218Us9AydC3Z5VCA452X3jUhVuh1fK9VQhr8D4RSv8p",
					Slot:               129612602,
				},
			}),
		},
		{
			name: "real json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				rpcUrl := os.Getenv("TEST_RPC_URL")
				if rpcUrl == "" {
					t.Skip("skipping test TEST_RPC_URL not defined")
				}
				return NewClient(rpcUrl), func() {}, func() {}
			},
			address: solana.MustPublicKeyFromBase58("5PzHeoZPEbsW8GuQZ8ct9feFhSPdwBJY4Qb8CBaDLzN7"),
			opts: &GetSignaturesForAddressOpts{
				Limit: &limit,
			},
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetSignaturesForAddress(test.address, test.opts)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if !isNil(test.expectOut) {
					assert.Equal(t, test.expectOut, out)
				}
				assertions()
			}
		})
	}
}
