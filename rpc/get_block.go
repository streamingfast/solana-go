package rpc

import (
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetBlockResult struct {
	// For blocs below ~ 900K RPC nodes often return nil
	BlockHeight *bin.Uint64 `json:"blockHeight"`
	// For blocs below ~ 900K RPC nodes often return nil
	BlockTime         *bin.Uint64      `json:"blockTime"`
	Blockhash         solana.PublicKey `json:"blockhash"`
	ParentSlot        bin.Uint64       `json:"parentSlot"`
	PreviousBlockhash solana.PublicKey `json:"previousBlockhash"`
	//Transactions      []struct {
	//	Meta        interface{} `json:"meta"`
	//	Transaction struct {
	//		Message struct {
	//			AccountKeys []string `json:"accountKeys"`
	//			Header      struct {
	//				NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	//				NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	//				NumRequiredSignatures       int `json:"numRequiredSignatures"`
	//			} `json:"header"`
	//			Instructions []struct {
	//				Accounts       []int  `json:"accounts"`
	//				Data           string `json:"data"`
	//				ProgramIdIndex int    `json:"programIdIndex"`
	//			} `json:"instructions"`
	//			RecentBlockhash string `json:"recentBlockhash"`
	//		} `json:"message"`
	//		Signatures []string `json:"signatures"`
	//	} `json:"transaction"`
	//} `json:"transactions"`
}

func (c *Client) GetBlock(slotNum uint64) (out *GetBlockResult, err error) {
	params := []interface{}{slotNum}
	err = c.DoRequest(&out, "getBlock", params)
	return
}
