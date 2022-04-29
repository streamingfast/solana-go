package rpc

import (
	bin "github.com/streamingfast/binary"
)

func (c *Client) GetSlot(commitment *CommitmentType) (uint64, error) {
	var params []interface{}
	if commitment != nil {
		params = append(params, string(*commitment))
	}

	var out bin.Uint64
	err := c.DoRequest(&out, "getSlot", params...)
	if err != nil {
		return 0, err
	}
	return uint64(out), nil
}
