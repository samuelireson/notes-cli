package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func convertFigure(figurePath string) {
	texFigurePath := filepath.Join(figurePath, "figure.tex")
	dviCommand := exec.Command(
		"latexmk",
		"-dvilua",
		"-cd",
		texFigurePath,
	)

	err := dviCommand.Run()
	if err != nil {
		panic(fmt.Errorf("fatal error compiling %s: %w", texFigurePath, err))
	}

	dviFigurePath := filepath.Join(figurePath, "figure.dvi")
	svgCommand := exec.Command("dvisvgm", "--font-format=TTF", dviFigurePath)

	err = svgCommand.Run()
	if err != nil {
		panic(fmt.Errorf("fatal error on converting %s: %w", dviFigurePath, err))
	}

	imgOutputDir := strings.TrimPrefix(figurePath, LaTeXDir)
	svgOutputPath := MDXDir + imgOutputDir + ".svg"

	err = os.Rename("figure.svg", svgOutputPath)
	if err != nil {
		panic(fmt.Errorf("fatal error moving figure"))
	}

	log.Printf("%s converted successfully", figurePath)
}

func convertFigures(coursePath string) {
	figureDir := filepath.Join(coursePath, "figures")

	figureNames, err := os.ReadDir(figureDir)
	if err == os.ErrNotExist {
		log.Print("No figures directory")
		return
	} else if err != nil {
		panic(fmt.Errorf("fatal error reading %s: %w", figureDir, err))
	}

	for _, figureName := range figureNames {
		figurePath := filepath.Join(figureDir, figureName.Name())
		convertFigure(figurePath)
	}
}
