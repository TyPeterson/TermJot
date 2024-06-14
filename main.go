package main

import (
	"log"
    "fmt"
	"github.com/TyPeterson/TermJot/cmd"
	"github.com/TyPeterson/TermJot/internal/core"
)

func main() {
	err := core.Init()
	if err != nil {
		log.Fatalf("Error initializing core: %v", err)
	}
	cmd.Execute()
    fmt.Println()
}
