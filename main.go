package main

import (
	"github.com/auturnn/kickshaw-coin/blockchain"
	"github.com/auturnn/kickshaw-coin/cli"
)

func main()  {
	// blockchain.BlockChain().AddBlock("First")
	// blockchain.BlockChain().AddBlock("second")
	// blockchain.BlockChain().AddBlock("third")
	blockchain.BlockChain()
	cli.Start()
}