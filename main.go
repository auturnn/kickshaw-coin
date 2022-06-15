package main

import (
	"github.com/auturnn/kickshaw-coin/cli"
	"github.com/auturnn/kickshaw-coin/utils"
)

func main() {
	utils.HasSystemPath()
	cli.Start()
}
