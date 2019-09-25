/*
Copyright © 2019 Sharon Vendrov <sharon.vendrov1@gmail.com>

*/
package cmd

import (
	"github.com/sharon-vendrov/sharoncli/pkg/logic"
	"github.com/spf13/cobra"
)

var pipelineName string

// testruntimeCmd represents the testruntime command
var testruntimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "execute pipeline",
	Long:  `execute pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		err := logic.ExecutePipeline(pipelineName)
		if err != nil {
			panic("fail to run pipeline")
		}
	},
}

func init() {
	testruntimeCmd.Flags().StringVar(&pipelineName, "name", "", "pipeline name")
	testCmd.AddCommand(testruntimeCmd)

}
