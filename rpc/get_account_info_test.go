package rpc

import (
	"encoding/json"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_GetAccountInfo(t *testing.T) {
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		key         solana.PublicKey
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"context":{"slot":1},"value":{"data":["dGVzdA==","base64"]}},"id":0}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getAccountInfo", "params": []interface{}{
						"7xLk17EQQ5KLDLDe44wCmupJKJjTGd8hs3eSVVhCx932",
						map[string]interface{}{"encoding": "base64"},
					}}, server.RequestBody(t))
				}
			},
			key:       solana.MustPublicKeyFromBase58("7xLk17EQQ5KLDLDe44wCmupJKJjTGd8hs3eSVVhCx932"),
			expectOut: &GetAccountInfoResult{RPCContext: RPCContext{Context{Slot: 1}}, Value: &Account{Data: []byte{0x74, 0x65, 0x73, 0x74}}},
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
			key:       solana.MustPublicKeyFromBase58("AeodNaL3t4bGmbjkCimjRHzbMDR7xWuFfFhUkzrVUY7b"),
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetAccountInfo(test.key)
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
