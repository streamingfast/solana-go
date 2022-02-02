package solana_test

import "os"

func getRPCURL() string {
	rpcURL := os.Getenv("SOLANA_GO_RPC_URL")
	if rpcURL != "" {
		return rpcURL
	}

	return "https://api.mainnet-beta.solana.com"
}

func getWSURL() string {
	wsURL := os.Getenv("SOLANA_GO_WS_URL")
	if wsURL != "" {
		return wsURL
	}

	return "ws://api.mainnet-beta.solana.com"
}
