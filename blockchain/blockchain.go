package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct{
	data string
	hash string
	prevHash string
}

type blockchain struct{
	block []*block
}

var bc *blockchain
var once sync.Once

func getLastHash() string {
	totalBlocks := len(GetBlockChain().block)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().block[totalBlocks-1].hash
}

func (b *block) getHash() {
	hash := sha256.Sum256([]byte(b.data+b.prevHash))
	b.hash = fmt.Sprintf("%x", hash)	
}

func createBlock(data string) *block {	
	newBlock := block{data: data, prevHash: getLastHash()}
	newBlock.getHash()
	return &newBlock
}

func (bc *blockchain) AddBlock(data string)  {
	bc.block = append(bc.block, createBlock(data))
}

func (bc *blockchain) AllBlocks() []*block {
	return bc.block
}

func (bc *blockchain) AllBlocksPrint()  {
	for i, block := range bc.block{
		fmt.Printf("Index number %d Block\n", i)
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("PrevHash: %s\n\n", block.prevHash)
	}
}

func (bc *blockchain) FindBlock(index int)  {
	if index >= len(bc.block){
		fmt.Printf("Error: %d is Over the length\n", index)
	} else{
		fmt.Println("Find Success!!!")
		fmt.Printf("Index number %d Block\n", index)
		fmt.Printf("Data: %s\n", bc.block[index].data)
		fmt.Printf("Hash: %s\n", bc.block[index].hash)
		fmt.Printf("PrevHash: %s\n\n", bc.block[index].prevHash)
	}
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