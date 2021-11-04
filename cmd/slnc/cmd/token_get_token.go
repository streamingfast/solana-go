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

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/token"
)

var tokenGetAccountCmd = &cobra.Command{
	Use:   "account {account_addr}",
	Short: "Retrieves token information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tokenAddress, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("decoding mint addr: %w", err)
		}

		client := getClient()

		acct, err := client.GetAccountInfo(ctx, tokenAddress)
		if err != nil {
			return fmt.Errorf("couldn't get account data: %w", err)
		}

		account := &token.Account{}
		if err := account.Decode(acct.Value.Data); err != nil {
			return fmt.Errorf("unable to retrieve int information: %w", err)
		}

		if !account.IsInitialized {
			fmt.Println("Uninitialized Account. Data length", len(acct.Value.Data))
			return nil
		}

		var out []string

		out = append(out, fmt.Sprintf("Amount | %d", account.Amount))
		out = append(out, fmt.Sprintf("Mint | %s", account.Mint.String()))
		out = append(out, fmt.Sprintf("Owner | %s", account.Owner.String()))

		fmt.Println(columnize.Format(out, nil))

		return nil
	},
}

func init() {
	tokenGetCmd.AddCommand(tokenGetAccountCmd)
}
