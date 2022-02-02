package solana_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streamingfast/solana-go/rpc"
)

func ExampleRPC_GetRecentBlockhash() {
	client := rpc.NewClient(getRPCURL())

	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		panic(fmt.Errorf("get info: %w", err))
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("json marshal response: %w", err))
	}

	fmt.Println(string(bytes))
}
