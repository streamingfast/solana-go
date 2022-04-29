package rpc

func (c *Client) GetMinimumBalanceForRentExemption(dataSize int) (lamport int, err error) {
	params := []interface{}{dataSize}
	err = c.DoRequest(&lamport, "getMinimumBalanceForRentExemption", params...)
	return
}
