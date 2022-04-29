package rpc

import (
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetSignaturesForAddressResult []*TransactionSignature

type TransactionSignature struct {
	BlockTime          bin.Uint64        `json:"blockTime"`
	ConfirmationStatus string            `json:"confirmationStatus"`
	Err                *TransactionError `json:"err"`
	Memo               *string           `json:"memo"`
	Signature          string            `json:"signature"`
	Slot               bin.Uint64        `json:"slot"`
}

type GetSignaturesForAddressOpts struct {
	Limit  *uint64 `json:"limit,omitempty"`
	Before *string `json:"before,omitempty"`
	Until  *string `json:"until,omitempty"`
}

func (c *Client) GetSignaturesForAddress(address solana.PublicKey, opts *GetSignaturesForAddressOpts) (out GetSignaturesForAddressResult, err error) {
	params := []interface{}{address.String()}
	if opts != nil {
		filter := map[string]interface{}{}
		if opts.Limit != nil {
			filter["limit"] = *opts.Limit
		}
		if opts.Before != nil {
			filter["before"] = *opts.Before
		}
		if opts.Until != nil {
			filter["until"] = *opts.Until
		}
		params = append(params, filter)
	}
	err = c.DoRequest(&out, "getSignaturesForAddress", params...)
	return
}

//func (c *Client) GetConfirmedSignaturesForAddress2(ctx context.Context, address solana.PublicKey, opts *GetConfirmedSignaturesForAddress2Opts) (out GetConfirmedSignaturesForAddress2Result, err error) {
//
//	params := []interface{}{address.String(), opts}
//
//	err = c.DoRequest(&out, "getConfirmedSignaturesForAddress2", params...)
//	return
//}
