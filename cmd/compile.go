/*
Copyright © 2024 Samuel Ireson <samuelireson@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile LaTeX notes to pdf chapters.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("compile called")
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
