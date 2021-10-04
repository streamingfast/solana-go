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

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/metaplex"
)

var metaplexGetMintCmd = &cobra.Command{
	Use:   "mint {mint_addr}",
	Short: "Get Metaplex metadata for a given mint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

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

		fmt.Println("Derive metadata account addr: ", metadataAddr.String())

		return getAndDisplayMetadata(ctx, metadataAddr)
	},
}

func init() {
	metaplexGetCmd.AddCommand(metaplexGetMintCmd)
}
