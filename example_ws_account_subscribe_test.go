package solana_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/rpc"
	"github.com/streamingfast/solana-go/rpc/ws"
)

func ExampleWS_AccountSubscribe() {
	// Let's this example run for 1m than shutdown everything, not required in your own production code of course
	ctx, shutdown := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Minute)
		shutdown()
	}()

	websocket := ws.NewClient(getWSURL(), false)

	fmt.Printf("Connecting to websocket %s\n", getWSURL())
	err := websocket.Dial(ctx)
	if err != nil {
		panic(fmt.Errorf("websocket dial: %w", err))
	}

	// Ensure we cleanup the connection before the program exits
	defer websocket.Close()

	serumEventQueueV3 := solana.MustPublicKeyFromBase58("D1hpxetuGzfz2mSf3US6F7QHjmmA3A5Q1EUJ3Qk5E1ZG")

	fmt.Printf("Subscribing to account data changes for %s\n", serumEventQueueV3)
	subscription, err := websocket.AccountSubscribe(serumEventQueueV3, rpc.CommitmentConfirmed)
	if err != nil {
		panic(fmt.Errorf("websocket dial: %w", err))
	}

	fmt.Println("Consuming account subscription messages")
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping consuming subsription event")
				return
			default:
				message, err := subscription.Recv(ctx)
				if err != nil {
					if err == context.Canceled {
						// Context canceled, we terminate the example execution normally
						return
					}

					panic(fmt.Errorf("account subsciption: %w", err))
				}

				out, err := json.Marshal(message.(*ws.AccountResult))
				if err != nil {
					panic(fmt.Errorf("json marshal response: %w", err))
				}

				fmt.Println(string(out))
			}
		}
	}()

	// Wait until the actual context terminates, should be something else in your own case
	<-ctx.Done()
}
