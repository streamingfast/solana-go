package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
)

var fromKeypairAddressToolsCmd = &cobra.Command{
	Use:   "from-keypair",
	Short: "Converts a keypair file to an address",
	Args:  cobra.ExactArgs(01),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		keypair, err := readKeypair(ctx, args[0])
		if err != nil {
			return fmt.Errorf("unable to open keypair: %w", err)
		}

		pkey := solana.PrivateKey(keypair[:])
		fmt.Printf("%s\n", pkey.PublicKey().String())
		return nil
	},
}

func init() {
	addressToolsCmd.AddCommand(fromKeypairAddressToolsCmd)
}
