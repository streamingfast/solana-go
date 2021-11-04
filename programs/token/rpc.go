// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

import (
	"context"
	"fmt"

	"github.com/streamingfast/solana-go/rpc/confirm"
	"github.com/streamingfast/solana-go/rpc/ws"

	associatedtokenaccount "github.com/streamingfast/solana-go/programs/associated-token-account"

	"github.com/streamingfast/solana-go"

	"github.com/streamingfast/solana-go/rpc"
)

//go:generate rice embed-go

func FetchMints(ctx context.Context, rpcCli *rpc.Client) (out []*Mint, err error) {
	resp, err := rpcCli.GetProgramAccounts(
		ctx,
		PROGRAM_ID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: MINT_SIZE,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp empty... program account not found")
	}

	for _, keyedAcct := range resp {
		acct := keyedAcct.Account

		m := &Mint{}
		if err := m.Decode(acct.Data); err != nil {
			return nil, fmt.Errorf("unable to decode mint %q: %w", acct.Owner.String(), err)
		}
		out = append(out, m)

	}
	return
}

func FetchAccountsForOwner(ctx context.Context, rpcCli *rpc.Client, owner solana.PublicKey) (out []*Account, err error) {
	resp, err := rpcCli.GetProgramAccounts(
		ctx,
		PROGRAM_ID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{DataSize: ACCOUNT_SIZE},
				{Memcmp: &rpc.RPCFilterMemcmp{Offset: 32, Bytes: owner[:]}},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp empty... program account not found")
	}

	for _, keyedAcct := range resp {
		acct := keyedAcct.Account

		a := &Account{}
		if err := a.Decode(keyedAcct.Pubkey, acct.Data); err != nil {
			return nil, fmt.Errorf("unable to decode mint %q: %w", acct.Owner.String(), err)
		}
		out = append(out, a)

	}
	return
}

func TransferToken(ctx context.Context, rpcCli *rpc.Client, wsCli *ws.Client, amount uint64, senderSPLTokenAccount, mint, recipient solana.PublicKey, sender *solana.Account) (solana.PublicKey, string, error) {
	blockHashResult, err := rpcCli.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.PublicKey{}, "", fmt.Errorf("unable retrieve recent block hash: %w", err)
	}

	recipientSPLTokenAccount := associatedtokenaccount.MustGetAssociatedTokenAddress(mint, PROGRAM_ID, recipient)
	trx, err := solana.NewTransaction(
		[]solana.Instruction{
			associatedtokenaccount.NewCreateInstruction(
				sender.PublicKey(),
				recipientSPLTokenAccount,
				recipient,
				mint,
				PROGRAM_ID,
			),
			NewTransferInstruction(amount, senderSPLTokenAccount, recipientSPLTokenAccount, sender.PublicKey()),
		},
		blockHashResult.Value.Blockhash,
	)
	if err != nil {
		return solana.PublicKey{}, "", fmt.Errorf("unable to craft transaction: %w", err)
	}

	_, err = trx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		// create account need to be signed by the private key of the new account
		// that is not in the vault and will be lost after the execution.
		if key == sender.PublicKey() {
			return &sender.PrivateKey
		}
		return nil
	})

	trxHash, err := confirm.SendAndConfirmTransaction(ctx, rpcCli, wsCli, trx)
	if err != nil {
		return solana.PublicKey{}, "", fmt.Errorf("unable to send transaction: %w", err)
	}

	return recipientSPLTokenAccount, trxHash, nil
}
