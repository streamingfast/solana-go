# Change log

The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html). See
[MAINTAINERS.md](./MAINTAINERS.md) for instructions to keep up to
date.

```
⚠️ solana-go works using SemVer but in 0 version, which means that the 'minor' will be changed when some broken changes are introduced into the application, and the 'patch' will be changed when a new feature with new changes is added or for bug fixing. As soon as v1.0.0 be released, solana-go will start to use SemVer as usual.
```

## Unreleased

## [v0.5.0](https://github.com/streamingfast/solana-go/releases/v0.4.0) (Feb 02, 2022)

### Breaking

* Moved `cmd/slnc` (and related packages `cli` and `vault`) to its own repository.

### Changes

* f05d17b Added listing holders for a given mint (token list holders)
* 73b6d19 Added mint token and update metadat support
* d8a67f2 Added mint verfied
* 1aaca1f Deleted missing slnc cmd
* 3b96fe1 Fixed a bit changelog
* 5d9b3b5 Fixed goreleaser project name
* 3f86ba9 Removed all .envrc files
* def18e6 Serialize Data the same way it is deserialized.
* 23a4683 Downgrade logging

## [v0.4.0](https://github.com/streamingfast/solana-go/releases/v0.4.0) (Nov 25, 2021)

* d541cb7 Reduced logging verbosity of some elements in `rpc.Client`.

## [v0.3.0](https://github.com/streamingfast/solana-go/releases/v0.3.0) (Nov 18, 2021)

* a361bc6 Reconnection logic when RPC nodes go down.

### Changes

## [v0.2.0](https://github.com/streamingfast/solana-go/releases/v0.2.0) (Nov 27, 2020)

* d64d73f Add program id for token registry mainnet
* e22f933 Added dex instruction decoding
* 37d477e Added parsed transaction support
* 8609824 Added todo
* 67de838 Added vault back
* ea20128 Adding dist to ignore
* 5069128 Fix missing signature on the new create account instruction
* de90e6c Fix readme instructions
* 5fa02f1 New text decoder New decode use into get transactions
* fa1eabd Updated rice boxes
* 0fd1895 move `slnc get token-meta` to `slnc spl get meta`
* edffba3 add more detail spl-token with string decoding
* 6bcac34 added programID resolver on token registry
* 037c9c5 added Decode and DecodeFromBase64 func to Event
* 062f5eb added Filled and Side func to Event
* 58bade2 added `token` program structs and decoding code
* bba400d added account keys for serum
* 90a4641 added ascii handling to tokenmeta
* 5069c6c added cmd request-airdrop
* b8e0e98 added decode func to OrderBook
* ef91119 added decode func to slotSubscribe and fix the unsubscribe functionality
* 965d5a6 added get transaction cli cmd with decoded instruction for serum and system.
* 5a21265 added more debug logs
* d7226cd added ping message so ws client do not get disconnected from the server
* 6ca7856 added return to error handler
* 1e94cdf added serum test and cleaned up struc
* 58d341c added serum test and cleaned up struc
* 7a8d266 added token-meta to slnc. still missing the program ID to be completed
* c666803 added website resize symbol and fix meta data size
* a519599 added website to register token instruction and fix cmd
* bbcd748 all test are passing
* 2f2a32f bump binary to fix instruction decoding
* 5947757 clean up serum types
* 026969c cleaned up error handling
* b781f1b downgrade logger
* d863779 fix MustPublicKeyFromBase58 that was not panicing
* ef7e72b fix bin= to bin: tag syntax
* 37a472e fix compilation issue
* 44cbbdb fix decoding of TokenMeta!
* f15f138 fix get token meta
* 3d0443c fix interface
* fcc62ec fix missing error handling
* 8e6e3c2 fix package import after moving program related stuff to folder programs
* cc80ff5 fix panic
* b893a1c fix payer key
* ea1c459 fix program entry
* 63e2f71 fix set account inerface
* c4868dd fix set accounts interface... for the LAST TIME
* f9fe8d5 fix variant def of token registry
* 4d145fe get spl token owner output
* eb1b380 instructions not associated with a decoder are now printed as `raw` instruction
* 1a6dcbd logo is not 64 byte long
* 8d95c0f major refactor cli
* adfd8f2 move graphql server to dfuse-solana
* ec87d95 move traceEnabled to logging.IsTraceEnabled
* 57107bb moved AccountSettable to solana interface.
* b18ae61 now able to decode EvenQueues
* d49de8c refactor to by able to unsubscribe
* 76687fa refactored ws
* 81f4290 refactored ws connection
* 0a7c9bc remove most of the remaining struc ref
* 60094b7 remove slnc dependancy
* 84b120e remove struc dependancy
* 6a84567 replace struc by dfuse bin. Some test are failing
* e513c61 rice boxed schema.graphql
* 8672711 serving graphql endpoint. still need to serve static stuff
* 5f67ae0 serving static stuff
* fb91647 set some log to debug
* 68f8c70 skip bin accountst to serum instruction
* 35b0380 slnc `spl token register` working.
* 6d56926 spl-token with string decoding

## [v0.1.0] 2020-11-09

First release

## Includes the following features:

* Get basic information from the chain about accounts, balances, etc.
* Issue SOL native token transfer
* Issue SPL token transfers
* Get Project SERUM markets list and live market data