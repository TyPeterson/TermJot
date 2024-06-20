package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TyPeterson/TermJot/cmd"
	"github.com/TyPeterson/TermJot/internal/core"
)

func main() {
	if os.Getenv("TERMJOT_TESTING") != "1" {
		err := core.Init()
		if err != nil {
			log.Fatalf("Error initializing core: %v", err)
		}
	}

	cmd.Execute()
	fmt.Println()
}
