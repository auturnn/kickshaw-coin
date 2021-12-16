package blockchain

import (
	"testing"

	"github.com/auturnn/kickshaw-coin/utils"
)

func TestMakeTx(t *testing.T) {
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {

			return utils.ToBytes("")
		},
		fakeLoadChain: func() []byte {
			bc = &blockchain{
				Height:            2,
				NewestHash:        "test",
				CurrentDifficulty: 1,
			}
			return utils.ToBytes(bc)
		},
	}
}
