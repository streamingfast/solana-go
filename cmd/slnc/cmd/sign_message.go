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
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
)

var signMessageCmd = &cobra.Command{
	Use:   "sign {keyfile} {message_string}",
	Short: "sign a message using keyfile",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		priv, err := solana.PrivateKeyFromSolanaKeygenFile(args[0])
		if err != nil {
			return err
		}
		sig, err := priv.Sign([]byte(args[1]))
		if err != nil {
			return err
		}
		sigb := [64]byte(sig)
		fmt.Println(hex.EncodeToString(sigb[:]))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(signMessageCmd)
}
