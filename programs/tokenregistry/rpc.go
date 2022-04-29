package tokenregistry

import (
	"fmt"

	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/rpc"
)

func GetTokenRegistryEntry(rpcCli *rpc.Client, mintAddress solana.PublicKey) (*TokenMeta, error) {
	resp, err := rpcCli.GetProgramAccounts(
		ProgramID(),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 5,
						Bytes:  mintAddress[:], // hackey to convert [32]byte to []byte
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp empty... cannot find account")
	}

	for _, keyedAcct := range resp {
		acct := keyedAcct.Account
		t, err := DecodeTokenMeta(acct.Data)
		if err != nil {
			return nil, fmt.Errorf("unable to decode token meta %q: %w", acct.Owner.String(), err)
		}
		return t, nil
	}
	return nil, rpc.ErrNotFound
}

func GetEntries(rpcCli *rpc.Client) (out []*TokenMeta, err error) {
	resp, err := rpcCli.GetProgramAccounts(
		ProgramID(),
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: TOKEN_META_SIZE,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp empty... cannot find accounts")
	}

	for _, keyedAcct := range resp {
		acct := keyedAcct.Account
		t, err := DecodeTokenMeta(acct.Data)
		if err != nil {
			return nil, fmt.Errorf("unable to decode token meta %q: %w", acct.Owner.String(), err)
		}
		out = append(out, t)
	}
	return out, nil
}
