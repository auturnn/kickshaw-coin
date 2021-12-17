package blockchain

import (
	"crypto/ecdsa"
	"reflect"
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

func TestCreateblock(t *testing.T) {
	dbStorage = fakeDB{}
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
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("test")
		if err == nil {
			t.Error("The block should not be found")
		}
	})
	t.Run("Block is found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}
		block, _ := FindBlock("test")
		if block.Height != 1 {
			t.Error("The block should be found")
		}
	})
}
