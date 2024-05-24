package main

import (
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/cmd"
	"github.com/spf13/cobra/doc"
)

const docPath = "/docs/documentation"

func main() {
	os.Mkdir("./docs/documentation", os.ModePerm)
	err := doc.GenMarkdownTree(cmd.RootCmd, "."+docPath)
	if err != nil {
		log.Fatal(err)
	}
}
