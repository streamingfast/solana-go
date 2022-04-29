package rpc

import (
	"encoding/json"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_GetLatestBlockHash(t *testing.T) {
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		key         solana.PublicKey
		commitment  CommitmentType
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"context":{"slot":131727818},"value":{"blockhash":"F3kFjvpvUig5C3yyudmaGMosoB2UCo6aKUikLV3o6LS8","lastValidBlockHeight":119478575}},"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getLatestBlockhash", "params": []interface{}{
						map[string]interface{}{"commitment": "processed"},
					}}, server.RequestBody(t))
				}
			},
			commitment: CommitmentProcessed,
			expectOut: &GetRecentBlockhashResult{RPCContext: RPCContext{Context{Slot: 131727818}}, Value: &BlockhashResult{
				Blockhash:            solana.MustPublicKeyFromBase58("F3kFjvpvUig5C3yyudmaGMosoB2UCo6aKUikLV3o6LS8"),
				LastValidBlockHeight: 119478575,
			}},
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
			commitment: CommitmentProcessed,
			expectOut:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetLatestBlockhash(test.commitment)
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
