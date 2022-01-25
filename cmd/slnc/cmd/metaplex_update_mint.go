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

package cmd

import (
	"fmt"
	"github.com/streamingfast/solana-go/rpc"
	"github.com/streamingfast/solana-go/rpc/confirm"
	"go.uber.org/zap"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/metaplex"
)

var metaplexUpdateMintCmd = &cobra.Command{
	Use:   "mint {mint_addr}",
	Short: "Get Metaplex metadata for a given mint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		vault := mustGetWallet()
		rpcClient := getClient()
		wsClient, err := getWsClient(ctx)
		if err != nil {
			return fmt.Errorf("unable to retrieve ws client: %w", err)
		}

		metaplexMetaProgramId := viper.GetString("metaplex-global-meta-program-id")
		programID, err := solana.PublicKeyFromBase58(metaplexMetaProgramId)
		if err != nil {
			return fmt.Errorf("unable to decode metaplex metadata programId %q: %w", metaplexMetaProgramId, err)
		}

		mintAddr, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("unable to decode mint addr: %w", err)
		}

		metadataAddr, err := metaplex.DeriveMetadataPublicKey(programID, mintAddr)
		if err != nil {
			return fmt.Errorf("unablt to decode metadata address: %w", err)
		}

		v := true
		updateAuthority := solana.MustPublicKeyFromBase58("EBXJg2Y7gdSbdsDWzM3tuu5fxkCGpg8HuDgGwgBgTXz3")
		creators := []metaplex.Creator{
			{
				Address:  updateAuthority,
				Verified: true,
				Share:    100,
			},
			{
				Address: solana.MustPublicKeyFromBase58("FsG3VK6mZtGRNzbxzGee7zTHuvFFbxKJYLDrFjCqfDNU"),
				Share:   0,
			},
		}
		updateMetadataInstruction := metaplex.NewUpdateMetadataAccountInstruction(
			programID,
			&metaplex.Data{
				Name:                 "Alfred 0042",
				Symbol:               "Alf",
				URI:                  "https://arweave.net/5kaIQKrgbLV5tztRttDB22YsuXj1aw4PZLU2HXy8qw8",
				SellerFeeBasisPoints: 500,
				Creators:             &creators,
			},
			nil,
			&v,
			metadataAddr,
			updateAuthority,
		)

		zlog.Debug("retrieving block hash")
		blockHashResult, err := rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			return fmt.Errorf("unable retrieve recent block hash: %w", err)
		}

		zlog.Debug("found block hash",
			zap.String("block_hash", blockHashResult.Value.Blockhash.String()),
		)

		trx, err := solana.NewTransaction([]solana.Instruction{
			updateMetadataInstruction,
		}, blockHashResult.Value.Blockhash)
		if err != nil {
			return fmt.Errorf("unable to craft transaction: %w", err)
		}

		zlog.Debug("signing metaplex transaction")
		_, err = trx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
			// create account need to be signed by the private key of the new account
			// that is not in the vault and will be lost after the execution.
			for _, k := range vault.KeyBag {
				if k.PublicKey() == key {
					return &k
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to sign transaction: %w", err)
		}

		zlog.Debug("sending transaction")

		trxHash, err := confirm.SendAndConfirmTransaction(ctx, rpcClient, wsClient, trx)
		if err != nil {
			return fmt.Errorf("unable to send transaction: %w", err)
		}

		fmt.Printf("Metaplex Metadata Updated, with transaction hash: %s\n", trxHash)
		fmt.Printf("  Mint Address: %s\n", mintAddr.String())
		fmt.Printf("  Metadata Address: %s\n", metadataAddr.String())
		fmt.Printf("Run `slnc metaplex get mint %s` to view metadata", mintAddr.String())
		return nil
	},
}

func init() {
	metaplexUpdateCmd.AddCommand(metaplexUpdateMintCmd)
}
