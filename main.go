package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct{
	data string
	hash string
	prevHash string
}
type blockchain struct{
	block []block
}

func (bc *blockchain) getLastHash() string {
	if len(bc.block) > 0{
		return bc.block[len(bc.block)-1].hash
	}

	return ""
}

func (bc *blockchain) addBlock(data string) {
	newBlock := block{
		data: data,
		prevHash: bc.getLastHash(),
	}
	hashByte := sha256.Sum256([]byte(newBlock.data+ newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hashByte)
	bc.block = append(bc.block, newBlock)
}

func (bc *blockchain) listBlocks()  {
	for _, block := range bc.block{
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("PrevHash: %s\n\n", block.prevHash)
	}
}

func main()  {
	chain := blockchain{}
	chain.addBlock("hi")
	chain.addBlock("hi")

	chain.listBlocks()
}
