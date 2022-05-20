package main

import (
	"github.com/auturnn/kickshaw-coin/cli"
	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/utils"
)

func main() {
	defer db.Close()
	utils.HasSystemPath()
	db.InitDB()
	cli.Start()
}
