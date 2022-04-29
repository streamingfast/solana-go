package rpc

import (
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetAccountInfoResult struct {
	RPCContext
	Value *Account `json:"value"`
}

type Account struct {
	Lamports   bin.Uint64       `json:"lamports"`
	Data       solana.Data      `json:"data"`
	Owner      solana.PublicKey `json:"owner"`
	Executable bool             `json:"executable"`
	RentEpoch  bin.Uint64       `json:"rentEpoch"`
}

func (c *Client) GetAccountDataIn(account solana.PublicKey, inVar interface{}) (err error) {
	resp, err := c.GetAccountInfo(account)
	if err != nil {
		return err
	}

	return bin.NewDecoder(resp.Value.Data).Decode(inVar)
}

func (c *Client) GetAccountInfo(account solana.PublicKey) (out *GetAccountInfoResult, err error) {
	obj := map[string]interface{}{
		"encoding": "base64",
	}
	params := []interface{}{account, obj}

	err = c.DoRequest(&out, "getAccountInfo", params...)
	if err != nil {
		return nil, err
	}

	if out.Value == nil {
		return nil, ErrNotFound
	}

	return out, nil
}
