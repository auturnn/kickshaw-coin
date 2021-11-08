package utils

import (
	"bytes"
	"encoding/gob"
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