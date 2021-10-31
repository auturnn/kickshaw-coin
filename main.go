package main

import (
	"github.com/auturnn/kickshaw-coin/blockchain"
)

func main()  {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("SecondBlock?")
	chain.AllBlocksPrint()
	chain.FindBlock(1)
}
