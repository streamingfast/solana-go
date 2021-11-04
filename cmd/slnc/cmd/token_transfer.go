package cmd

import (
	"fmt"
	"strconv"

	"github.com/streamingfast/solana-go/rpc"

	"github.com/spf13/cobra"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/token"
)

var tokenTransferCmd = &cobra.Command{
	Use:   "transfer {recipient} {spl_token_account} {amount}",
	Short: "Transfer's a token to someone's associated token account",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		rpcCli := getClient(rpc.WithDebug())
		wsCli, err := getWsClient(ctx)
		if err != nil {
			return fmt.Errorf("unable to setup websocket client: %w", err)
		}
		vault := mustGetWallet()
		recipient, err := solana.PublicKeyFromBase58(args[0])
		if err != nil {
			return fmt.Errorf("decoding recipient addr: %w", err)
		}

		splTokenAccount, err := solana.PublicKeyFromBase58(args[1])
		if err != nil {
			return fmt.Errorf("decoding owner spl token accout: %w", err)
		}

		amount, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to decode amount: %w", err)
		}

		acct, err := rpcCli.GetAccountInfo(ctx, splTokenAccount)
		if err != nil {
			return fmt.Errorf("couldn't get spl token account data: %w", err)
		}

		account := &token.Account{}
		if err := account.Decode(splTokenAccount, acct.Value.Data); err != nil {
			return fmt.Errorf("unable to decode account information: %w", err)
		}

		if !account.IsInitialized {
			fmt.Println("uninitialized SPL token account. Data length", len(acct.Value.Data))
			return fmt.Errorf("uninitialized SPL token account. data length: %w", err)
		}

		var sender *solana.Account
		for _, privateKey := range vault.KeyBag {
			if privateKey.PublicKey() == account.Owner {
				sender = &solana.Account{PrivateKey: privateKey}
			}
		}

		if sender == nil {
			return fmt.Errorf("spl token account owner %q must be present in the vault to sign the send transaction", account.Owner.String())
		}
		fmt.Printf("Sending %d token %s (mint: %s) to %q from %q\n", amount, account.Key, account.Mint, recipient.String(), account.Owner.String())

		recipientSplTokenAccount, trxHash, err := token.TransferToken(ctx, rpcCli, wsCli, 1, account.Key, account.Mint, recipient, sender)
		if err != nil {
			return fmt.Errorf("unable to send transaction: %w", err)
		}

		fmt.Printf("Token Transfer successfull, with transaction hash: %s\n", trxHash)
		fmt.Printf("  Recipient SPL Token Account: %s\n", recipientSplTokenAccount.String())
		return nil
	},
}

func init() {
	tokenCmd.AddCommand(tokenTransferCmd)
}
