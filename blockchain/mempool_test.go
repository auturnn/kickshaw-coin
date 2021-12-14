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
	addr     string = "7fbe04dd8fb0acab6b3b9cba531f103e45f23a67b6f12be7150e6ca2122de14c8660a582a6a8d8450d96fd96353d44ab424a88fb8e0435312294243bfc63a615"
	testAddr        = "00a247a0604f3383e6f176a00f5f5ab10806c311e61c757909fe775c96e6ca96"
)

type fakeWallet struct {
	addr string
	priv *ecdsa.PrivateKey
}

var fw *fakeWallet

type fakeWalletLayer struct {
	fakeGetAddress func() string
	fakeGetPrivKey func() *ecdsa.PrivateKey
	fakeInitWallet func()
}

func (fakeWalletLayer) InitWallet() {}

func (f fakeWalletLayer) GetAddress() string {
	return f.fakeGetAddress()
}

func (f fakeWalletLayer) GetPrivKey() *ecdsa.PrivateKey {
	return f.fakeGetPrivKey()
}

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

	t.Run("should add transaction", func(t *testing.T) {
		tx, _ := mp.AddTx("test", 50)
		if tx.TxOuts[0].Amount != 50 {
			t.Errorf("AddTx Sholud return amount of 50, got %d", tx.TxOuts[0].Amount)
		}
	})

	t.Run("should failed add transaction", func(t *testing.T) {
		memOnce = *new(sync.Once)
		_, err := mp.AddTx("test", 200)
		if err == nil {
			t.Error("AddTx() Sholud not return transaction")
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
	mp = Mempool()
	mp = &mempool{
		Txs: map[string]*Tx{
			"test": {
				TxIns: []*TxIn{
					{TxID: "testTxID", Index: 0, Signature: ""},
				},
			},
		},
	}

	if !isOnMempool(utxOut) {
		t.Error("isOnMempool() should return of true, got false")
	}
}
