package blockchain

import (
	"crypto/ecdsa"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/auturnn/kickshaw-coin/utils"
)

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
	blocks := []*Block{
		{PrevHash: "test"},
		{PrevHash: ""},
	}

	fakeBlocks := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlocks++
			}()
			return utils.ToBytes(blocks[fakeBlocks])
		},
	}

	bc := &blockchain{}
	blocksResult := Blocks(bc)
	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of block")
	}
}

func TestGetDifficulty(t *testing.T) {
	t.Run("getDefficulty sholud CurrentDifficulty +1", func(t *testing.T) {
		type test struct {
			height int
			want   int
		}
		blocks := []*Block{
			{PrevHash: "test"},
			{PrevHash: "test"},
			{PrevHash: "test"},
			{PrevHash: "test"},
			{PrevHash: ""},
		}

		fakeblock := 0
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				defer func() {
					fakeblock++
				}()
				return utils.ToBytes(blocks[fakeblock])
			},
		}
		tests := []test{
			{height: 0, want: defaultDiffculty},
			{height: 2, want: defaultDiffculty},
			{height: 5, want: 3},
		}
		for _, tc := range tests {
			bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDiffculty}
			got := getDifficulty(bc)
			if got != tc.want {
				t.Errorf("getDifiiculty() should return %d got %d", tc.want, got)
			}
		}
	})

	t.Run("getDefficulty sholud CurrentDifficulty -1", func(t *testing.T) {
		type test struct {
			height int
			want   int
		}
		blocks := []*Block{
			{PrevHash: "test", Timestamp: int(time.Now().Unix() - 300)},
			{PrevHash: "test", Timestamp: int(time.Now().Unix() - 300)},
			{PrevHash: "test", Timestamp: int(time.Now().Unix() - 300)},
			{PrevHash: "test", Timestamp: int(time.Now().Unix() - 300)},
			{PrevHash: ""},
		}

		fakeblock := 0
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				defer func() {
					fakeblock++
				}()
				return utils.ToBytes(blocks[fakeblock])
			},
		}
		tests := []test{
			{height: 0, want: defaultDiffculty},
			{height: 2, want: defaultDiffculty},
			{height: 4, want: defaultDiffculty},
			{height: 5, want: 1},
		}
		for _, tc := range tests {
			bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDiffculty}
			got := getDifficulty(bc)
			if got != tc.want {
				t.Errorf("getDifiiculty() should return %d got %d", tc.want, got)
			}
		}
	})
}

func TestAddPeerBlock(t *testing.T) {
	t.Run("", func(t *testing.T) {
		mp.Txs["test"] = &Tx{}
		newBlock := &Block{
			Difficulty: 2,
			Hash:       "test",
			Transactions: []*Tx{
				{ID: "test"},
			},
		}
		bc = &blockchain{}
		bc.AddPeerBlock(newBlock)
		if mp.Txs["test"] != nil {
			t.Error("AddPeerBlock() should do not has transaction")
		}
	})
}

func TestReplace(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	blocks := []*Block{
		{Difficulty: 2, Hash: "test"},
		{Difficulty: 2, Hash: "test"},
	}
	bc.Replace(blocks)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("Replace() should mutate the blockchain")
	}
}

func TestUTxOutsByAddress(t *testing.T) {
	fw = &fakeWallet{}
	fw.priv = createPrivKey()
	fw.addr = addrFromKey(fw.priv)

	w = fakeWalletLayer{
		fakeGetAddress: func() string {
			return fw.addr
		},
		fakeGetPrivKey: func() *ecdsa.PrivateKey {
			return fw.priv
		},
	}

	tx := makeCoinbaseTx(w.GetAddress())
	txs := []*Tx{
		tx,
		{
			ID: "",
			TxIns: []*TxIn{
				{TxID: tx.ID, Index: 0},
			},
			TxOuts: []*TxOut{
				{Address: "to", Amount: 50},
			},
		},
	}

	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			block := &Block{
				PrevHash:     "",
				Transactions: txs,
			}
			return utils.ToBytes(block)
		},
	}

	utxOuts := UTxOutsByAddress(w.GetAddress(), BlockChain())
	total := 0
	for _, utxOut := range utxOuts {
		total += utxOut.Amount
	}
	if total != 50 {
		t.Error("Error!")
	}
}
