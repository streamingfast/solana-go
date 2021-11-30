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
	"strings"

	"github.com/spf13/viper"

	"github.com/streamingfast/solana-go"

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go/programs/token"
)

var tokenListHoldersCmd = &cobra.Command{
	Use:   "holders {mint}",
	Short: "Lists token holders for a given mint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcCli := getClient()
		toCSV := viper.GetBool("token-list-holders-cmd-to-csv")

		ownerAddr, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("decoding owner addr: %w", err)
		}

		accounts, err := token.FetchAccountHolders(cmd.Context(), rpcCli, ownerAddr)
		if err != nil {
			return fmt.Errorf("unable to retrieve mints: %w", err)
		}

		if toCSV {
			fmt.Println("address,mint,owner,amount")
			for _, a := range accounts {
				line := []string{
					fmt.Sprintf("%s", a.Key.String()),
					fmt.Sprintf("%s", a.Mint.String()),
					fmt.Sprintf("%s", a.Owner.String()),
					fmt.Sprintf("%d", a.Amount),
				}
				fmt.Println(strings.Join(line, ","))
			}
		} else {
			out := []string{"Address | Mint | Owner | Amount"}
			for _, a := range accounts {
				line := []string{
					fmt.Sprintf("%s", a.Key.String()),
					fmt.Sprintf("%s", a.Mint.String()),
					fmt.Sprintf("%s", a.Owner.String()),
					fmt.Sprintf("%d", a.Amount),
				}
				out = append(out, strings.Join(line, " | "))
			}
			fmt.Println(columnize.Format(out, nil))
		}
		return nil
	},
}

func init() {
	tokenListHoldersCmd.Flags().Bool("to-csv", false, "outputs the data in csv format")
	tokenListCmd.AddCommand(tokenListHoldersCmd)
}
