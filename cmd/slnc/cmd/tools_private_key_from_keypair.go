package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/streamingfast/dstore"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
)

var fromKeypairPrivateKeyToolsCmd = &cobra.Command{
	Use:   "from-keypair",
	Short: "Converts keypair file to a base58 private key",
	Args:  cobra.ExactArgs(01),
	RunE: func(cmd *cobra.Command, args []string) error {
		keypair, err := readKeypair(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("unable to open keypair: %w", err)
		}
		pkey := solana.PrivateKey(keypair[:])
		fmt.Printf("%s\n", pkey.String())
		return nil
	},
}

func init() {
	privateKeytoolsCmd.AddCommand(fromKeypairPrivateKeyToolsCmd)
}

func readKeypair(ctx context.Context, keypairFilePath string) ([]uint8, error) {
	file, _, _, err := dstore.OpenObject(ctx, keypairFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open keypair: %w", err)
	}

	cnt, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file pair content: %w", err)
	}

	var values []uint8
	if err = json.Unmarshal([]byte(cnt), &values); err != nil {
		return nil, fmt.Errorf("unable to unmarshall byte array: %w", err)
	}
	return values, nil
}
