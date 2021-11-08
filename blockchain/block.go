package blockchain

import (
	"crypto/sha256"
	"errors"
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

var ErrNotFound = errors.New("block not found")

func (b *Block) persist()  {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte)  {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(data, prevHash string, height int) *Block {
	block := &Block{
		Data: data,
		PrevHash: prevHash,
		Hash: "",
		Height: height,
	}
	payload := block.Data + block.Hash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}