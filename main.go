package main

import (
	"fmt"

	"github.com/verbalsaintmars/rmqctl/cmd"
)

func main() {
	if err := cmd.Cmd(); err != nil {
		fmt.Printf("rmqctl error: %s\n", err.Error())
	}
}
