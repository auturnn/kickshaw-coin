package main

import (
	"github.com/auturnn/kickshaw-coin/cli"
	"github.com/auturnn/kickshaw-coin/db"
)

func main() {
	defer db.Close()
	cli.Start()
	//http stateless ws stateful
}
