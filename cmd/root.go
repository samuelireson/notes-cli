/*
Copyright © 2024 Samuel Ireson <samuelireson@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "notes",
	Short: "A CLI to manage notes.",
	Long: `A CLI for managing notes.

Convert your notes to publish them on a website.

Bulk compile LaTeX notes.

Scaffold projects quickly.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("texDir", "tex", "LaTeX notes directory")
	rootCmd.PersistentFlags().String("mdxDir", "mdx", "MDX notes directory")
	cobra.CheckErr(viper.BindPFlag("texDir", rootCmd.PersistentFlags().Lookup("texDir")))
	cobra.CheckErr(viper.BindPFlag("mdxDir", rootCmd.PersistentFlags().Lookup("mdxDir")))
	viper.SetDefault("texDir", "PATH TO BASE LATEX NOTES DIRECTORY")
	viper.SetDefault("mdxDir", "PATH TO BASE MDX NOTES DIRECTORY")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Writing the default config to config.yaml")
			err := viper.SafeWriteConfig()
			if err != nil {
				panic(fmt.Errorf("fatal error writing config file: %w", err))
			}
			log.Println("Edit the config file to run further commands")
			os.Exit(0)
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}
