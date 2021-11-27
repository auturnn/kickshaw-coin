package blockchain

import (
	"fmt"
	"sync"

	"github.com/auturnn/kickshaw-coin/db"
	"github.com/auturnn/kickshaw-coin/utils"
)

const (
	defaultDiffculty   int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var bc *blockchain
var once sync.Once

func (bc *blockchain) recalculrateDifficulty() int {
	allBlocks := bc.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculratedBlock := allBlocks[difficultyInterval-1]
	//actualTime : 현재 생성되는 블럭의 생성주기
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculratedBlock.Timestamp / 60)
	//expectedTime : 의도한 블럭생성주기
	expectedTime := difficultyInterval * blockInterval
	if actualTime < (expectedTime - allowedRange) {
		return bc.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return bc.CurrentDifficulty - 1
	}
	return bc.CurrentDifficulty
}

func (bc *blockchain) difficulty() int {
	if bc.Height == 0 {
		return defaultDiffculty
	} else if bc.Height%difficultyInterval == 0 {
		//recalculrate the difficulty
		return bc.recalculrateDifficulty()
	} else {
		return bc.CurrentDifficulty
	}

}

func (bc *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(bc))
}

func (bc *blockchain) AddBlock() {
	block := createBlock(bc.NewestHash, bc.Height+1)
	bc.NewestHash = block.Hash
	bc.Height = block.Height
	bc.CurrentDifficulty = block.Difficulty
	bc.persist()
}

func (bc *blockchain) restore(data []byte) {
	utils.FromBytes(bc, data)
}

func (bc *blockchain) BalanceByAddress(address string) (amount int) {
	txOuts := bc.UTxOutsByAddress(address)
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

//Unspent Transaction Outputs By Address
func (bc *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range bc.Blocks() {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BlockChain() *blockchain {
	if bc == nil {
		once.Do(func() {
			bc = &blockchain{Height: 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				bc.AddBlock()
			} else {
				bc.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash: %s\n", bc.NewestHash)
	return bc
}

func (bc *blockchain) Blocks() (blocks []*Block) {
	hashCursor := bc.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
