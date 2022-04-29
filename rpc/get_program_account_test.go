package rpc

import (
	"encoding/base64"
	"encoding/json"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func mustB64Decode(in string) []byte {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return out

}

func TestClient_GetProgramAccount(t *testing.T) {
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		address     solana.PublicKey
		opts        *GetProgramAccountsOpts
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":[{"account":{"data":["AhOU6/6wEDbzB7bCI8xU0VWxKmtDSsWS/TDN2++3BcsNhkqF4Ap2IaxWTbaO4JD487e8X7fW1AaEGH4O1jC8WvolWTx0H1SZjrxVWgWXegZXZeUa2lesJqydsCntsH+NCanbpTq0K5kc+A1JWBrOz4gZ4PTGcGEdT+8xtBtWwTsaenqKBgAAAABaDo8GAAAAAEpYkQYAAAAAAV50vXRaTZ4SCvMW5ILY+q6ny5HYyplwjssIaHWTqCBAAddA1EogKYgLs5xWnA6UviRJaC9d7wYdTp63QPHEKDaUAd/SGZF98DTiFIRD7cq1/OEo5vl7BSykHeq3Cc4TzoR7gPD6AgAAAAABASAAAABhbGllbnVzIGFyZ2VudHVtIGNhZWN1cwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==","base64"],"executable":false,"lamports":3173760,"owner":"HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs","rentEpoch":304},"pubkey":"J5DMRcNQaR1AKVSVrNLgYT2d7Y8f391YFBqMge9eAygc"},{"account":{"data":["A8z1AfFt0iBFT+dWq87Dq7G8PKpz7XqWnOLQrxNuaYpBw3Oh4FowHmZZUpb7HR0xF0EKB8s9dROqhKwOjoohh7pWUc2ShzcqOY2aotLuahKLA8thkH3ha1Lq+RRPFuPFvA==","base64"],"executable":false,"lamports":1566000,"owner":"HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs","rentEpoch":304},"pubkey":"BEBB9n89kYuXRCQR6FHgBbM8xWr4eaPeUvE81WBPc8XR"},{"account":{"data":["AmBUoEyNaxnFVSMkqZ6d9VciuQdBc9D/kMgDjMaIZXqsmVrdKtfSPeoMQVw4zPBEo6Eaj7O8dtg++7j4mzx46SwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHuls0Jd1F7WPcCJQ3WSK0fM3MK03rcODXoLIen/LtxjenqKBgAAAABaDo8GAAAAAEpYkQYAAAAAAAAAAAAAAAAAAAABACAAAABkaXZpbnVzIGJsbGljdXMgdnVsZ2l2YWd1cwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==","base64"],"executable":false,"lamports":3173760,"owner":"HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs","rentEpoch":305},"pubkey":"2XPHKoJBH33uwf9RNkLn93ji1MML2v14PSrLhMoLPEoq"}]}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getProgramAccounts", "params": []interface{}{
						"HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs",
						map[string]interface{}{
							"encoding": "base64",
						},
					}}, server.RequestBody(t))
				}
			},
			address: solana.MustPublicKeyFromBase58("HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs"),
			expectOut: GetProgramAccountsResult([]*KeyedAccount{
				{
					Pubkey: solana.MustPublicKeyFromBase58("J5DMRcNQaR1AKVSVrNLgYT2d7Y8f391YFBqMge9eAygc"),
					Account: &Account{
						Lamports:   3173760,
						Data:       mustB64Decode("AhOU6/6wEDbzB7bCI8xU0VWxKmtDSsWS/TDN2++3BcsNhkqF4Ap2IaxWTbaO4JD487e8X7fW1AaEGH4O1jC8WvolWTx0H1SZjrxVWgWXegZXZeUa2lesJqydsCntsH+NCanbpTq0K5kc+A1JWBrOz4gZ4PTGcGEdT+8xtBtWwTsaenqKBgAAAABaDo8GAAAAAEpYkQYAAAAAAV50vXRaTZ4SCvMW5ILY+q6ny5HYyplwjssIaHWTqCBAAddA1EogKYgLs5xWnA6UviRJaC9d7wYdTp63QPHEKDaUAd/SGZF98DTiFIRD7cq1/OEo5vl7BSykHeq3Cc4TzoR7gPD6AgAAAAABASAAAABhbGllbnVzIGFyZ2VudHVtIGNhZWN1cwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="),
						Owner:      solana.MustPublicKeyFromBase58("HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs"),
						Executable: false,
						RentEpoch:  304,
					},
				},
				{
					Pubkey: solana.MustPublicKeyFromBase58("BEBB9n89kYuXRCQR6FHgBbM8xWr4eaPeUvE81WBPc8XR"),
					Account: &Account{
						Lamports:   1566000,
						Data:       mustB64Decode("A8z1AfFt0iBFT+dWq87Dq7G8PKpz7XqWnOLQrxNuaYpBw3Oh4FowHmZZUpb7HR0xF0EKB8s9dROqhKwOjoohh7pWUc2ShzcqOY2aotLuahKLA8thkH3ha1Lq+RRPFuPFvA=="),
						Owner:      solana.MustPublicKeyFromBase58("HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs"),
						Executable: false,
						RentEpoch:  304,
					},
				},
				{
					Pubkey: solana.MustPublicKeyFromBase58("2XPHKoJBH33uwf9RNkLn93ji1MML2v14PSrLhMoLPEoq"),
					Account: &Account{
						Lamports:   3173760,
						Data:       mustB64Decode("AmBUoEyNaxnFVSMkqZ6d9VciuQdBc9D/kMgDjMaIZXqsmVrdKtfSPeoMQVw4zPBEo6Eaj7O8dtg++7j4mzx46SwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHuls0Jd1F7WPcCJQ3WSK0fM3MK03rcODXoLIen/LtxjenqKBgAAAABaDo8GAAAAAEpYkQYAAAAAAAAAAAAAAAAAAAABACAAAABkaXZpbnVzIGJsbGljdXMgdnVsZ2l2YWd1cwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="),
						Owner:      solana.MustPublicKeyFromBase58("HRBRbMQF38Z2hQSUCDFhMKQo6FK3qWQ3d1p3fa5TKKSs"),
						Executable: false,
						RentEpoch:  305,
					},
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
			address:   solana.MustPublicKeyFromBase58("5PzHeoZPEbsW8GuQZ8ct9feFhSPdwBJY4Qb8CBaDLzN7"),
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetProgramAccounts(test.address, test.opts)
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
