/*
Copyright © 2024 Samuel Ireson <samuelireson@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert [path to notes]",
	Short: "Convert notes from LaTeX to MDX.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("convert called")
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
