package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/auturnn/kickshaw-coin/utils"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (fakeDB) SaveBlock(hash string, data []byte) {}
func (fakeDB) SaveChain(data []byte)              {}
func (fakeDB) DeleteAllBlocks()                   {}

func TestBlockChain(t *testing.T) {
	//if
	t.Run("Should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}
		bc := BlockChain()
		if bc.Height != 1 {
			t.Error("BlockChain() should create a blockchain")
		}
	})
	//else
	t.Run("Should restore blockchain", func(t *testing.T) {
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					Height:            2,
					NewestHash:        "test",
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}

		bc := BlockChain()
		if bc.Height != 2 {
			t.Errorf("BlockChain() should restore a blockchain with a height of %d, got %d", 2, bc.Height)
		}
	})
}

func TestBlocks(t *testing.T) {
	fakeBlocks := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			var b *Block
			if fakeBlocks == 0 {
				b = &Block{
					Height:   1,
					PrevHash: "test",
				}
			}
			if fakeBlocks == 1 {
				b = &Block{
					Height: 1,
				}
			}
			fakeBlocks++
			return utils.ToBytes(b)
		},
	}

	bc := &blockchain{}
	blocks := Blocks(bc)
	if reflect.TypeOf(blocks) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of block")
	}
}

const (
	addr string = "40af6376aabb80090ed110a49e70db1861054228eeacb977969d755c1e6e3a281f05d92649ac55b59792b03152aab2bb593399c5662bb96ce17ddebd9d5d686b"
)

func TestMakeTx(t *testing.T) {
	t.Run("Has not enough money", func(t *testing.T) {
		var txs []*Tx
		for range [1]int{} {
			txs = append(txs, makeCoinbaseTx(addr))
		}

		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Transactions: txs,
				}
				return utils.ToBytes(b)
			},
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					NewestHash:        txs[len(txs)-1].ID,
					Height:            1,
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}
		BlockChain()
		_, err := makeTx(addr, "to", 100)
		if err == nil {
			t.Error("makeTx() should return nil")
		}
	})

	t.Run("Has enough money", func(t *testing.T) {
		once = *new(sync.Once)
		var txs []*Tx
		for range [3]int{} {
			txs = append(txs, makeCoinbaseTx(addr))
		}
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Transactions: txs,
				}
				return utils.ToBytes(b)
			},
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					NewestHash:        txs[len(txs)-1].ID,
					Height:            3,
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}
		BlockChain()
		_, err := makeTx(addr, "to", 50)
		if err != nil {
			t.Error("makeTx() should return nil")
		}
	})

	t.Run("addr err", func(t *testing.T) {
		once = *new(sync.Once)
		var txs []*Tx
		for range [3]int{} {
			tx := makeCoinbaseTx(addr)
			txs = append(txs, tx)
		}
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Transactions: txs,
				}
				return utils.ToBytes(b)
			},
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					Height:            3,
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}
		BlockChain()
		_, err := makeTx(addr, "to", 150)
		if err != nil {
			t.Error("makeTx() should return ErrNoMoney")
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("vaildate has error verify", func(t *testing.T) {
		once = *new(sync.Once)
		var txs []*Tx
		for range [3]int{} {
			tx := makeCoinbaseTx(addr)
			txs = append(txs, tx)
		}
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Transactions: txs,
				}
				return utils.ToBytes(b)
			},
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					Height:            3,
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}
		BlockChain()
		tx, _ := makeTx(addr, "to", 150)
		tx.ID = "10af6376aabb80090ed110a49e70db1861054228eeacb977969d755c1e6e3a281f05d92649ac55b59792b03152aab2bb593399c5662bb96ce17ddebd9d5d686b"
		if validate(tx) {
			t.Error("vaildate() sholud be return Verify's false")
		}
	})
	t.Run("vaildate has error prevTx", func(t *testing.T) {
		once = *new(sync.Once)
		var txs []*Tx
		for range [3]int{} {
			tx := makeCoinbaseTx(addr)
			txs = append(txs, tx)
		}
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Transactions: txs,
				}
				return utils.ToBytes(b)
			},
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					Height:            3,
					CurrentDifficulty: 1,
				}
				return utils.ToBytes(bc)
			},
		}
		BlockChain()
		tx, _ := makeTx(addr, "to", 150)
		tx.TxIns[0].TxID = "test"
		if validate(tx) {
			t.Error("이건 진짜 모르것네?")
		}
	})
}

func TestUTxOutsByAddress(t *testing.T) {
	t.Run("UTxOutsByAddress should be break COINBASE", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
					Transactions: []*Tx{
						{
							ID: "test1",
							TxIns: []*TxIn{
								{TxID: "", Index: -1, Signature: "COINBASE"},
							},
							TxOuts: []*TxOut{
								{Address: "testAddress", Amount: 50},
							},
						},
					},
				}
				return utils.ToBytes(b)
			},
		}
		bc := &blockchain{}
		utxOut := UTxOutsByAddress("testAddress", bc)
		if utxOut[0].Amount != 50 {
			t.Error("UTxOutsByAddress should be return amount:50")
		}
	})

	t.Run("UTxOutsByAddress should return NOTCOINBASE", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 3,
					Transactions: []*Tx{
						{
							ID: "test1",
							TxIns: []*TxIn{
								{TxID: "", Index: -1, Signature: "COINBASE"},
								{TxID: "test1", Index: 0, Signature: "NOTCOINBASE1"},
							},
							TxOuts: []*TxOut{
								{Address: "testAddress1", Amount: 1000},
								{Address: "testAddress2", Amount: 1000},
							},
						},
						{
							ID: "test2",
							TxIns: []*TxIn{
								{TxID: "test1", Index: 1, Signature: "NOTCOINBASE2"},
							},
							TxOuts: []*TxOut{
								{Address: "testAddress2", Amount: 100},
							},
						},
					},
				}
				return utils.ToBytes(b)
			},
		}
		bc := &blockchain{NewestHash: "test2"}
		utxOut := UTxOutsByAddress("testAddress2", bc)
		if !(utxOut[0].Amount == 1000 && utxOut[1].Amount == 100) {
			t.Error("UtxOutsByAddress should be utxOuts[0].Amount = 1000 && utxOuts[1].Amount = 100")
		}
	})

}

// func TestFindTx(t *testing.T) {
// 	t.Run("Tx not found", func(t *testing.T) {
// 		dbStorage = fakeDB{
// 			fakeFindBlock: func() []byte {
// 				b := &Block{
// 					Height:       2,
// 					Transactions: []*Tx{},
// 				}
// 				return utils.ToBytes(b)
// 			},
// 		}
// 		tx := FindTx(&blockchain{NewestHash: "testNewestHash"}, "testID")
// 		if tx != nil {
// 			t.Error("Tx should not found")
// 		}
// 	})
// 	t.Run("Tx should be found", func(t *testing.T) {
// 		dbStorage = fakeDB{
// 			fakeFindBlock: func() []byte {
// 				b := &Block{
// 					Height: 2,
// 					Transactions: []*Tx{
// 						{ID: "testID"},
// 					},
// 				}
// 				return utils.ToBytes(b)
// 			},
// 		}
// 		tx := FindTx(&blockchain{NewestHash: "newestHash"}, "testID")
// 		if tx == nil {
// 			t.Error("Transaction should be found")
// 		}
// 	})
// }
