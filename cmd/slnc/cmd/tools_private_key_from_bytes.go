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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
)

var fromBytesPrivateKeytoolsCmd = &cobra.Command{
	Use:   "from-bytes",
	Short: "Converts a private key to byte array",
	Args:  cobra.ExactArgs(01),
	RunE: func(cmd *cobra.Command, args []string) error {
		privateKey := args[0]

		var values []uint8
		err := json.Unmarshal([]byte(privateKey), &values)
		if err != nil {
			return fmt.Errorf("unable to unmarshall byte array: %w", err)
		}

		pkey := solana.PrivateKey(values[:])
		fmt.Printf("%s\n", pkey.String())
		return nil
	},
}

func init() {
	privateKeytoolsCmd.AddCommand(fromBytesPrivateKeytoolsCmd)
}
