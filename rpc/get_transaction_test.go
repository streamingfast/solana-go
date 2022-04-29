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

func TestClient_GetTransaction(t *testing.T) {
	tests := []struct {
		name        string
		clientFunc  func(t *testing.T) (*Client, func(), func())
		signature   string
		expectError bool
		expectOut   interface{}
	}{
		{
			name: "mock json rpc request success transaction",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"blockTime":1650568375,"meta":{"err":null,"fee":10000,"innerInstructions":[{"index":2,"instructions":[{"accounts":[0,2],"data":"11119os1e9qSs2u7TsThXqkBSRVFxhmYaFKFZ1waB2X7armDmvK3p5GmLdUxYdg3h7QSrL","programIdIndex":9},{"accounts":[2,1],"data":"6dYJoXAVVvkRZc9RR2n57FF1dbE9CfReS65esF3PWLn5u","programIdIndex":10}]},{"index":4,"instructions":[{"accounts":[0,3],"data":"3Bxs4EMbRQoDyoj5","programIdIndex":9},{"accounts":[3],"data":"9krTDUMpjBo4wxLP","programIdIndex":9},{"accounts":[3],"data":"SYXsBkG6yKW2wWDcW8EDHR6D3P82bKxJGPpM65DD8nHqBfMP","programIdIndex":9},{"accounts":[0,4],"data":"3Bxs48v9NdVhakdd","programIdIndex":9},{"accounts":[4],"data":"9krTDgje7Fnho7ps","programIdIndex":9},{"accounts":[4],"data":"SYXsBkG6yKW2wWDcW8EDHR6D3P82bKxJGPpM65DD8nHqBfMP","programIdIndex":9},{"accounts":[1,0,0],"data":"biy1qNRB7L1cTkpmsoJhJYaVResvoqZqNJ8qQWcQ3rBgRoF","programIdIndex":10}]}],"logMessages":["Program 11111111111111111111111111111111 invoke [1]","Program 11111111111111111111111111111111 success","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]","Program log: Instruction: InitializeMint","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2343 of 200000 compute units","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success","Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL invoke [1]","Program log: Create","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Initialize the associated token account","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]","Program log: Instruction: InitializeAccount3","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2604 of 186319 compute units","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success","Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL consumed 16944 of 200000 compute units","Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL success","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]","Program log: Instruction: MintTo","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2515 of 200000 compute units","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success","Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s invoke [1]","Program log: Instruction: Mint New Edition from Master Edition Via Token","Program log: Transfer 5616720 lamports to the new account","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Allocate space for the account","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Assign the account to the owning program","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Transfer 2568240 lamports to the new account","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Allocate space for the account","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Assign the account to the owning program","Program 11111111111111111111111111111111 invoke [2]","Program 11111111111111111111111111111111 success","Program log: Setting mint authority","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]","Program log: Instruction: SetAuthority","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 1833 of 140267 compute units","Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success","Program log: Setting freeze authority","Program log: Skipping freeze authority because this mint has none","Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s consumed 62955 of 200000 compute units","Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s success"],"postBalances":[5343593681,1461600,2039280,5616720,2568240,2853600,1113600,1009200,31383018,1,953185920,2039280,5616720,853073280,1141440],"postTokenBalances":[{"accountIndex":2,"mint":"HHVpMURCU6gkLnW8tkwGzw6YJoWCj9RuLLESiTALY48x","owner":"HQc8axxhdu9jLfKtwcsmmGaF6LZPFgNjPAV1kThh3dew","uiTokenAmount":{"amount":"1","decimals":0,"uiAmount":1.0,"uiAmountString":"1"}},{"accountIndex":11,"mint":"F35m318ScNzFAb8iizXKpHque7n9d1p6pvyfAJ3ZRZzd","owner":"6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3","uiTokenAmount":{"amount":"1","decimals":0,"uiAmount":1.0,"uiAmountString":"1"}}],"preBalances":[5355289521,0,0,0,0,2853600,1113600,1009200,31383018,1,953185920,2039280,5616720,853073280,1141440],"preTokenBalances":[{"accountIndex":11,"mint":"F35m318ScNzFAb8iizXKpHque7n9d1p6pvyfAJ3ZRZzd","owner":"6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3","uiTokenAmount":{"amount":"1","decimals":0,"uiAmount":1.0,"uiAmountString":"1"}}],"rewards":[],"status":{"Ok":null}},"slot":130742131,"transaction":{"message":{"accountKeys":["6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3","HHVpMURCU6gkLnW8tkwGzw6YJoWCj9RuLLESiTALY48x","C8pR5bkLtKhUuXwh1oDz2R8mXS2etpYrgVKEMcgpFUwA","5a38uBMkkSB9bKWRobarhLZKiMCvCLT3HkZTn1bo3zDN","AeodNaL3t4bGmbjkCimjRHzbMDR7xWuFfFhUkzrVUY7b","3M4kjMG1mFzedjYV2ioSDN7g6PxqTfPErY4xBREYqnhF","526iyd7MvjXuEJsBr5ugxsgdkhXCgsaUJSV6KEGAjept","SysvarRent111111111111111111111111111111111","HQc8axxhdu9jLfKtwcsmmGaF6LZPFgNjPAV1kThh3dew","11111111111111111111111111111111","TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA","CVjPMikca93xzhm7JK76wRGapoib6fwVuyNvfbMgQywN","E61wZ3tYzA3KzwJQRghVuZyR9RsEbHL3ubsPHBED4mz6","ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL","metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":8,"numRequiredSignatures":2},"instructions":[{"accounts":[0,1],"data":"11114XtYk9gGfZoo968fyjNUYQJKf9gdmkGoaoBpzFv4vyaSMBn3VKxZdv7mZLzoyX5YNC","programIdIndex":9},{"accounts":[1,7],"data":"11TF6jTf6sYCva7BULe694SiHcRGsNY9kM8JPmMQKTbP11q","programIdIndex":10},{"accounts":[0,2,8,1,9,10,7],"data":"1","programIdIndex":13},{"accounts":[1,2,0],"data":"6AuM4xMCPFhR","programIdIndex":10},{"accounts":[3,4,5,1,6,0,0,0,11,0,12,10,9,7],"data":"9EsUmmETCEMm","programIdIndex":14}],"recentBlockhash":"5rMbS2L6ZhhJHJS9m8jscYuwbeLsKgbTj7USY7cpKKGo"},"signatures":["29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs","63s6CPYA4VnAhNkuF2U88vUqcWeEWBY5QJJRoUmckrSbnL37BJy1GaJQxPgMQH8apQmVbVbmcAGL5BiCasWegbtp"]}},"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getTransaction", "params": []interface{}{"29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs", "json"}}, server.RequestBody(t))
				}
			},
			signature: "29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs",
			expectOut: &GetTransactionResponse{
				Slot:      0x7caf773,
				BlockTime: puint64(1650568375),
				Transaction: &Transaction{
					Signatures: []string{
						"29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs",
						"63s6CPYA4VnAhNkuF2U88vUqcWeEWBY5QJJRoUmckrSbnL37BJy1GaJQxPgMQH8apQmVbVbmcAGL5BiCasWegbtp",
					},
					Message: &Message{
						AccountKeys: []solana.PublicKey{
							solana.MustPublicKeyFromBase58("6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3"),
							solana.MustPublicKeyFromBase58("HHVpMURCU6gkLnW8tkwGzw6YJoWCj9RuLLESiTALY48x"),
							solana.MustPublicKeyFromBase58("C8pR5bkLtKhUuXwh1oDz2R8mXS2etpYrgVKEMcgpFUwA"),
							solana.MustPublicKeyFromBase58("5a38uBMkkSB9bKWRobarhLZKiMCvCLT3HkZTn1bo3zDN"),
							solana.MustPublicKeyFromBase58("AeodNaL3t4bGmbjkCimjRHzbMDR7xWuFfFhUkzrVUY7b"),
							solana.MustPublicKeyFromBase58("3M4kjMG1mFzedjYV2ioSDN7g6PxqTfPErY4xBREYqnhF"),
							solana.MustPublicKeyFromBase58("526iyd7MvjXuEJsBr5ugxsgdkhXCgsaUJSV6KEGAjept"),
							solana.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111"),
							solana.MustPublicKeyFromBase58("HQc8axxhdu9jLfKtwcsmmGaF6LZPFgNjPAV1kThh3dew"),
							solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
							solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"),
							solana.MustPublicKeyFromBase58("CVjPMikca93xzhm7JK76wRGapoib6fwVuyNvfbMgQywN"),
							solana.MustPublicKeyFromBase58("E61wZ3tYzA3KzwJQRghVuZyR9RsEbHL3ubsPHBED4mz6"),
							solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"),
							solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s"),
						},
						Header: MessageHeader{
							NumReadonlySignedAccounts:   0,
							NumReadonlyUnsignedAccounts: 8,
							NumRequiredSignatures:       2,
						},
						Instructions: []Instruction{
							{
								ProgramIdIndex: 9,
								Data:           "11114XtYk9gGfZoo968fyjNUYQJKf9gdmkGoaoBpzFv4vyaSMBn3VKxZdv7mZLzoyX5YNC",
								Accounts: []bin.Uint64{
									0, 1,
								},
							},
							{
								ProgramIdIndex: 10,
								Data:           "11TF6jTf6sYCva7BULe694SiHcRGsNY9kM8JPmMQKTbP11q",
								Accounts: []bin.Uint64{
									1, 7,
								},
							},
							{
								ProgramIdIndex: 13,
								Data:           "1",
								Accounts: []bin.Uint64{
									0, 2, 8, 1, 9, 10, 7,
								},
							},
							{
								ProgramIdIndex: 10,
								Data:           "6AuM4xMCPFhR",
								Accounts: []bin.Uint64{
									1, 2, 0,
								},
							},
							{
								ProgramIdIndex: 14,
								Data:           "9EsUmmETCEMm",
								Accounts: []bin.Uint64{
									3, 4, 5, 1, 6, 0, 0, 0, 11, 0, 12, 10, 9, 7,
								},
							},
						},
						RecentBlockhash: solana.MustPublicKeyFromBase58("5rMbS2L6ZhhJHJS9m8jscYuwbeLsKgbTj7USY7cpKKGo"),
					},
				},
				Meta: &Meta{
					Fee:          10000,
					PreBalances:  []bin.Uint64{5355289521, 0, 0, 0, 0, 2853600, 1113600, 1009200, 31383018, 1, 953185920, 2039280, 5616720, 853073280, 1141440},
					PostBalances: []bin.Uint64{5343593681, 1461600, 2039280, 5616720, 2568240, 2853600, 1113600, 1009200, 31383018, 1, 953185920, 2039280, 5616720, 853073280, 1141440},
					InnerInstructions: []*InnerInstruction{
						{
							Index: 2,
							Instructions: []InstructionMeta{
								{
									Accounts:       []bin.Uint64{0, 2},
									Data:           "11119os1e9qSs2u7TsThXqkBSRVFxhmYaFKFZ1waB2X7armDmvK3p5GmLdUxYdg3h7QSrL",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{2, 1},
									Data:           "6dYJoXAVVvkRZc9RR2n57FF1dbE9CfReS65esF3PWLn5u",
									ProgramIdIndex: 10,
								},
							},
						},
						{
							Index: 4,
							Instructions: []InstructionMeta{
								{
									Accounts:       []bin.Uint64{0, 3},
									Data:           "3Bxs4EMbRQoDyoj5",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{3},
									Data:           "9krTDUMpjBo4wxLP",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{3},
									Data:           "SYXsBkG6yKW2wWDcW8EDHR6D3P82bKxJGPpM65DD8nHqBfMP",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{0, 4},
									Data:           "3Bxs48v9NdVhakdd",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{4},
									Data:           "9krTDgje7Fnho7ps",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{4},
									Data:           "SYXsBkG6yKW2wWDcW8EDHR6D3P82bKxJGPpM65DD8nHqBfMP",
									ProgramIdIndex: 9,
								},
								{
									Accounts:       []bin.Uint64{1, 0, 0},
									Data:           "biy1qNRB7L1cTkpmsoJhJYaVResvoqZqNJ8qQWcQ3rBgRoF",
									ProgramIdIndex: 10,
								},
							},
						},
					},
					PostTokenBalances: []*TokeBalance{
						{
							AccountIndex: 2,
							Mint:         solana.MustPublicKeyFromBase58("HHVpMURCU6gkLnW8tkwGzw6YJoWCj9RuLLESiTALY48x"),
							Owner:        solana.MustPublicKeyFromBase58("HQc8axxhdu9jLfKtwcsmmGaF6LZPFgNjPAV1kThh3dew"),
						},
						{
							AccountIndex: 11,
							Mint:         solana.MustPublicKeyFromBase58("F35m318ScNzFAb8iizXKpHque7n9d1p6pvyfAJ3ZRZzd"),
							Owner:        solana.MustPublicKeyFromBase58("6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3"),
						},
					},
					PreTokenBalances: []*TokeBalance{
						{
							AccountIndex: 11,
							Mint:         solana.MustPublicKeyFromBase58("F35m318ScNzFAb8iizXKpHque7n9d1p6pvyfAJ3ZRZzd"),
							Owner:        solana.MustPublicKeyFromBase58("6wrL8rQzDWSH7PJyZRGsdBiNcrpD8Wd6vJzGVBinuCL3"),
						},
					},
					LogMessages: []string{
						"Program 11111111111111111111111111111111 invoke [1]",
						"Program 11111111111111111111111111111111 success",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]",
						"Program log: Instruction: InitializeMint",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2343 of 200000 compute units",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
						"Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL invoke [1]",
						"Program log: Create",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Initialize the associated token account",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]",
						"Program log: Instruction: InitializeAccount3",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2604 of 186319 compute units",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
						"Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL consumed 16944 of 200000 compute units",
						"Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL success",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]",
						"Program log: Instruction: MintTo",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 2515 of 200000 compute units",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
						"Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s invoke [1]",
						"Program log: Instruction: Mint New Edition from Master Edition Via Token",
						"Program log: Transfer 5616720 lamports to the new account",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Allocate space for the account",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Assign the account to the owning program",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Transfer 2568240 lamports to the new account",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Allocate space for the account",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Assign the account to the owning program",
						"Program 11111111111111111111111111111111 invoke [2]",
						"Program 11111111111111111111111111111111 success",
						"Program log: Setting mint authority",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]",
						"Program log: Instruction: SetAuthority",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 1833 of 140267 compute units",
						"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
						"Program log: Setting freeze authority",
						"Program log: Skipping freeze authority because this mint has none",
						"Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s consumed 62955 of 200000 compute units",
						"Program metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s success",
					},
					Rewards: []interface{}{},
				},
			},
		},
		{
			name: "mock json rpc request success transaction",
			clientFunc: func(t *testing.T) (*Client, func(), func()) {
				server, closer := mockJSONRPC(t, json.RawMessage(`{"jsonrpc":"2.0","result":{"blockTime":1651153595,"meta":{"err":{"InstructionError":[0,{"Custom":41}]},"fee":5000,"innerInstructions":[{"index":0,"instructions":[{"accounts":[5,6,7,3,4,2,8,0,9],"data":"19krTD1r94kgZ9vMd","programIdIndex":10}]}],"logMessages":["Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk invoke [1]","Program log: Instruction: CancelPerpOrder","Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK invoke [2]","Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK consumed 15194 of 173545 compute units","Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK failed: custom program error: 0x29","Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk consumed 41649 of 200000 compute units","Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk failed: custom program error: 0x29"],"postBalances":[11012972500841,83944560,6410160,32141280,23636160,10440000,457104960,457104960,1825496640,253692000,1141440,1141440],"postTokenBalances":[],"preBalances":[11012972505841,83944560,6410160,32141280,23636160,10440000,457104960,457104960,1825496640,253692000,1141440,1141440],"preTokenBalances":[],"rewards":[],"status":{"Err":{"InstructionError":[0,{"Custom":41}]}}},"slot":131683246,"transaction":{"message":{"accountKeys":["GjirmTYq3N2MepKgGiafFN6u294CT79cuek7p8fY3dY2","6kSuaA8Fi16ks53ye8eqgr7jZn2uus59W7gMSasgtFon","2GwQEz5p9sLJShcUf7TPWJFZ1KckWSCD7zT3H4kvV8Kp","3z5HfN7PtvCNLwcNrwWWPrD4JpByNJxfwKoWv1rsV6ro","3QuwRpf9rk1W3kRXf3tLRbdjhbA6zTeyPfVAy379tFN5","9sRguTCVrmoTFsFJfDtpd9AhaB4ZAC5qQx8ddUhQ1vJM","DCYEq1BhmEdtjP3AK537vEwjezv9JUEVETCPMSA2pem1","GFU7oh9eQEwkDbifzB7cyB4GgQL9i1Aw8vihaXXjhQmS","CF61UE6CbgpvJDZDzPmxpS52EykvKZ1BKMfgn9qf7e8N","71yykwxq1zQqy99PgRsgZJXi2HHK2UDx9G4va7pH6qRv","ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK","Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk"],"header":{"numReadonlySignedAccounts":0,"numReadonlyUnsignedAccounts":3,"numRequiredSignatures":1},"instructions":[{"accounts":[9,1,0,2,3,4,5,6,7,8,10],"data":"YYdLQNbX79Q79th9UYsxcVCCQf","programIdIndex":11}],"recentBlockhash":"DCQ49uwWuUzniFpScSD4YhVUAq2wUtUeMUm51ji6vdNG"},"signatures":["1XGtQ2XJe9gs5Dysp3WZa5MYGRPNpMQXvbDpYXTNaB4HttQAoQ1mhLYmRyNLq9kY8bCkPCzei4DjbdE8QoKSG4q"]}},"id":1}`))
				client := newTestClient(server.URL)
				return client, closer, func() {
					assert.Equal(t, map[string]interface{}{"id": float64(0), "jsonrpc": "2.0", "method": "getTransaction", "params": []interface{}{"29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs", "json"}}, server.RequestBody(t))
				}
			},
			signature: "29RTUcaTCA48QBsQmYxjBofUogEyk8q6cJiThCWqYkCPNEvSGZgYyrwtRXJKDW4VWMLN8qeLjyH28cwXQzsKwExs",
			expectOut: &GetTransactionResponse{
				Slot:      131683246,
				BlockTime: puint64(1651153595),
				Transaction: &Transaction{
					Signatures: []string{
						"1XGtQ2XJe9gs5Dysp3WZa5MYGRPNpMQXvbDpYXTNaB4HttQAoQ1mhLYmRyNLq9kY8bCkPCzei4DjbdE8QoKSG4q",
					},
					Message: &Message{
						AccountKeys: []solana.PublicKey{
							solana.MustPublicKeyFromBase58("GjirmTYq3N2MepKgGiafFN6u294CT79cuek7p8fY3dY2"),
							solana.MustPublicKeyFromBase58("6kSuaA8Fi16ks53ye8eqgr7jZn2uus59W7gMSasgtFon"),
							solana.MustPublicKeyFromBase58("2GwQEz5p9sLJShcUf7TPWJFZ1KckWSCD7zT3H4kvV8Kp"),
							solana.MustPublicKeyFromBase58("3z5HfN7PtvCNLwcNrwWWPrD4JpByNJxfwKoWv1rsV6ro"),
							solana.MustPublicKeyFromBase58("3QuwRpf9rk1W3kRXf3tLRbdjhbA6zTeyPfVAy379tFN5"),
							solana.MustPublicKeyFromBase58("9sRguTCVrmoTFsFJfDtpd9AhaB4ZAC5qQx8ddUhQ1vJM"),
							solana.MustPublicKeyFromBase58("DCYEq1BhmEdtjP3AK537vEwjezv9JUEVETCPMSA2pem1"),
							solana.MustPublicKeyFromBase58("GFU7oh9eQEwkDbifzB7cyB4GgQL9i1Aw8vihaXXjhQmS"),
							solana.MustPublicKeyFromBase58("CF61UE6CbgpvJDZDzPmxpS52EykvKZ1BKMfgn9qf7e8N"),
							solana.MustPublicKeyFromBase58("71yykwxq1zQqy99PgRsgZJXi2HHK2UDx9G4va7pH6qRv"),
							solana.MustPublicKeyFromBase58("ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK"),
							solana.MustPublicKeyFromBase58("Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk"),
						},
						Header: MessageHeader{
							NumReadonlySignedAccounts:   0,
							NumReadonlyUnsignedAccounts: 3,
							NumRequiredSignatures:       1,
						},
						Instructions: []Instruction{
							{
								ProgramIdIndex: 11,
								Data:           "YYdLQNbX79Q79th9UYsxcVCCQf",
								Accounts:       []bin.Uint64{9, 1, 0, 2, 3, 4, 5, 6, 7, 8, 10},
							},
						},
						RecentBlockhash: solana.MustPublicKeyFromBase58("DCQ49uwWuUzniFpScSD4YhVUAq2wUtUeMUm51ji6vdNG"),
					},
				},
				Meta: &Meta{
					Err: &TransactionError{
						Raw: map[string]interface{}{
							"InstructionError": []interface{}{
								float64(0),
								map[string]interface{}{
									"Custom": float64(41),
								},
							},
						},
						InstructionIndex:     0,
						InstructionErrorCode: "unknown",
						InstructionErrorType: "Custom",
					},
					Fee: 5000,
					PreBalances: []bin.Uint64{11012972505841,
						83944560,
						6410160,
						32141280,
						23636160,
						10440000,
						457104960,
						457104960,
						1825496640,
						253692000,
						1141440,
						1141440,
					},
					PostBalances: []bin.Uint64{11012972500841,
						83944560,
						6410160,
						32141280,
						23636160,
						10440000,
						457104960,
						457104960,
						1825496640,
						253692000,
						1141440,
						1141440,
					},
					InnerInstructions: []*InnerInstruction{
						{
							Index: 0,
							Instructions: []InstructionMeta{
								{
									Accounts:       []bin.Uint64{5, 6, 7, 3, 4, 2, 8, 0, 9},
									Data:           "19krTD1r94kgZ9vMd",
									ProgramIdIndex: 10,
								},
							},
						},
					},
					PostTokenBalances: []*TokeBalance{},
					PreTokenBalances:  []*TokeBalance{},
					LogMessages: []string{
						"Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk invoke [1]",
						"Program log: Instruction: CancelPerpOrder",
						"Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK invoke [2]",
						"Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK consumed 15194 of 173545 compute units",
						"Program ZDx8a8jBqGmJyxi1whFxxCo5vG6Q9t4hTzW2GSixMKK failed: custom program error: 0x29",
						"Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk consumed 41649 of 200000 compute units",
						"Program Zo1ggzTUKMY5bYnDvT5mtVeZxzf2FaLTbKkmvGUhUQk failed: custom program error: 0x29",
					},
					Rewards: []interface{}{},
				},
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
			signature: "1XGtQ2XJe9gs5Dysp3WZa5MYGRPNpMQXvbDpYXTNaB4HttQAoQ1mhLYmRyNLq9kY8bCkPCzei4DjbdE8QoKSG4q",
			expectOut: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, cleanup, assertions := test.clientFunc(t)
			defer cleanup()
			out, err := client.GetTransaction(test.signature)
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
