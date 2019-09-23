/*
Copyright © 2019 Sharon Vendrov <sharon.vendrov1@gmail.com>

*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "TODO",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Provide item to the create command")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
