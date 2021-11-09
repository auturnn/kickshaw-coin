package db

import (
	"fmt"

	"github.com/auturnn/kickshaw-coin/utils"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbName = "blockchain.db"
	dataBucket = "data"
	blocksBucket = "blocks"
	checkpoint = "checkpoint"
)

func DB() *bolt.DB {
	if db == nil{
		dbPointer, err := bolt.Open(dbName, 0600, nil)
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

func Close(){
	DB().Close()
}

func SaveBlock(hash string, data []byte)  {
	fmt.Printf("Saving Block: %s \n", hash)
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		return bucket.Put([]byte(hash), data)
	})
	utils.HandleError(err)
}

func SaveBlockChain(data []byte)   {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		return bucket.Put([]byte(checkpoint), data)
	})
	utils.HandleError(err)
}

//Block is get bucket and search to hash data
func Block(hash string) []byte  {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		buk := t.Bucket([]byte(blocksBucket))
		data = buk.Get([]byte(hash))
		return nil
	})
	return data
}

func Checkpoint() []byte  {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}