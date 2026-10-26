package main

import (
	"fmt"
	"os"

	"lem-in/colony"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}
	data, err := os.ReadFile(os.Args[1])
	if err == nil {
		err = colony.Run(data)
	}
	if err != nil {
		fmt.Println("ERROR: invalid data format")
	}
}
