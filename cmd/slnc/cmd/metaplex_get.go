package cmd

import (
	"github.com/spf13/cobra"
)

var metaplexGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Metaplex objects",
}

func init() {
	metaplexCmd.AddCommand(metaplexGetCmd)
}
