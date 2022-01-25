package cmd

import (
	"github.com/spf13/cobra"
)

var metaplexUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Metaplex objects",
}

func init() {
	metaplexCmd.AddCommand(metaplexUpdateCmd)
}
