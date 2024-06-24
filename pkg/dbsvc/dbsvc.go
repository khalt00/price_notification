package dbsvc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/khalt00/price_notification/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type PriceNotification struct {
	Symbol string
	High   float64
	Low    float64
}

func PutLevelDB(db *leveldb.DB, key int64, value []PriceNotification) error {
	// Create a bytes buffer to hold the encoded value
	var buffer bytes.Buffer
	// Create a new encoder and encode the value into the buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("failed to encode value: %v", err)
	}

	// Convert key to bytes
	keyBytes, err := utils.Int64ToBytes(key)
	if err != nil {
		log.Fatal(err)
	}

	// Put the encoded value into the database
	if err := db.Put(keyBytes, buffer.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to put value in leveldb: %v", err)
	}

	return nil
}

func GetLevelDB(db *leveldb.DB, key int64) ([]PriceNotification, error) {
	var value []PriceNotification

	// Convert key to bytes
	keyBytes, err := utils.Int64ToBytes(key)

	// Get the value from the database
	data, err := db.Get(keyBytes, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return value, fmt.Errorf("key not found in leveldb")
		}
		return value, fmt.Errorf("failed to get value from leveldb: %v", err)
	}

	// Create a bytes buffer from the retrieved data
	buffer := bytes.NewBuffer(data)
	// Create a new decoder and decode the data into the value
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("failed to decode value: %v", err)
	}

	return value, nil
}
