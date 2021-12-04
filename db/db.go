package db

import (
	"fmt"
	"os"

	"github.com/auturnn/kickshaw-coin/utils"
	bolt "go.etcd.io/bbolt"
)

//data races => notion에 정리하기
var db *bolt.DB

const (
	dbName       = "kickshaw"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func getDBName() string {
	return fmt.Sprintf("%s_%s.db", dbName, os.Args[1][6:])
}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		utils.HandleError(err)
		db = dbPointer

		err = db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleError(err)

			_, err = t.CreateBucketIfNotExists([]byte(dataBucket))
			return err
		})
		utils.HandleError(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	fmt.Printf("Saving Block: %s \n", hash)
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		return bucket.Put([]byte(hash), data)
	})
	utils.HandleError(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		return bucket.Put([]byte(checkpoint), data)
	})
	utils.HandleError(err)
}

//Block is get bucket and search to hash data
func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		buk := t.Bucket([]byte(blocksBucket))
		data = buk.Get([]byte(hash))
		return nil
	})
	return data
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleError(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleError(err)
		return nil
	})
}
