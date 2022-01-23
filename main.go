package main

import (
	"log"
	"os"

	"github.com/dungnh3/guardrails-assignment/cmd"
)

func main() {
	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
