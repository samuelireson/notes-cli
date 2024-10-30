/*
Copyright © 2024 Samuel Ireson <samuelireson@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Takes as input the full chapter path, e.g., notes/mlnn/chapters/introduction.tex
// Outputs the path to the converted mdx file, e.g., site/src/content/docs/mlnn/introduction.mdx
func generateMDXOutputPath(texChapterPath string) string {
	chapterDir, texChapterName := filepath.Split(texChapterPath)
	chapterName := strings.TrimSuffix(texChapterName, filepath.Ext(texChapterName))
	mdxOutputDir := strings.Replace(chapterDir, LaTeXDir, MDXDir, 1)

	err := os.MkdirAll(mdxOutputDir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("fatal error creating output directory: %w", err))
	}

	mdxChapterName := chapterName + ".mdx"
	mdxOutputPath := filepath.Join(mdxOutputDir, mdxChapterName)
	return mdxOutputPath
}

// Takes as input the full chapter path, and the course bibliography
// Write the converted file to the output path
func processFile(texChapterPath string, courseBibliography bibliography) error {
	mdxOutputPath := generateMDXOutputPath(texChapterPath)

	texChapter, err := os.ReadFile(texChapterPath)
	if err != nil {
		return err
	}

	content := string(texChapter)
	content = convertTeXToMDX(content)
	content = convertCitationsToFootnotes(courseBibliography, content)
	content = addDownloadLinks(content, texChapterPath)
	mdxChapter := []byte(content)

	err = os.WriteFile(mdxOutputPath, mdxChapter, 0644)
	if err != nil {
		return err
	}

	log.Printf("%s converted successfully", texChapterPath)
	return nil
}

func processDirectory(coursePath string, courseBibliography bibliography) {
	chaptersPath := filepath.Join(coursePath, "chapters")
	chapters, err := os.ReadDir(chaptersPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading course directory: %w", err))
	}

	for _, chapter := range chapters {
		if filepath.Ext(chapter.Name()) == ".tex" {
			texChapterPath := filepath.Join(chaptersPath, chapter.Name())
			err := processFile(texChapterPath, courseBibliography)
			if err != nil {
				panic(fmt.Errorf("fatal error processing %s: %w", texChapterPath, err))
			}
		}
	}
}

var convertCmd = &cobra.Command{
	Use:   "convert [path to notes]",
	Short: "Convert notes from LaTeX to MDX.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		coursePath := args[0]
		courseBibliography := parseBibliography(coursePath)
		processDirectory(coursePath, courseBibliography)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
