package blockchain

import (
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/auturnn/kickshaw-coin/wallet"
)

const (
	minerReward int = 50
)

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
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, w.GetPrivKey())
	}
}

var ErrorNoMoney = errors.New("not enough Money")
var ErrorNotValid = errors.New("Tx Invaild")

type walletLayer interface {
	GetAddress() string
	GetPrivKey() *ecdsa.PrivateKey
	// InitWallet()
}

var w walletLayer = wallet.WalletLayer{}

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

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, BlockChain()) < amount {
		return nil, ErrorNoMoney
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

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
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
