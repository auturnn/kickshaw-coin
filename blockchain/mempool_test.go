package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"

	"github.com/auturnn/kickshaw-coin/utils"
)

const (
	testAddr string = "00a247a0604f3383e6f176a00f5f5ab10806c311e61c757909fe775c96e6ca96"
)

func createPrivKey() *ecdsa.PrivateKey {
	prk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return prk
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func addrFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func TestAddTx(t *testing.T) {
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

	var txs []*Tx
	for range [3]int{} {
		tx := makeCoinbaseTx(w.GetAddress())
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

	t.Run("should create transaction", func(t *testing.T) {
		mp = Mempool()
		tx, _ := mp.AddTx("xx", 30)
		if tx.TxOuts[0].Amount != 20 {
			t.Errorf("AddTx Sholud return amount of 50, got %d", tx.TxOuts[0].Amount)
		}
	})

	t.Run("should failed add transaction", func(t *testing.T) {
		memOnce = *new(sync.Once)
		mp = Mempool()
		_, err := mp.AddTx("xx", 200)
		if err == nil {
			t.Errorf("AddTx() Sholud not return Error got %s", err)
		}
	})
}

func TestIsOnMempool(t *testing.T) {
	var utxOut *UTxOut
	utxOut = &UTxOut{
		TxID:   "testTxID",
		Index:  0,
		Amount: 50,
	}

	memOnce = *new(sync.Once)
	Mempool().Txs["test"] = &Tx{
		TxIns: []*TxIn{
			{TxID: "testTxID", Index: 0},
		},
	}

	if !isOnMempool(utxOut) {
		t.Error("isOnMempool() should return of true, got false")
	}
}

func TestAddPeerTx(t *testing.T) {
	memOnce = *new(sync.Once)
	mp = Mempool()
	tx := &Tx{
		ID: "test",
	}
	mp.AddPeerTx(tx)
	if _, ok := mp.Txs["test"]; !ok {
		t.Error("AddPeerTx() sholud create map[tx.ID].tx")
	}
}
