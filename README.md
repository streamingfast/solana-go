#  StreamingFast Solana library for Go
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/solana-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go library to interface with Solana nodes's JSON-RPC interface, Solana's SPL tokens various instructions decoding for popular programs.

> :warning: `solana-go` works using SemVer but in 0 version, which means that the 'minor' will be changed when some broken changes are introduced into the application, and the 'patch' will be changed when a new feature with new changes is added or for bug fixing. As soon as v1.0.0 be released, `solana-go` will start to use SemVer as usual.

## Installation

```
go get github.com/streamingfast/solana-go
```

> All development happens on the `develop` branch so if you need to check if there is any fixes in there, try updating to it with `go get github.com/streamingfast/solana-go@develop`

## Usage

Loading an SPL mint

```golang

import "github.com/streamingfast/solana-go/rpc"
import "github.com/streamingfast/solana-go/token"

	addr := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	cli := rpc.NewClient("https://api.mainnet-beta.solana.com")

	var m token.Mint
	err := cli.GetAccountDataIn(context.Background(), addr, &m)
	// handle `err`

	json.NewEncoder(os.Stdout).Encode(m)
	// {"OwnerOption":1,
	//  "Owner":"2wmVCSfPxGPjrnMMn7rchp4uaeoTqN39mXFC2zhPdri9",
	//  "Decimals":128,
	//  "IsInitialized":true}

```


Getting any account's data:

```golang

import "github.com/streamingfast/solana-go/rpc"
import "github.com/streamingfast/solana-go/token"

	addr := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	cli := rpc.NewClient("https://api.mainnet-beta.solana.com")

	acct, err := cli.GetAccountInfo(context.Background(), addr)
	// handle `err`

	json.NewEncoder(os.Stdout).Encode(m)
// {
//   "context": {
//     "Slot": 47836700
//   },
//   "value": {
//     "lamports": 1461600,
//     "data": {
//       "data": "AQAAABzjWe1aAS4E+hQrnHUaHF6Hz9CgFhuchf/TG3jN/Nj2gCa3xLwWAAAGAQEAAAAqnl7btTwEZ5CY/3sSZRcUQ0/AjFYqmjuGEQXmctQicw==",
//       "encoding": "base64"
//     },
//     "owner": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
//     "executable": false,
//     "rentEpoch": 109
//   }
// }

```

## Examples

### Reference

 * RPC
	* [Get Recent Blockhash](./example_rpc_get_recent_blockhash_test.go)
 * WebSocket
	* [Account Subscribe](./example_ws_account_subscribe_test.go)

### Running

The easiest way to see the actual output for a given example is to add a line
`// Output: any` at the very end of the test, looks like this for
`ExampleRPC_GetRecentBlockhash` file ([example_rpc_get_recent_blockhash_test.go](./example_rpc_get_recent_blockhash_test.go)):

```
	...

    fmt.Println(string(bytes))
    // Output: any
}
```

This tells `go test` that it can execute this test correctly. Then, simply
run only this example:

    go test -run ExampleRPC_GetRecentBlockhash

Replacing `ExampleRPC_GetRecentBlockhash` with the actual example name you want to try
out where line `// Output: any` was added.

This will run the example and compares the standard output with the `any` which
will fail. But it's ok an expected, so you can see the actual output
printed to your terminal.

> WebSocket examples runs for a 1 minute then exits, you will not see anything until the example finish
> RPC URL to use can be specified by using environment variable `SOLANA_GO_RPC_URL`, defaults to https://api.mainnet-beta.solana.com.
> WS URL to use can be specified by using environment variable `SOLANA_GO_WS_URL`, defaults to ws://api.mainnet-beta.solana.com.

## Contributing

**Issues and PR in this repo related strictly to the solana go library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

This codebase uses unit tests extensively, please write and run tests.

## License

[Apache 2.0](LICENSE)
