package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct{
	Height int `json:"height"`
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
}

type blockchain struct{
	blocks []*Block
}

var bc *blockchain
var once sync.Once

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func (b *Block) getHash() {
	hash := sha256.Sum256([]byte(b.Data+b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)	
}

func createBlock(data string) *Block {	
	newBlock := Block{Data: data, PrevHash: getLastHash(), Height: len(GetBlockChain().blocks)+1}
	newBlock.getHash()
	return &newBlock
}

func (bc *blockchain) AddBlock(data string)  {
	bc.blocks = append(bc.blocks, createBlock(data))
}

func (bc *blockchain) AllBlocks() []*Block {
	return bc.blocks
}

func (bc *blockchain) GetBlock(height int) *Block {
	return bc.blocks[height-1]
}

func GetBlockChain() *blockchain {
	if bc == nil{
		once.Do(func ()  {
			bc = &blockchain{}
			bc.AddBlock("FirstBlock!!")
		})
	}
	return bc
}