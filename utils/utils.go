package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(err error)  {
	if err != nil{
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte  {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	HandleError(encoder.Encode(i))
	return buf.Bytes()
}

func FromBytes(i interface{}, data []byte)  {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprintf("%v",i)
	hash:= sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}