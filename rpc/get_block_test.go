package rpc

import (
	"encoding/json"
	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_GetBlock(t *testing.T) {
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		slotNum     uint64
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"blockHeight":80300230,"blockTime":1627988608,"blockhash":"DUCT8VSgk2BXkMhQfxKVYfikEZCQf4dZ4ioPdGdaVxMN","parentSlot":429,"previousBlockhash":"HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs","transactions":[{"meta":null,"transaction":{"message":{"accountKeys":["GdnSyH3YtwcxFvQrVVJMm1JhTS4QVX7MFsX56uJLUfiZ","sCtiJieP8B3SwYnXemiLpRFRR8KJLMtsMVN25fAFWjW","SysvarS1otHashes111111111111111111111111111","SysvarC1ock11111111111111111111111111111111","Vote111111111111111111111111111111111111111"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":3,"numRequiredSignatures":1},"instructions":[{"accounts":[1,2,3,0],"data":"37u9WtQpcm6ULa3WRQHmj49EPs5hZfet7nKmHTYD6Vww48rtUK99VNvA7fVscfBpL6KyjwiF","programIdIndex":4}],"recentBlockhash":"HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs"},"signatures":["MsdZAVaCjHcVWs8zMJinXvntufdXwtHJWCRLSyw9zeAZuNDec6s41H12KFFyPHbq3uj98wRjMa86z6nW2kUv1Zs"]}},{"meta":null,"transaction":{"message":{"accountKeys":["CakcnaRDHka2gXyfbEd2d3xsvkJkqsLw2akB3zsN1D2S","9bRDrYShoQ77MZKYTMoAsoCkU7dAR24mxYCBjXLpfEJx","SysvarS1otHashes111111111111111111111111111","SysvarC1ock11111111111111111111111111111111","Vote111111111111111111111111111111111111111"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":3,"numRequiredSignatures":1},"instructions":[{"accounts":[1,2,3,0],"data":"37u9WtQpcm6ULa3WRQHmj49EPs5hZfet7nKmHTYD6Vww48rtUK99VNvA7fVscfBpL6KyjwiF","programIdIndex":4}],"recentBlockhash":"HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs"},"signatures":["RetEYHPRuKmR5qyXqWxSusma3PEqeYdteSahAkasiMicsq312toEddsc16jQdPr1NFGkPvnQ9MkmossLTsEmeXM"]}},{"meta":null,"transaction":{"message":{"accountKeys":["7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2","4785anyR2rYSas6cQGHtykgzwYEtChvFYhcEgdDw3gGL","SysvarS1otHashes111111111111111111111111111","SysvarC1ock11111111111111111111111111111111","Vote111111111111111111111111111111111111111"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":3,"numRequiredSignatures":1},"instructions":[{"accounts":[1,2,3,0],"data":"37u9WtQpcm6ULa3WRQHmj49EPs5hZfet7nKmHTYD6Vww48rtUK99VNvA7fVscfBpL6KyjwiF","programIdIndex":4}],"recentBlockhash":"HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs"},"signatures":["5iwUiSfvaiD2JEeHZKg72wMBcUatg8tW7qSh3ryHqTYuAprpskVZLbfKqPrQ7VMjEVj9oLAGW84JHNKJ99TNBpdV"]}},{"meta":null,"transaction":{"message":{"accountKeys":["DE1bawNcRJB9rVm3buyMVfr8mBEoyyu73NBovf2oXJsJ","8XgHUtBRY6qePVYERxosyX3MUq8NQkjtmFDSzQ2WpHTJ","SysvarS1otHashes111111111111111111111111111","SysvarC1ock11111111111111111111111111111111","Vote111111111111111111111111111111111111111"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":3,"numRequiredSignatures":1},"instructions":[{"accounts":[1,2,3,0],"data":"37u9WtQpcm6ULa3WRQHmj49EPs5hZfet7nKmHTYD6Vww48rtUK99VNvA7fVscfBpL6KyjwiF","programIdIndex":4}],"recentBlockhash":"HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs"},"signatures":["5Vf9ppLf9saAwEoAxfo3TJfuUbE7LyC7nBNgM9iovjSTt8SSMYsJjxhvvmSRsMjans5SqGXxuT1xbnm8MuK6kcSp"]}}]},"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getBlock", "params": []interface{}{float64(430)}}, server.RequestBody(t))
				}
			},
			slotNum: 430,
			expectOut: &GetBlockResult{
				BlockHeight:       puint64(80300230),
				BlockTime:         puint64(1627988608),
				Blockhash:         solana.MustPublicKeyFromBase58("DUCT8VSgk2BXkMhQfxKVYfikEZCQf4dZ4ioPdGdaVxMN"),
				ParentSlot:        429,
				PreviousBlockhash: solana.MustPublicKeyFromBase58("HA2fJgGqmQezCXJRVNZAWPbRMXCPjUyo7VjRF47JGdYs"),
				//Transactions:      nil,
			},
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
			slotNum:   430,
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetBlock(test.slotNum)
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
