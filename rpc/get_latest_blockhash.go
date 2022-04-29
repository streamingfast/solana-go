package rpc

import (
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetRecentBlockhashResult struct {
	RPCContext
	Value *BlockhashResult `json:"value"`
}

type BlockhashResult struct {
	Blockhash            solana.PublicKey `json:"blockhash"` /* make this a `Hash` type, which is a copy of the PublicKey` type */
	LastValidBlockHeight bin.Uint64       `json:"lastValidBlockHeight"`
}

func (c *Client) GetLatestBlockhash(commitment CommitmentType) (out *GetRecentBlockhashResult, err error) {
	params := []interface{}{map[string]string{
		"commitment": string(commitment),
	}}
	err = c.DoRequest(&out, "getLatestBlockhash", params)
	return
}
