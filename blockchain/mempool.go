package blockchain

import "sync"

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
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

func (mp *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(w.GetAddress(), to, amount)
	if err != nil {
		return nil, err
	}
	mp.Txs[tx.ID] = tx
	return tx, nil
}

func (mp *mempool) TxToConfirm() []*Tx {
	//coinbase의 모든 거래내역을 가져옴
	coinbase := makeCoinbaseTx(w.GetAddress())
	//거래내역에 coinbase 거래내역을 추가
	var txs []*Tx
	for _, tx := range mp.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	//confirm이 끝나면 memory pool에서 비워주어야함
	mp.Txs = make(map[string]*Tx)
	return txs
}

func (mp *mempool) AddPeerTx(tx *Tx) {
	mp.m.Lock()
	defer mp.m.Unlock()

	mp.Txs[tx.ID] = tx
}
