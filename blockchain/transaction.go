package blockchain

import (
	"errors"
	"time"

	"github.com/auturnn/kickshaw-coin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
OuterLoop: // label
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break OuterLoop
			}

		}
	}
	return exists
}

/*
	ex) 내 거래내역이 50, 50을 얻었다고 가정.
	이전거래내역의 Output은 [50, 50], 이중 70을 A에게 보내고 싶다면
	Output[:1]까지의 거래내역을 Input에 추가한뒤, 잔금을 돌려받는 시스템.

	단순히 숫자를 통한 계산이 아닌 Block(거래매물,화폐)을 통한 시스템이기 때문에
	인터넷 뱅킹과 같은 시스템이 아닌, 실물화폐를 생각할 경우가 이해하기 쉽다.
	지갑에 5000원 2장이 있을때 7000원 상당의 물건을 거래할 경우
	5000원 2장을 모두 꺼내어 거래상대에게 주고, 잔금을 치뤄받는 형식.
*/
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
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	// from(wallet.Wallet().Address)은 나중에 추가.
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

//채굴시 작성
func (m *mempool) TxToConfirm() []*Tx {
	//coinbase의 모든 거래내역을 가져옴
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	//모든 mempool내역을 가져온다
	txs := m.Txs
	//거래내역에 coinbase 거래내역을 추가
	txs = append(txs, coinbase)
	//confirm이 끝나면 memory pool에서 비워주어야함
	m.Txs = nil
	return txs
}
