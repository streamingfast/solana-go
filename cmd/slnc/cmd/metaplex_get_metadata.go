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
	"context"
	"fmt"

	"github.com/streamingfast/solana-go/programs/metaplex"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
)

var metaplexGetMetadataCmd = &cobra.Command{
	Use:   "metadata {mint_addr}",
	Short: "Get Metaplex metadata account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		metadataAddr, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("decoding market addr: %w", err)
		}

		return getAndDisplayMetadata(ctx, metadataAddr)
	},
}

func init() {
	metaplexGetCmd.AddCommand(metaplexGetMetadataCmd)
}

func getAndDisplayMetadata(ctx context.Context, metadataAddr solana.PublicKey) error {
	cli := getClient()
	acc, err := cli.GetAccountInfo(ctx, metadataAddr)
	if err != nil {
		return fmt.Errorf("unable to retrieve account: %w", err)
	}
	metadata := &metaplex.Metadata{}
	err = metadata.Decode(acc.Value.Data)
	if err != nil {
		return fmt.Errorf("unable to decode metadata: %w", err)
	}

	fmt.Println()
	fmt.Printf("Metadata Addr: %s\n", metadataAddr.String())
	fmt.Printf("Mint: %s\n", metadata.Mint.String())
	fmt.Printf("Udpate Authority: %s\n", metadata.UpdateAuthority.String())
	fmt.Printf("Primary Sale Happened: %t\n", metadata.PrimarySaleHappened)
	fmt.Printf("Is Mutable: %t\n", metadata.IsMutable)
	fmt.Println("Data")
	fmt.Printf("> Name:%s \n", metadata.Data.Name)
	fmt.Printf("> Symbol: %s\n", metadata.Data.Symbol)
	fmt.Printf("> URI: %s\n", metadata.Data.URI)
	fmt.Printf("> Seller Basis Points: %d\n", metadata.Data.SellerFeeBasisPoints)
	if metadata.Data.Creators != nil && len(*metadata.Data.Creators) > 0 {
		fmt.Printf("> %d creators\n", len(*metadata.Data.Creators))
		for _, creator := range *metadata.Data.Creators {
			verified := "[ ]"
			if creator.Verified {
				verified = "[✅]"
			}

			fmt.Printf("> %s %d %s\n", creator.Address.String(), creator.Share, verified)
		}
	} else {
		fmt.Println("> No creators found")
	}

	return nil
}
