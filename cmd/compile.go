/*
Copyright © 2024 Samuel Ireson <samuelireson@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func compileMaster(courseDir string) {
	masterPath := filepath.Join(courseDir, "master.tex")
	compileCommand := exec.Command(
		"latexmk",
		"-lualatex",
		"-cd",
		"-g",
		masterPath,
	)

	err := compileCommand.Run()
	if err != nil {
		panic(fmt.Errorf("fatal error compiling %s: %w", masterPath, err))
	}
	log.Printf("%s compiled successfully", masterPath)

	outputDir := filepath.Join(courseDir, "output")
	os.Mkdir(outputDir, os.ModePerm)

	masterPDFPath := strings.TrimSuffix(masterPath, filepath.Ext(masterPath)) + ".pdf"
	outputPDFPath := filepath.Join(outputDir, "master.pdf")
	err = os.Rename(masterPDFPath, outputPDFPath)
	if err != nil {
		log.Fatal(err)
	}
}

func compileChapter(chapter, courseDir string) {
	includeChapterPath := filepath.Join("chapters", chapter)
	chapterPath := filepath.Join(courseDir, includeChapterPath)
	useTeX := `-usepretex="\\includeonly{` + includeChapterPath + `}"`
	masterPath := filepath.Join(courseDir, "master.tex")
	compileCommand := exec.Command(
		"latexmk",
		"-lualatex",
		"-cd",
		"-g",
		useTeX,
		masterPath,
	)

	err := compileCommand.Run()
	if err != nil {
		panic(fmt.Errorf("fatal error compiling %s: %w", chapterPath, err))
	}

	outputDir := filepath.Join(courseDir, "output")
	os.Mkdir(outputDir, os.ModePerm)

	chapterName := strings.TrimSuffix(chapter, filepath.Ext(chapter))
	outputPDFPath := filepath.Join(outputDir, chapterName+".pdf")
	masterPDFPath := strings.TrimSuffix(masterPath, filepath.Ext(masterPath)) + ".pdf"
	err = os.Rename(masterPDFPath, outputPDFPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s compiled successfully", filepath.Join(courseDir, chapterPath))
}

func compileCourse(courseDir string) {
	chapterDir := filepath.Join(courseDir, "chapters")
	chapters, err := os.ReadDir(chapterDir)
	if err != nil {
		log.Fatal(err)
	}

	compileMaster(courseDir)

	for _, chapter := range chapters {
		if filepath.Ext(chapter.Name()) == ".tex" {
			compileChapter(chapter.Name(), courseDir)
		}
	}

	outputPath := filepath.Join(courseDir, "output")
	publicDownloadPath := PDFDir + courseDir

	err = os.RemoveAll(publicDownloadPath)
	if err != nil {
		panic(fmt.Errorf("fatal error on removing PDF directory: %w", err))
	}

	err = os.Rename(outputPath, publicDownloadPath)
	if err != nil {
		panic(fmt.Errorf("fatal error on moving PDFs: %w", err))
	}

	log.Println("Download paths synced")
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile LaTeX notes to pdf chapters.",
	Run: func(cmd *cobra.Command, args []string) {
		courseDir := args[0]
		compileCourse(courseDir)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
