package main

import (
	"fmt"
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
hide:
  - navigation
---
`

const docPath = "/docs/documentation"

func main() {

	filePrepender := func(filename string) string {
		return fmt.Sprintf(fmTemplate)
	}

	identity := func(s string) string { return s }

	os.Mkdir("./docs/documentation", os.ModePerm)
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "."+docPath, filePrepender, identity)
	if err != nil {
		log.Fatal(err)
	}
}
