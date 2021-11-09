package main

import (
	"github.com/auturnn/kickshaw-coin/cli"
	"github.com/auturnn/kickshaw-coin/db"
)

func main()  {
	// blockchain.BlockChain().AddBlock("First")
	// blockchain.BlockChain().AddBlock("second")
	// blockchain.BlockChain().AddBlock("third")
	defer db.Close()
	cli.Start()
}