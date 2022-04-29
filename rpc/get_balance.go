package rpc

import (
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetBalanceResult struct {
	RPCContext
	Value bin.Uint64 `json:"value"`
}

func (c *Client) GetBalance(publicKey solana.PublicKey, commitment *CommitmentType) (out *GetBalanceResult, err error) {
	params := []interface{}{publicKey.String()}
	if commitment != nil {
		params = append(params, map[string]string{
			"commitment": string(*commitment),
		})
	}
	err = c.DoRequest(&out, "getBalance", params...)
	return
}
