package rpc

import (
	"fmt"
	"github.com/streamingfast/solana-go"
)

func (c *Client) RequestAirdrop(account *solana.PublicKey, lamport uint64, commitment CommitmentType) (signature string, err error) {

	obj := map[string]interface{}{
		"commitment": commitment,
	}

	params := []interface{}{
		account.String(),
		lamport,
		obj,
	}

	if err := c.DoRequest(&signature, "requestAirdrop", params...); err != nil {
		return "", fmt.Errorf("send transaction: rpc send: %w", err)
	}
	return
}
