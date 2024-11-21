package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type regexPattern struct {
	captureGroup *regexp.Regexp
	replacement  string
}

type stringPattern struct {
	old string
	new string
}

var header = fmt.Sprintf(`---
title: $1
---
import Aside from '%s';
import Comments from '@components/Comments/Comments.tsx';
import { Tabs, TabItem, LinkButton } from '@astrojs/starlight/components';

`, AsidePath)

var basicRegexPatterns = []regexPattern{
	// document organisation
	{regexp.MustCompile(`\\chapter\{(.*?)\}`), header},
	{regexp.MustCompile(`\\section\{(.*?)\}`), "## $1"},
	{regexp.MustCompile(`\\subsection\{(.*?)\}`), "### $1"},
	{regexp.MustCompile(`\\chapterauthor\{.*?\}`), ""},

	// theorem environments
	{regexp.MustCompile(`\\begin\{corollary\}(\[.*?\])?`), "<Aside type='result' title='Corollary' name='$1'>"},
	{regexp.MustCompile(`\\begin\{definition\}(\[.*?\])?`), "<Aside type='definition' title='Definition' name='$1'>"},
	{regexp.MustCompile(`\\begin\{example\}(\[.*?\])?`), "<Aside type='example' title='Example' name='$1'>"},
	{regexp.MustCompile(`\\begin\{lemma\}(\[.*?\])?`), "<Aside type='result' title='Lemma' name='$1'>"},
	{regexp.MustCompile(`\\begin\{nonexample\}(\[.*?\])?`), "<Aside type='example' title='Non-example' name='$1'>"},
	{regexp.MustCompile(`\\begin\{notation\}(\[.*?\])?`), "<Aside type='comment' title='Notation' name='$1'>"},
	{regexp.MustCompile(`\\begin\{proposition\}(\[.*?\])?`), "<Aside type='result' title='Proposition' name='$1'>"},
	{regexp.MustCompile(`\\begin\{remark\}(\[.*?\])?`), "<Aside type='comment' title='Remark' name='$1'>"},
	{regexp.MustCompile(`\\begin\{theorem\}(\[.*?\])?`), "<Aside type='result' title='Theorem' name='$1'>"},
	{regexp.MustCompile(`\\end\{(definition|theorem|lemma|proposition|corollary|example|nonexample|notation|remark)\}`), "</Aside>"},

	// fonts and ligatures
	{regexp.MustCompile(`\\textbf\{(.*?)\}`), "<b> $1 </b>"},
	{regexp.MustCompile(`\\textit\{(.*?)\}`), "<em> $1 </em>"},

	// definitions - later extract to chapterwise index.
	{regexp.MustCompile(`\\defined\{(.*?)\}`), "<em> $1 </em>"},

	// figure fluff
	{regexp.MustCompile(`\\(begin|end)\{figure\}(\[!htb\])?`), ""},
	{regexp.MustCompile(`\s*?\\centering`), ""},
	{regexp.MustCompile(`\s*?\\caption\{(.*?)\}\n`), "\n<div style='width: 80%; font-style: italic; margin-inline: auto;'>Caption: $1 </div>"},
	{regexp.MustCompile(`\\includegraphics\{(.*?)/figure\.pdf\}`), "![$1](../figures/$1.svg)"},

	// maths environments
	{regexp.MustCompile(`\s*?\\begin\{align\*\}`), "\n$$$$\n\\begin{align*}"},
	{regexp.MustCompile(`\s*?\\end\{align\*\}`), "\n\\end{align*}\n$$$$"},
}

var stringPatterns = []stringPattern{
	//document organisation
	{"\\begin{chout}", "<div style='text-align: center; font-style: italic;'>"},
	{"\\end{chout}", "</div>"},

	// exercises
	{"\\begin{exercise}", "<Tabs>"},
	{"\\end{exercise}", "</Tabs>"},
	{"\\begin{problem}", "<TabItem label='Problem'>"},
	{"\\begin{solution}", "<TabItem label='Solution'>"},
	{"\\end{problem}", "</TabItem>"},
	{"\\end{solution}", "</TabItem>"},

	// badges
	{"\\basic", ":badge[Basic]{variant=success}"},
	{"\\intermediate", ":badge[Intermediate]{variant=caution}"},
	{"\\challenging", ":badge[Challenging]{variant=danger}"},

	// unordered lists
	{"\\begin{itemize}", ""},
	{"\\item", "-"},
	{"\\end{itemize}", ""},

	// fonts and ligatures
	{"`", "'"},

	// maths environments
	// {"\\begin{align*}", "$$\n\\begin{align*}"},
	// {"\\end{align*}", "\\end{align*}\n$$"},
}

func convertTeXToMDX(content string) string {
	for _, element := range stringPatterns {
		content = strings.ReplaceAll(content, element.old, element.new)
	}

	for _, element := range basicRegexPatterns {
		content = element.captureGroup.ReplaceAllString(content, element.replacement)
	}
	return content
}

func addDownloadLinks(content, texChapterPath string) string {

	chapterDir, texChapterName := filepath.Split(texChapterPath)
	courseDir := strings.TrimSuffix(chapterDir, "chapters/")
	chapterName := strings.TrimSuffix(texChapterName, filepath.Ext(texChapterName))
	baseDownloadPath := "/" + SiteURL + courseDir
	chapterDownloadPath := filepath.Join(baseDownloadPath, chapterName+".pdf")
	courseDownloadPath := filepath.Join(baseDownloadPath, "master.pdf")

	downloadLinkTemplate := fmt.Sprintf("<div style='display: flex; justify-content: space-around;'>\n\t<LinkButton target=\"_blank\" href=\"%s\" variant=\"secondary\" icon=\"document\" >Download</LinkButton>\n\t<LinkButton target=\"_blank\" href=\"%s\" variant=\"primary\" icon=\"open-book\" >Download</LinkButton>\n</div>", chapterDownloadPath, courseDownloadPath)

	content = strings.Replace(content, "\n\n", "\n\n"+downloadLinkTemplate, 1)

	return content
}
