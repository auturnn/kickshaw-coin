package blockchain

import (
	"errors"
	"time"

	"github.com/auturnn/kickshaw-coin/utils"
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

//이전에 생성된 TxOut을 찾을 수 있을수있는 수단 => TxID, Index
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

/*
Tx = transaction
내가 블록체인에서 거래를 진행할때. 코드는 이전의 블록체인에서 일어난 거래들의 Output을 확인하여 현재의 Balance를 계산하고,
나의 거래 Output이 들어있는 블록들을 송금액보다 크거나, 같을때까지 가져와서 그안의 Amount를 더하여 계산한다.
즉. 송금시(Input) Blockchain의 모든 Output중에 나의 이름이 담긴 블록들을 찾고, 블록안의 Amount를 더한뒤(total) 송금할 금액을 뺀뒤에
거스름돈과 누구에게 얼마가 송금되었는지에 대한 명세를 작성한다.(transaction Output)
Output에는 송금자의 Address(Publick Key)가 들어있다.
거래의 과정을 살펴보면 이전 Tx을 참조하게 되는데.
이전의 Tx로 가게만들어주는 이정표가 바로 TxIn의 TxID, Index이다
구체적으로는 Blockchain전체의 Block들 안에 들어있는 Tx를 전부 가져온후,
TxInput의 TxID(거래ID)와 일치하는 Tx를 찾는다.
Tx를 찾게되면 해당 Tx의 TxOutputs의 Index를 살펴본다(Tx.TxOutputs[TxIn.Index]).
TxOuts에는 Address(PublicKey)가 들어있고 이것을 현재 송금하는 사람의 Sign(송금자PrivateKey, TxID)를 검증한다.=> verify(sign,tx.ID,addr)
->(TxIn을 만들사람이 정말 TxInput에 참조되어있는 TxOutput을 소유한 사람인지를 검증하는것)

ex) 내가 누군가에게 송금할때. 이전에 거래한 내역들과, 나의 사인을 같이 보낸다. (Tx 생성시에 모든 Input에 서명)
시스템은 송금자(Input)인 내가 진짜 해당 금액을 보유하고있는지 확인하고싶다.
시스템은 내가 이번 송금에 동봉한 이전거래내역( Tx(TxIn.TxID).TxOuts[TxIn.Index] )를 통해 나의 이전의 거래내역(Block)들에 적혀있는 잔금명세(Output)를 찾고,
이전거래내역을 찾은 시스템은 이 거래를 진짜 내가 진행한것이 맞는지 확인하기 위해
내가 동봉한 Sign(송금자PrivateKey, TxID)과 이전거래내역에 적혀있는 Address(PubK), tx.ID를 가지고 검증한다.
*검증이 가능한 이유 = 송금(Input)Sign이 나의 PrvKey와 TxID로 이루어져있고, 이전 거래명세(Tx.TxOut)의 Address가 나의 공개키로 이루어져있기 때문에 대조가 가능하다.
(PrivateKey, PublicKey가 wallet 생성시에 자동으로 한쌍이 생성된다.)
검증결과가 True일 경우 나는 시스템이 찾은 이전 거래내역(내 잔고를 확인할수있는.)을 사용할 수 있는 권한이 있다는 것을 증명가능.

나는 진짜 이 이전의 거래내역들이 내가 거래한것들이 맞는지를 증명해야하는것.

(블록 하나에 여러개의 Tx가 존재)
*/

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
	tx.sign()
	if !validate(tx) {
		return nil, ErrorNoMoney
	}

	return tx, nil
}

var ErrorNoMoney = errors.New("not enough Money")
var ErrorNotValid = errors.New("Tx Invaild")

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
