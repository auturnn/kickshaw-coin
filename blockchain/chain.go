package blockchain

import (
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

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(bc *blockchain, targetID string) *Tx {
	for _, tx := range Txs(bc) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

func getDifficulty(bc *blockchain) int {
	if bc.Height == 0 {
		return defaultDiffculty
	} else if bc.Height%difficultyInterval == 0 {
		//recalculrate the difficulty
		return recalculrateDifficulty(bc)
	} else {
		return bc.CurrentDifficulty
	}
}

func recalculrateDifficulty(bc *blockchain) int {
	allBlocks := Blocks(bc)
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

func persistBlockchain(bc *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(bc))
}

func (bc *blockchain) AddBlock() {
	block := createBlock(bc.NewestHash, bc.Height+1, getDifficulty(bc))
	bc.NewestHash = block.Hash
	bc.Height = block.Height
	bc.CurrentDifficulty = block.Difficulty
	persistBlockchain(bc)
}

func (bc *blockchain) restore(data []byte) {
	utils.FromBytes(bc, data)
}

func BalanceByAddress(address string, bc *blockchain) (amount int) {
	txOuts := UTxOutsByAddress(address, bc)
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func UTxOutsByAddress(address string, bc *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(bc) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(bc, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
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
	once.Do(func() {
		bc = &blockchain{Height: 0}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			bc.AddBlock()
		} else {
			bc.restore(checkpoint)
		}
	})
	return bc
}

func Blocks(bc *blockchain) (blocks []*Block) {
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
