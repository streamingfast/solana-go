package cmd

import (
	"fmt"

	"github.com/streamingfast/solana-go/rpc"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/token"
)

var tokenCloseAccountCmd = &cobra.Command{
	Use:   "close-account {account} {destination} {owner}",
	Short: "Close an account by transferring all its SOL to the destination account, Non-native accounts may only be closed if its token amount is zero.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		rpcCli := getClient(rpc.WithDebug())
		wsCli, err := getWsClient(ctx)
		if err != nil {
			return fmt.Errorf("unable to setup websocket client: %w", err)
		}
		vault := mustGetWallet()
		accountKey, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("decoding account key: %w", err)
		}

		destinationKey, err := solana.PublicKeyFromBase58(args[1])
		if err != nil {
			return fmt.Errorf("decoding destination key: %w", err)
		}

		ownerKey, err := solana.PublicKeyFromBase58(args[2])
		if err != nil {
			return fmt.Errorf("decoding owner key: %w", err)
		}

		if _, err = rpcCli.GetAccountInfo(ctx, accountKey); err != nil {
			return fmt.Errorf("couldn't get account data: %w", err)
		}

		var signer *solana.Account
		for _, privateKey := range vault.KeyBag {
			if privateKey.PublicKey() == ownerKey {
				signer = &solana.Account{PrivateKey: privateKey}
			}
		}

		if signer == nil {
			return fmt.Errorf("spl token account owner %q must be present in the vault to sign the send transaction", ownerKey.String())
		}
		fmt.Printf("Closing account %s, sending remaining lamports to %s\n", accountKey.String(), destinationKey.String())

		trxHash, err := token.DoCloseAccount(ctx, rpcCli, wsCli, accountKey, destinationKey, ownerKey, signer)
		if err != nil {
			return fmt.Errorf("unable to send transaction: %w", err)
		}

		fmt.Printf("Close Account successful, with transaction hash: %s\n", trxHash)
		return nil
	},
}

func init() {
	tokenCmd.AddCommand(tokenCloseAccountCmd)
}
