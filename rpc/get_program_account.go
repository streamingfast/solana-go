package rpc

import (
	"github.com/streamingfast/solana-go"
)

type GetProgramAccountsResult []*KeyedAccount

type KeyedAccount struct {
	Pubkey  solana.PublicKey `json:"pubkey"`
	Account *Account         `json:"account"`
}

type GetProgramAccountsOpts struct {
	Commitment CommitmentType `json:"commitment,omitempty"`
	// Filter on accounts, implicit AND between filters
	Filters []RPCFilter `json:"filters,omitempty"`
}

func (c *Client) GetProgramAccounts(publicKey solana.PublicKey, opts *GetProgramAccountsOpts) (out GetProgramAccountsResult, err error) {
	obj := map[string]interface{}{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if len(opts.Filters) != 0 {
			obj["filters"] = opts.Filters
		}
	}

	params := []interface{}{publicKey, obj}

	err = c.DoRequest(&out, "getProgramAccounts", params...)
	return
}
