package main

import (
	"fmt"
	"os"

	"lem-in/colony"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: invalid data format, missing file argument")
		return
	}
	if len(os.Args) > 2 {
		fmt.Println("ERROR: invalid data format, expected only one file argument")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: invalid data format, cannot read file: %v\n", err)
		return
	}

	if err := colony.Run(data); err != nil {
		fmt.Printf("ERROR: invalid data format, %v\n", err)
	}
}