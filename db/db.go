package db

import (
	"fmt"
	"os"

	"github.com/auturnn/kickshaw-coin/utils"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const (
	dbName       = "kickshaw"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
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
	saveCheckpoint(data)
}

func (DB) DeleteAllBlocks() {
	emptyBlocks()
}

func getDBName() string {
	return fmt.Sprintf("%s_%s.db", dbName, os.Args[1][6:])
}

func InitDB() {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		db = dbPointer
		utils.HandleError(err)

		err = db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleError(err)

			_, err = t.CreateBucketIfNotExists([]byte(dataBucket))
			return err
		})
		utils.HandleError(err)
	}
}

func Close() {
	db.Close()
}

func saveBlock(hash string, data []byte) {
	fmt.Printf("Saving Block: %s \n", hash)
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		return bucket.Put([]byte(hash), data)
	})
	utils.HandleError(err)
}

func saveCheckpoint(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		return bucket.Put([]byte(checkpoint), data)
	})
	utils.HandleError(err)
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

func emptyBlocks() {
	db.Update(func(t *bolt.Tx) error {
		utils.HandleError(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleError(err)
		return nil
	})
}
