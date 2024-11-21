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
	Use:   "notes-cli [command] [flags]",
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

var AsidePath string
var BibliographyPath string
var LaTeXDir string
var MDXDir string
var PDFDir string
var SiteURL string

func init() {
	rootCmd.PersistentFlags().StringVar(&BibliographyPath, "bibliographyPath", "bibliography.bib", "Bibliography path")
	rootCmd.PersistentFlags().StringVar(&LaTeXDir, "texDir", "PATH TO LATEX NOTES", "LaTeX note directory")
	rootCmd.PersistentFlags().StringVar(&MDXDir, "mdxDir", "PATH TO MDX NOTES", "MDX note directory")
	rootCmd.PersistentFlags().StringVar(&PDFDir, "pdfDir", "PATH TO PDF DIRECTORY", "Path for pdf output")
	rootCmd.PersistentFlags().StringVar(&SiteURL, "siteUrl", "SITE URL", "Url for public website")
	cobra.CheckErr(viper.BindPFlag("bibliographyPath", rootCmd.PersistentFlags().Lookup("bibliographyPath")))
	cobra.CheckErr(viper.BindPFlag("mdxDir", rootCmd.PersistentFlags().Lookup("mdxDir")))
	cobra.CheckErr(viper.BindPFlag("pdfDir", rootCmd.PersistentFlags().Lookup("pdfDir")))
	cobra.CheckErr(viper.BindPFlag("siteUrl", rootCmd.PersistentFlags().Lookup("siteUrl")))
	cobra.CheckErr(viper.BindPFlag("texDir", rootCmd.PersistentFlags().Lookup("texDir")))

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

	BibliographyPath = viper.GetString("bibliographyPath")
	LaTeXDir = viper.GetString("texDir")
	MDXDir = viper.GetString("mdxDir")
	PDFDir = viper.GetString("pdfDir")
	SiteURL = viper.GetString("siteUrl")
}
