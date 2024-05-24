package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/flagship-io/abtasty-cli/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: %s
title: "%s"
filename: %s
---
`

const docVersion = "v1.0"
const docPath = "/docs/" + docVersion

func main() {
	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), filename)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))

		return fmt.Sprintf("%s/%s.md", docPath, base)
	}

	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "."+docPath, filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
