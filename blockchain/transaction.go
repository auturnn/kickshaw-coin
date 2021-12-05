package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/auturnn/kickshaw-coin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	mt  sync.Mutex
}

var mp *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		mp = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return mp
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		//check. need to money in blockchain.
		prevTx := FindTx(BlockChain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		addr := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, addr)
		if !valid {
			break
		}
	}

	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
OuterLoop: // label
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break OuterLoop
			}

		}
	}
	return exists
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, bc) < amount {
		return nil, errors.New("not enough money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0

	uTxOuts := UTxOutsByAddress(from, BlockChain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	tx.sign()

	if !validate(tx) {
		return nil, ErrorNoMoney
	}

	return tx, nil
}

var ErrorNoMoney = errors.New("not enough Money")
var ErrorNotValid = errors.New("Tx Invaild")

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	//coinbase의 모든 거래내역을 가져옴
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	//거래내역에 coinbase 거래내역을 추가
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	//confirm이 끝나면 memory pool에서 비워주어야함
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.mt.Lock()
	defer m.mt.Unlock()

	m.Txs[tx.ID] = tx
}
