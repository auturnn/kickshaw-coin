package db

import (
	"fmt"

	"github.com/auturnn/kickshaw-coin/utils"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const (
	dbName       string = "kickshaw"
	dataBucket   string = "data"
	blocksBucket string = "blocks"
	checkpoint   string = "checkpoint"
)

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}

func (DB) LoadChain() []byte {
	return loadChain()
}

func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}

func (DB) SaveChain(data []byte) {
	saveChain(data)
}

func (DB) DeleteAllBlocks() {
	emptyBlocks()
}

//Block is get bucket and search to hash data
func findBlock(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		buk := t.Bucket([]byte(blocksBucket))
		data = buk.Get([]byte(hash))
		return nil
	})
	return data
}

func loadChain() []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func saveBlock(hash string, data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		return bucket.Put([]byte(hash), data)
	})
	utils.HandleError(err, utils.ErrSaveBlock)
}

func saveChain(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		return bucket.Put([]byte(checkpoint), data)
	})
	utils.HandleError(err, utils.ErrSaveChain)
}

func emptyBlocks() {
	db.Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket([]byte(blocksBucket))
		utils.HandleError(err, utils.ErrReplaceBlockchain)

		_, err = t.CreateBucket([]byte(blocksBucket))
		utils.HandleError(err, utils.ErrCreateDB)
		return nil
	})
}

func getDBName() string {
	return fmt.Sprintf("./%s.db", dbName)
}

func Close() {
	db.Close()
}

func InitDB() {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		db = dbPointer
		utils.HandleError(err, utils.ErrLoadDB)

		err = db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleError(err, utils.ErrCreateBlockChain)

			_, err = t.CreateBucketIfNotExists([]byte(dataBucket))
			return err
		})
		utils.HandleError(err, utils.ErrLoadDB)
	}
}
