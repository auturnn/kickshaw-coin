package blockchain

import (
	"sync"
)

type blockchain struct{
	NewestHash string `json:"newestHash"`
	Height 	   int `json:"height"`
}

var bc *blockchain
var once sync.Once

func (bc *blockchain) AddBlock(data string)  {
	block := createBlock(data, bc.NewestHash, bc.Height)
	bc.NewestHash = block.Hash
	bc.Height = bc.Height

}

func BlockChain() *blockchain {
	if bc == nil{
		once.Do(func ()  {
			bc = &blockchain{NewestHash: "", Height: 0}
			bc.AddBlock("hi")
		})
	}
	return bc
}