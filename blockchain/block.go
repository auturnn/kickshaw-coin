package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/utils"
)

const difficulty int = 2

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

var ErrNotFound = errors.New("block not found")

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
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

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("Hash:=%s\nTarget:=%s\nNonce:=%d\n\n", hash, target, b.Nonce)
		if !strings.HasPrefix(hash, target) {
			b.Nonce++
		} else {
			b.Hash = hash
			break
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		PrevHash:   prevHash,
		Hash:       "",
		Height:     height,
		Difficulty: BlockChain().difficulty(),
		Nonce:      0,
	}
	block.mine()
	//mining이 언제 끝날지 모르기 때문에 끝난후에 이떄 추가
	block.Transactions = Mempool.TxToConfirm()
	block.persist()
	return block
}

func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

//Coinbase는 Blockchain에서 채굴자에게 주는 보상
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}

	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	return &tx
}
