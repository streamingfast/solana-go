package rpc

import (
	"encoding/json"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_GetBalance(t *testing.T) {
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
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"context":{"slot":131678970},"value":5465913541},"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {}
			},
			key:       solana.MustPublicKeyFromBase58("6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3"),
			expectOut: &GetBalanceResult{RPCContext: RPCContext{Context{Slot: 0x7d942fa}}, Value: bin.Uint64(5465913541)},
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
			key:       solana.MustPublicKeyFromBase58("6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3"),
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetBalance(test.key, nil)
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
