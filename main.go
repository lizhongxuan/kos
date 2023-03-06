package main

import (
	"kos/cmd"
	"kos/pkg/utils"
)

var logger = utils.GetLogger("kos")

func main() {
	err := cmd.Main()
	if err != nil {
		logger.Fatal(err)
	}
}
