package blockchain

import (
	"fmt"
	"sync"

	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/utils"
)
type blockchain struct{
	NewestHash string `json:"newestHash"`
	Height 	   int `json:"height"`
}

var bc *blockchain
var once sync.Once

func (bc *blockchain) persist()  {
	db.SaveBlockChain(utils.ToBytes(bc))
}

func (bc *blockchain) AddBlock(data string)  {
	block := createBlock(data, bc.NewestHash, bc.Height+1)
	bc.NewestHash = block.Hash
	bc.Height = block.Height
	bc.persist()
}

func (bc *blockchain) restore(data []byte)   {
	utils.FromBytes(bc, data)
}

func BlockChain() *blockchain {
	if bc == nil{
		once.Do(func ()  {
			bc = &blockchain{NewestHash: "", Height: 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil{
				bc.AddBlock("Genesis")
			} else {
				bc.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash: %s\n", bc.NewestHash)
	return bc
}

func (bc *blockchain) Blocks() (blocks []*Block)  {
	hashCursor := bc.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else{
			break
		}
	}
	return blocks
}