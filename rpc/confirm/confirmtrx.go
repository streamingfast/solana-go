package confirm

import (
	"context"
	"fmt"

	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/rpc"
	"github.com/streamingfast/solana-go/rpc/ws"
)

func SendAndConfirmTransaction(ctx context.Context, rppClient *rpc.Client, wsClient *ws.Client, transaction *solana.Transaction) (signature string, err error) {
	sig, err := rppClient.SendTransaction(
		transaction,
		&rpc.SendTransactionOptions{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		})
	if err != nil {
		return "", fmt.Errorf("unable to send transction: %w", err)
	}

	sub, err := wsClient.SignatureSubscribe(
		sig,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return sig, fmt.Errorf("failed to subscribe to signature: %w", err)
	}
	defer sub.Unsubscribe()
	for {
		res, err := sub.Recv(ctx)
		if err != nil {
			return sig, err
		}
		signResult := res.(*ws.SignatureResult)
		if signResult.Value.Err != nil {
			return sig, fmt.Errorf("transaction confirmation failed: %v", signResult.Value.Err)
		} else {
			return sig, nil
		}
	}
}
