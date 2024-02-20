package main

import (
	"log"
	"os"

	"github.com/scottjr632/sequoia/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
