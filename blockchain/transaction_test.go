package blockchain

import (
	"testing"

	"github.com/auturnn/kickshaw-coin/utils"
)

const (
	testPayload = "00a247a0604f3383e6f176a00f5f5ab10806c311e61c757909fe775c96e6ca96"
	testSign    = "83fc24fc0f552361cbb94c931669451fbabc793c287f6fda94a638b2784ab129718e32235382bfee58ef257212f22c9348b2ba631ba3cd34a26fdb5a3b331bdf"
)

func TestValidate(t *testing.T) {
	t.Run("vaildate FindTx return nil", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				block := &Block{}
				return utils.ToBytes(block)
			},
		}
		tx := &Tx{
			ID: "xx",
			TxIns: []*TxIn{
				{TxID: "", Index: 0, Signature: ""},
			},
		}
		if validate(tx) {
			t.Error("vaildate sholud be return false")
		}
	})

	t.Run("vaildate false Verify", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				block := &Block{
					Transactions: []*Tx{
						{
							ID: testPayload,
							TxOuts: []*TxOut{
								{Address: w.GetAddress()},
							},
						},
					},
				}
				return utils.ToBytes(block)
			},
		}
		tx := &Tx{
			ID: "",
			TxIns: []*TxIn{
				{
					TxID:      testPayload,
					Index:     0,
					Signature: testSign,
				},
			},
		}
		if validate(tx) {
			t.Error("vaildate sholud be return false")
		}
	})

}
