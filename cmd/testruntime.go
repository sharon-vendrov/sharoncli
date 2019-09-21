/*
Copyright © 2019 Sharon Vendrov <sharon.vendrov1@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// testruntimeCmd represents the testruntime command
var testruntimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("testruntime called")
	},
}

func init() {
	testCmd.AddCommand(testruntimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testruntimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testruntimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
