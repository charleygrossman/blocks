package main

import (
	"blocks/cmd"
	"fmt"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
