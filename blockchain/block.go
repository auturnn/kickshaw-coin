package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/utils"
)

type Block struct{
	Data 	 string `json:"data"`
	Hash 	 string `json:"hash"`
	PrevHash string `json:"prevHash"`
	Height   int 	`json:"height"`
}

func (b *Block) toBytes() []byte  {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	err := encoder.Encode(b)
	utils.HandleError(err)
	return blockBuffer.Bytes()
}

func (b *Block) persist()  {
	db.SaveBlock(b.Hash, b.toBytes())
}

func createBlock(data, prevHash string, height int) *Block {
	block := Block{
		Data: data,
		PrevHash: prevHash,
		Hash: "",
		Height: height,
	}
	payload := block.Data + block.Hash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return &block
}